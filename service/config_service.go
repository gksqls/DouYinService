package service

import "github.com/gin-gonic/gin"

type Config struct {
	Context *gin.Context
}

type TConfig struct {
	Id         int    `json:"id" xorm:"not null integer"`
	ServerPort int    `json:"server_port" xorm:"integer"`
	DouyinRoom string `json:"douyin_room" xorm:"text"`
}

func (c *Config) FindConfig() *TConfig {
	conf := new(TConfig)
	_, err := Db.Where("id=?", 1).Get(conf)
	if err != nil {
		panic(err)
	}
	return conf
}

func (c *Config) UpdateConfig() {
	conf := new(TConfig)
	c.Context.BindJSON(&conf)
	_, err := Db.Where("id=?", 1).Update(conf)
	if err != nil {
		panic(err)
	}
}
