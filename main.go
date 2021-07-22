package main

import (
	"go_im/config"
	"go_im/im"
	"go_im/pkg/db"
)

func init() {
	config.Init()
	db.Init()
}

func main() {

	im.Run()
}
