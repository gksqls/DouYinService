package service

import (
	rdb "DouYinService/cache"
	"context"

	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
	"xorm.io/xorm/names"
)

var (
	Redis *rdb.Cache
	Db    *xorm.Engine
)

func initial_db() {
	engine, err := xorm.NewEngine("sqlite3", "nicosoft.db")
	if err != nil {
		panic(err)
	}
	engine.ShowSQL(true)
	engine.Logger().SetLevel(log.LOG_INFO)
	engine.SetTableMapper(names.SnakeMapper{})
	engine.SetColumnMapper(names.SnakeMapper{})
	err = engine.Ping()
	if err != nil {
		panic(err)
	}
	Db = engine
}

func Initial(cxt context.Context) {
	Redis = &rdb.Cache{Context: cxt}
	// Redis.Initial()
	initial_db()
}
