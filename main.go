package main

import (
	"go_im/im"
	"go_im/pkg/db"
)

func main() {

	db.Init()
	im.Run()
}
