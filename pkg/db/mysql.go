package db

import (
	"fmt"
	"github.com/glide-im/glideim/config"
	l "github.com/glide-im/glideim/pkg/logger"
	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"

	"runtime"
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
	DB, err = gorm.Open(mysql.Open(url), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "im_",
			SingularTable: true,
			//NameReplacer:  nil,
			//NoLowerCase:   false,
		},
	})
	if err != nil {
		panic(err)
	}
	db, err := DB.DB()
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(10000)
	db.SetMaxIdleConns(1000)
	db.SetConnMaxLifetime(time.Minute * 10)
	db.SetConnMaxIdleTime(time.Minute * 6)

	//DB.LogMode(true)
	//DB.SingularTable(true)
	//gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
	//	return "im_" + defaultTableName
	//}
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

	l.D("redis init: %s", Redis.Ping().Val())
}
