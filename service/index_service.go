package service

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
)

type Index struct {
	Context *gin.Context
}

type PluginList struct {
	Wooden struct {
		Enable bool `json:"enable"`
	} `json:"wooden"`
	Barrage struct {
		Enable bool `json:"enable"`
	} `json:"barrage"`
	Countdown struct {
		Enable bool `json:"enable"`
	} `json:"countdown"`
	Overtime struct {
		Enable bool `json:"enable"`
	} `json:"overtime"`
}

const (
	MANAGE_PLUGIN_LIST = "MANAGE:PLUGIN:LIST"
)

func (i *Index) LoadPluginList() *PluginList {
	data := Redis.Get(MANAGE_PLUGIN_LIST)
	list := new(PluginList)
	if data != "" {
		err := json.Unmarshal([]byte(data), &list)
		if err != nil {
			log.Printf("[ERROR] LoadPluginList unmarshal json error -> %s", err.Error())
		}
	} else {
		list.Barrage.Enable = false
		list.Wooden.Enable = false
		list.Countdown.Enable = false
		list.Overtime.Enable = false
		val, err := json.Marshal(list)
		if err != nil {
			log.Printf("[ERROR] LoadPluginList marshal json error -> %s", err.Error())
		} else {
			Redis.Set(MANAGE_PLUGIN_LIST, string(val), 0)
		}
	}
	return list
}
