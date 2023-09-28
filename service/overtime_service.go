package service

import (
	"DouYinService/socket"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Overtime struct {
	Context *gin.Context
}

type TOvertime struct {
	Id        int `json:"id" xorm:"not null integer"`
	GiftAddId int `json:"gift_add_id" xorm:"not null integer"`
	GiftSubId int `json:"gift_sub_id" xorm:"not null integer"`
	TimeAdd   int `json:"time_add" xorm:"not null integer"`
	TimeSub   int `json:"time_sub" xorm:"not null integer"`
}

func (o *Overtime) Set() {
	jsonMap := make(map[string]interface{})
	o.Context.BindJSON(&jsonMap)
	t, _ := strconv.Atoi(jsonMap["type"].(string))
	cd, _ := strconv.Atoi(jsonMap["time"].(string))
	message := &socket.WsMessage{Type: t, Countdown: cd}
	message.Pust()
}

func (o *Overtime) Find() *TOvertime {
	overtime := new(TOvertime)
	_, err := Db.Where("id=?", 1).Get(overtime)
	if err != nil {
		panic(err)
	}
	fmt.Println(overtime)
	return overtime
}

func (o *Overtime) Save() {
	overtime := new(TOvertime)
	o.Context.BindJSON(&overtime)
	_, err := Db.Where("id=?", 1).Update(overtime)
	if err != nil {
		panic(err)
	}
	message := &socket.WsMessage{
		Type:      1021,
		GiftAddId: overtime.GiftAddId,
		GiftSubId: overtime.GiftSubId,
		TimeAdd:   overtime.TimeAdd,
		TimeSub:   overtime.TimeSub,
	}
	message.Pust()
}
