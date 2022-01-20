package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
)

const configEnv = "IM_CONFIG"

var (
	MySql          *MySqlConf
	Redis          *RedisConf
	IMService      *IMServiceConf
	ApiHttpService *ApiHttpServiceConf
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

type IMServiceConf struct {
	Addr    string
	Service string
	Port    int
}

type ApiHttpServiceConf struct {
	Addr string
	Port int
}

type config struct {
	MySql     MySqlConf
	Redis     RedisConf
	IMService IMServiceConf
	ApiHttp   ApiHttpServiceConf
}

func init() {
	var conf config
	env, b := os.LookupEnv(configEnv)
	if !b {
		panic("the config file location is not configured in env, please configure env IM_CONFIG")
	}
	_, err := toml.DecodeFile(env, &conf)

	//_, err := toml.DecodeFile("config.toml", &conf)
	if err != nil {
		panic(fmt.Sprintf("error on load config: %s", err.Error()))
	}
	MySql = &conf.MySql
	Redis = &conf.Redis
	IMService = &conf.IMService
	ApiHttpService = &conf.ApiHttp
}
