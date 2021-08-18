package main

import (
	"go_im/config"
	"go_im/im"
	"go_im/pkg/db"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func init() {
	config.Init()
	db.Init()
}

func main() {
	logToHttpServe()
	im.Run()
}

func logToHttpServe() {
	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:8081", nil))
	}()
}
