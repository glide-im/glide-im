package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

var (
	MySql *MySqlConf
)

type MySqlConf struct {
	Host     string
	Port     int
	Username string
	Password string
	Db       string
	Charset  string
}

type config struct {
	MySql MySqlConf
}

func init() {
	var conf config
	_, err := toml.DecodeFile("config.toml", &conf)
	if err != nil {
		panic(fmt.Sprintf("error on load config: %s", err.Error()))
	}
	MySql = &conf.MySql
}
