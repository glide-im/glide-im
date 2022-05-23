package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"go_im/pkg/logger"
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
	Addr      string
	Port      int
	IMService struct {
		Addr string
		Port int
		Etcd []string
		Name string
	}
}

type config struct {
	MySql     MySqlConf
	Redis     RedisConf
	IMService IMServiceConf
	ApiHttp   ApiHttpServiceConf
}

func init() {
	var conf config

	c := getConfigPath()
	logger.D("config path: %s", c)

	_, err := toml.DecodeFile(c, &conf)
	if err != nil {
		panic(fmt.Sprintf("error on load config: %s", err.Error()))
	}
	MySql = &conf.MySql
	Redis = &conf.Redis
	IMService = &conf.IMService
	ApiHttpService = &conf.ApiHttp
}

func getConfigPath() string {
	configPath := ""
	for i, arg := range os.Args {
		println(i, arg)
	}
	if len(os.Args) == 2 {
		configPath = os.Args[1]
	} else {
		configPath = "config.toml"
	}
	f, e := os.Open(configPath)
	if e == nil {
		_ = f.Close()
		return configPath
	}
	return readEnv()
}

func readEnv() string {
	configPath, b := os.LookupEnv(configEnv)
	if !b {
		panic("the config file location is not configured in env, please configure env IM_CONFIG")
	}
	return configPath
}
