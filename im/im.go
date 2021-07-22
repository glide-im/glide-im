package im

import "go_im/im/dao"

func Run() {

	dao.InitUserDao()
	NewWsServer(nil).Run()
}
