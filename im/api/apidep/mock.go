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

type MockGroupManager struct {
}

func (m *MockGroupManager) MemberOnline(gid int64, uid int64) error {
	logger.D("MemberOnline, gid=%d, uid=%d", gid, uid)
	return nil
}

func (m *MockGroupManager) MemberOffline(gid int64, uid int64) error {
	logger.D("MemberOffline, gid=%d, uid=%d", gid, uid)
	return nil
}
func (g *MockGroupManager) UpdateMember(gid int64, uid int64, flag int64) error {
	logger.D("UpdateMember, gid=%d, uid=%d, flag=%d", gid, uid, flag)
	return nil
}

func (m *MockGroupManager) PutMember(gid int64, mb []int64) error {
	logger.D("PutMember, gid=%d, mb=%v", gid, mb)
	return nil
}

func (m *MockGroupManager) RemoveMember(gid int64, uid ...int64) error {
	logger.D("RemoveMember, gid=%d, uid=%v", gid, uid)
	return nil
}

func (m *MockGroupManager) CreateGroup(gid int64) error {
	logger.D("CreateGroup, gid=%d", gid)
	return nil
}

func (m *MockGroupManager) DissolveGroup(gid int64) error {
	logger.D("DissolveGroup, gid=%d", gid)
	return nil
}

func (m *MockGroupManager) MuteGroup(gid int64, mute bool) error {
	logger.D("MuteGroup, gid=%d, mute=%v", gid, mute)
	return nil
}

func (m *MockGroupManager) DispatchNotifyMessage(gid int64, message *message.Message) error {
	logger.D("DispatchNotifyMessage, gid=%d, message=%v", gid, message)
	return nil
}
