package main

import (
	"go_im/config"
	"go_im/pkg/db"
	_ "net/http/pprof"
)

func init() {
	config.Init()
	db.Init()
}

func main() {

}
