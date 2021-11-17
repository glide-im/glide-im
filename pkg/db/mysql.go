package db

import (
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"go_im/config"
	"go_im/pkg/logger"
	"runtime"
	"time"
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

	DB.DB().SetMaxOpenConns(10000)
	DB.DB().SetMaxIdleConns(1000)
	DB.DB().SetConnMaxIdleTime(time.Minute * 5)

	if err != nil {
		panic(err)
	}
	//DB.LogMode(true)
	DB.SingularTable(true)
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "im_" + defaultTableName
	}
	initRedis()
}

func initRedis() {

	conf := config.Redis
	Redis = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password:     conf.Password,
		DB:           conf.Db,
		PoolSize:     runtime.NumCPU() * 30,
		MinIdleConns: 10,
	})

	logger.D("redis init: %s", Redis.Ping().Val())
}
