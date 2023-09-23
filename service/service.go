package service

import (
	"DouYinService/db"
	"context"
)

var (
	Redis *db.Cache
)

func Initial(cxt context.Context) {
	Redis = &db.Cache{Context: cxt}
	Redis.Initial()
}
