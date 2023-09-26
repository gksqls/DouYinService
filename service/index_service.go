package service

import (
	"github.com/gin-gonic/gin"
)

type Index struct {
	Context *gin.Context
}

type PluginList struct {
	Name   string `json:"name"`
	Uri    string `json:"uri"`
	Enable bool   `json:"enable"`
}

const (
	MANAGE_PLUGIN_LIST = "MANAGE:PLUGIN:"
)

func (i *Index) LoadPluginList() *[]PluginList {
	ret := []PluginList{}
	return &ret
}
