package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GiftDanceLeft(c *gin.Context) {
	c.HTML(http.StatusOK, "gift_left.tmpl", nil)
}

func GiftDanceFoot(c *gin.Context) {
	c.HTML(http.StatusOK, "gift_foot.tmpl", nil)
}

func GiftGameLeft(c *gin.Context) {
	c.HTML(http.StatusOK, "gift_game_left.tmpl", nil)
}

func GiftGameFoot(c *gin.Context) {
	c.HTML(http.StatusOK, "gift_game_foot.tmpl", nil)
}
