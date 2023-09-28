package douyin

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	dyproto "DouYinService/protobuf"
	"DouYinService/service"
	"DouYinService/socket"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"golang.org/x/exp/maps"
)

type Room struct {
	Url       string
	Ttwid     string
	RoomStore string
	RoomId    string
	RoomTitle string
	wsConnect *websocket.Conn
}

var (
	giftCountCache = make(map[string]map[int]time.Time)
	lock           sync.Mutex
	welcomes       []TWelcome
)

type TWelcome struct {
	Id      int    `xorm:"not null pk autoincr INTEGER"`
	User    string `xorm:"not null TEXT"`
	Content string `xorm:"not null TEXT"`
}

func gift_count_timer_elapsed() {
	now := time.Now()
	for {
		for k1 := range giftCountCache {
			val := maps.Values(giftCountCache[k1])[0]
			if now.Add(10 * time.Second).After(val) {
				lock.Lock()
				delete(giftCountCache, k1)
				lock.Unlock()
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func NewRoom(u string) (*Room, error) {
	h := map[string]string{
		"accept":     "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
		"cookie":     "__ac_nonce=0638733a400869171be51",
	}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for k, v := range h {
		req.Header.Set(k, v)
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer res.Body.Close()
	data := res.Cookies()
	var ttwid string
	for _, c := range data {
		if c.Name == "ttwid" {
			ttwid = c.Value
			break
		}
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	resText := string(body)
	re := regexp.MustCompile(`roomId\\":\\"(\d+)\\"`)
	match := re.FindStringSubmatch(resText)
	if match == nil || len(match) < 2 {
		log.Println("No match found")
		return nil, err
	}
	liveRoomId := match[1]
	return &Room{
		Url:    u,
		Ttwid:  ttwid,
		RoomId: liveRoomId,
	}, nil
}

func (r *Room) Connect() error {
	wsUrl := "wss://webcast3-ws-web-lq.douyin.com/webcast/im/push/v2/?app_name=douyin_web&version_code=180800&webcast_sdk_version=1.3.0&update_version_code=1.3.0&compress=gzip&internal_ext=internal_src:dim|wss_push_room_id:%s|wss_push_did:%s|dim_log_id:202302171547011A160A7BAA76660E13ED|fetch_time:1676620021641|seq:1|wss_info:0-1676620021641-0-0|wrds_kvs:WebcastRoomStatsMessage-1676620020691146024_WebcastRoomRankMessage-1676619972726895075_AudienceGiftSyncData-1676619980834317696_HighlightContainerSyncData-2&cursor=t-1676620021641_r-1_d-1_u-1_h-1&host=https://live.douyin.com&aid=6383&live_id=1&did_rule=3&debug=false&endpoint=live_pc&support_wrds=1&im_path=/webcast/im/fetch/&user_unique_id=%s&device_platform=web&cookie_enabled=true&screen_width=1440&screen_height=900&browser_language=zh&browser_platform=MacIntel&browser_name=Mozilla&browser_version=5.0%20(Macintosh;%20Intel%20Mac%20OS%20X%2010_15_7)%20AppleWebKit/537.36%20(KHTML,%20like%20Gecko)%20Chrome/110.0.0.0%20Safari/537.36&browser_online=true&tz_name=Asia/Shanghai&identity=audience&room_id=%s&heartbeatDuration=0&signature=00000000"
	wsUrl = strings.Replace(wsUrl, "%s", r.RoomId, -1)
	h := http.Header{}
	h.Set("cookie", "ttwid="+r.Ttwid)
	h.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
	wsConn, _, err := websocket.DefaultDialer.Dial(wsUrl, h)
	if err != nil {
		return err
	}
	r.wsConnect = wsConn
	go r.read()
	go r.send()
	go gift_count_timer_elapsed()
	service.Db.Find(&welcomes)
	return nil
}

func (r *Room) Close() error {
	return r.wsConnect.Close()
}

func (r *Room) read() {
	for {
		_, data, err := r.wsConnect.ReadMessage()

		if err != nil {
			panic(err.Error())
		}

		var msgPack dyproto.PushFrame
		_ = proto.Unmarshal(data, &msgPack)
		decompressed, _ := degzip(msgPack.Payload)
		var payloadPackage dyproto.Response
		_ = proto.Unmarshal(decompressed, &payloadPackage)

		if payloadPackage.NeedAck {
			r.sendAck(msgPack.LogId, payloadPackage.InternalExt)
		}

		for _, msg := range payloadPackage.MessagesList {
			switch msg.Method {
			case "WebcastChatMessage":
				parseChatMsg(msg.Payload)
			case "WebcastGiftMessage":
				parseGiftMsg(msg.Payload)
			case "WebcastLikeMessage":
				parseLikeMsg(msg.Payload)
			case "WebcastMemberMessage":
				parseEnterMsg(msg.Payload)
			}
		}
	}
}

func (r *Room) send() {
	for {
		pingPack := &dyproto.PushFrame{
			PayloadType: "bh",
		}
		data, _ := proto.Marshal(pingPack)
		err := r.wsConnect.WriteMessage(websocket.BinaryMessage, data)
		if err != nil {
			panic(err.Error())
		}
		time.Sleep(time.Second * 10)
	}
}

func (r *Room) sendAck(logId uint64, iExt string) {
	ackPack := &dyproto.PushFrame{
		LogId:       logId,
		PayloadType: iExt,
	}
	data, _ := proto.Marshal(ackPack)
	err := r.wsConnect.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		panic(err.Error())
	}
}

func degzip(data []byte) ([]byte, error) {
	b := bytes.NewReader(data)
	var out bytes.Buffer
	r, err := gzip.NewReader(b)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(&out, r)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func parseChatMsg(msg []byte) {
	var chatMsg dyproto.ChatMessage
	_ = proto.Unmarshal(msg, &chatMsg)
	message := &socket.WsMessage{
		Type:       2,
		NickName:   chatMsg.User.NickName,
		HeadImgUrl: chatMsg.User.AvatarThumb.UrlListList[0],
		Content:    chatMsg.Content,
	}
	message.Pust()
}

func parseGiftMsg(msg []byte) {
	var giftMsg dyproto.GiftMessage
	_ = proto.Unmarshal(msg, &giftMsg)

	key := strconv.FormatUint(giftMsg.GiftId, 10) + "-" + strconv.FormatUint(giftMsg.GroupId, 10)

	if giftMsg.RepeatEnd == 1 {
		if _, exists := giftCountCache[key]; exists && giftMsg.GroupId > 0 {
			lock.Lock()
			delete(giftCountCache, key)
			lock.Unlock()
		}
		return
	}

	lastCount := 0
	currCount := int(giftMsg.RepeatCount)
	backward := currCount <= lastCount

	if currCount <= 0 {
		currCount = 1
	}

	if _, exists := giftCountCache[key]; exists {
		keys := maps.Keys(giftCountCache[key])[0]
		lastCount = keys
		backward = currCount <= lastCount
		if !backward {
			lock.Lock()
			giftCountCache[key] = map[int]time.Time{currCount: time.Now()}
			lock.Unlock()
		}
	} else {
		if giftMsg.GroupId > 0 && !backward {
			lock.Lock()
			giftCountCache[key] = map[int]time.Time{currCount: time.Now()}
			lock.Unlock()
		}
	}

	if backward {
		return
	}

	count := currCount - lastCount

	message := &socket.WsMessage{
		Type:      1,
		GiftName:  giftMsg.Gift.Name,
		GiftCount: count,
		GiftId:    giftMsg.GiftId,
		NickName:  giftMsg.User.NickName,
	}
	message.Pust()
}

func parseLikeMsg(msg []byte) {
	var likeMsg dyproto.LikeMessage
	_ = proto.Unmarshal(msg, &likeMsg)
	message := &socket.WsMessage{
		Type:     3,
		NickName: likeMsg.User.NickName,
		Count:    likeMsg.Count,
	}
	message.Pust()
}

func parseEnterMsg(msg []byte) {
	var enterMsg dyproto.MemberMessage
	_ = proto.Unmarshal(msg, &enterMsg)
	for _, welcome := range welcomes {
		if strings.ContainsAny(enterMsg.User.NickName, welcome.User) {
			message := &socket.WsMessage{
				Type:       4,
				NickName:   enterMsg.User.NickName,
				HeadImgUrl: enterMsg.User.AvatarThumb.UrlListList[0],
				Content:    fmt.Sprintf(welcome.Content, enterMsg.User.NickName),
			}
			message.Pust()
		}
	}
}
