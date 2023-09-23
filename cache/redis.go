package cache

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v2"
)

type config struct {
	Redis struct {
		Address  string `yaml:"address"`
		Password string `yaml:"password"`
		Db       int    `yaml:"db"`
	} `yaml:"redis"`
}

var (
	conf config
	rdb  *redis.Client
)

type Cache struct {
	Context context.Context
}

// 初始化Redis
func (c *Cache) Initial() {
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		panic(err)
	}
	rdb = redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Address,
		Password: conf.Redis.Password,
		DB:       conf.Redis.Db,
	})
	r, _ := rdb.Ping(c.Context).Result()
	log.Printf("[INFO] Redis ping result -> %s", r)
}

func (c *Cache) Get(key, val string) string {
	ret, err := rdb.Get(c.Context, "key").Result()
	if err != nil {
		panic(err)
	}
	return ret
}

func (c *Cache) Set(key string, val interface{}, timeout time.Duration) error {
	return rdb.Set(c.Context, key, val, timeout).Err()
}

func (c *Cache) Del(key string) error {
	return rdb.Del(c.Context, key).Err()
}
