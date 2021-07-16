package im

import (
	"go_im/im/dao"
	"go_im/im/entity"
)

type groupApi struct{}

func (m *groupApi) GetGroupMember(msg *ApiMessage, gid uint64) error {
	members := GroupManager.GetGroup(gid).member

	resp := entity.NewMessage(msg.seq, entity.RespActionSuccess)
	if err := resp.SetData(members); err != nil {
		return err
	}
	ClientManager.GetClient(msg.uid)

	ClientManager.EnqueueMessage(msg.uid, resp)
	return nil
}

func (m *groupApi) GetGroupInfo(msg *ApiMessage, gid uint64) {
	dao.GroupDao.GetGroup(gid)
	ClientManager.EnqueueMessage(msg.uid, entity.NewSimpleMessage(msg.seq, entity.RespActionSuccess, "get group info"))
}

func (m *groupApi) RemoveMember(msg *ApiMessage, gid uint64, uid int64) error {

	if msg.uid != uid {
		// check permission
	}

	GroupManager.GetGroup(gid).Unsubscribe(uid)
	err := dao.GroupDao.RemoveMember(gid, uid)

	if err != nil {
		return err
	}

	resp := entity.NewSimpleMessage(msg.seq, entity.RespActionSuccess, "remove member success")
	if msg.uid == uid {
		resp.Data = []byte("exit group success")
	} else {
		resp1 := entity.NewSimpleMessage(0, entity.RespActionGroupRemoved, "you have been removed from the group xxx by xxx")
		ClientManager.EnqueueMessage(uid, resp1)
	}

	ClientManager.EnqueueMessage(msg.uid, resp)
	return nil
}

func (m *groupApi) AddMember(msg *ApiMessage, gid uint64, uid int64) {

}
