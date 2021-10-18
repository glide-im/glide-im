package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

var (
	MySql *MySqlConf
	Redis *RedisConf
)

type MySqlConf struct {
	Host     string
	Port     int
	Username string
	Password string
	Db       string
	Charset  string
}

type RedisConf struct {
	Host     string
	Port     int
	Password string
	Db       int
}

type config struct {
	MySql MySqlConf
	Redis RedisConf
}

func Init() {
	var conf config
	_, err := toml.DecodeFile("config.toml", &conf)
	if err != nil {
		panic(fmt.Sprintf("error on load config: %s", err.Error()))
	}
	MySql = &conf.MySql
	Redis = &conf.Redis
}
