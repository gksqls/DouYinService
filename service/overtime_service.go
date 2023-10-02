package service

import (
	"DouYinService/socket"
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

type TGift struct {
	Name string `json:"name" xorm:"not null text"`
	Code int    `json:"code" xorm:"not null int"`
}

func (t *TGift) Save() {
	gift := new(TGift)
	has, err := Db.Where("code=?", t.Code).Get(gift)
	if err != nil {
		panic(err)
	}
	if !has {
		_, err = Db.InsertOne(t)
		if err != nil {
			panic(err)
		}
	}
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

func (o *Overtime) FindGiftList() *[]TGift {
	var list []TGift
	err := Db.Find(&list)
	if err != nil {
		panic(err)
	}
	return &list
}
