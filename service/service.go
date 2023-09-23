package service

import (
	"DouYinService/db"
	"context"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
	"xorm.io/xorm/names"
)

type config struct {
	MySQL struct {
		Host     string `yaml:"host"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Db       string `yaml:"db"`
	} `yaml:"mysql"`
}

var (
	conf  config
	Redis *db.Cache
	MySQL *xorm.Engine
)

func initial_mysql() {
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		panic(err)
	}
	engine, err := xorm.NewEngine(
		"mysql",
		fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8", conf.MySQL.User, conf.MySQL.Password, conf.MySQL.Host, conf.MySQL.Db),
	)
	if err != nil {
		panic(err)
	}

	engine.ShowSQL(true)
	engine.Logger().SetLevel(log.LOG_INFO)
	engine.SetTableMapper(names.SameMapper{})
	engine.SetColumnMapper(names.SnakeMapper{})

	err = engine.Ping()
	if err != nil {
		panic(err)
	}
	_, err = engine.Exec("SELECT VERSION()")
	if err != nil {
		panic(err)
	}
	MySQL = engine
}

func Initial(cxt context.Context) {
	Redis = &db.Cache{Context: cxt}
	Redis.Initial()
	initial_mysql()
}
