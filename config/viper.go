package config

import "github.com/spf13/viper"

var (
	MySql       *MySqlConf
	Redis       *RedisConf
	WsServer    *WsServerConf
	ApiHttp     *ApiHttpConf
	IMRpcServer *IMRpcServerConf
)

type WsServerConf struct {
	Addr string
	Port int
}

type ApiHttpConf struct {
	Addr string
	Port int
}

type IMRpcServerConf struct {
	Addr        string
	Port        int
	Network     string
	Etcd        []string
	Name        string
	EnableGroup bool
}

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

func Load() error {

	viper.SetConfigName("config.toml")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/")
	viper.AddConfigPath("$HOME/.config/")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	c := struct {
		MySql       *MySqlConf
		Redis       *RedisConf
		WsServer    *WsServerConf
		ApiHttp     *ApiHttpConf
		IMRpcServer *IMRpcServerConf
	}{}

	err = viper.Unmarshal(&c)
	if err != nil {
		return err
	}
	Redis = c.Redis
	MySql = c.MySql
	WsServer = c.WsServer
	ApiHttp = c.ApiHttp
	IMRpcServer = c.IMRpcServer

	return err
}
