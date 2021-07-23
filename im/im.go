package im

import "go_im/im/dao"

func Run() {

	dao.Init()
	NewWsServer(nil).Run()
}
