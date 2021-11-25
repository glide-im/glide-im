package apidep

import (
	"go_im/im/message"
	"go_im/pkg/logger"
)

type MockClientManager struct {
}

func (MockClientManager) ClientSignIn(oldUid int64, uid int64, device int64) {
	logger.D("ClientSignIn, oldUid=%d, uid=%d, device=%d", oldUid, uid, device)
}

func (MockClientManager) ClientLogout(uid int64, device int64) {
	logger.D("ClientLogout, uid=%d, device=%d", uid, device)
}

func (MockClientManager) EnqueueMessage(uid int64, device int64, message *message.Message) {
	logger.D("EnqueueMessage, uid=%d, device=%d, msg=%v", uid, device, message)
}

func (MockClientManager) IsDeviceOnline(uid, device int64) bool {
	logger.D("IsDeviceOnline, uid=%d, device=%d", uid, device)
	return false
}

func (MockClientManager) IsOnline(uid int64) bool {
	logger.D("IsOnline, uid=%d", uid)
	return false
}

func (MockClientManager) AllClient() []int64 {
	logger.D("AllClient")
	return []int64{}
}
