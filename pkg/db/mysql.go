package db

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"go_im/config"

	_ "github.com/go-sql-driver/mysql"
)

var (
	DB    *gorm.DB
	Redis *redis.Client
)

func Init() {
	conf := config.MySql
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true",
		conf.Username, conf.Password, conf.Host, conf.Port, conf.Db, conf.Charset)
	var err error
	DB, err = gorm.Open("mysql", url)
	if err != nil {
		panic(err)
	}
}
