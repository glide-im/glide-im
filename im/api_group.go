package im

import (
	"errors"
	"go_im/im/dao"
	"go_im/im/entity"
)

type groupApi struct{}

func (m *groupApi) CreateGroup(msg *ApiMessage, request *entity.CreateGroupRequest) error {

	group, err := dao.GroupDao.NewGroup(request.Name, msg.uid)
	if err != nil {
		return err
	}

	resp := entity.NewMessage(msg.seq, entity.RespActionSuccess)
	if err = resp.SetData(group); err != nil {
		return err
	}
	ClientManager.EnqueueMessage(msg.uid, resp)
	return nil
}

func (m *groupApi) GetGroupMember(msg *ApiMessage, request *entity.GetGroupMemberRequest) error {

	members, err := dao.GroupDao.GetMembers(request.Gid)
	if err != nil {
		return err
	}

	ms := make([]*entity.GroupMemberResponse, len(members))

	for _, member := range members {
		ms = append(ms, &entity.GroupMemberResponse{
			Uid:        member.Uid,
			Nickname:   "",
			RemarkName: member.Remark,
			Type:       member.Type,
			Online:     ClientManager.IsOnline(member.Uid),
		})
	}

	resp := entity.NewMessage(msg.seq, entity.RespActionSuccess)
	if err = resp.SetData(ms); err != nil {
		return err
	}

	ClientManager.EnqueueMessage(msg.uid, resp)
	return nil
}

func (m *groupApi) GetGroupInfo(msg *ApiMessage, gid int64) {
	ClientManager.EnqueueMessage(msg.uid, entity.NewSimpleMessage(msg.seq, entity.RespActionSuccess, "get group info"))
}

func (m *groupApi) RemoveMember(msg *ApiMessage, request *entity.RemoveMemberRequest) error {

	for _, uid := range request.Uid {
		GroupManager.GetGroup(request.Gid).Unsubscribe(uid)
		err := dao.GroupDao.RemoveMember(request.Gid, uid)

		if err != nil {
			return err
		}
		notifyResp := entity.NewSimpleMessage(-1, entity.RespActionGroupRemoved, "you have been removed from the group xxx by xxx")
		ClientManager.EnqueueMessage(uid, notifyResp)
	}

	resp := entity.NewSimpleMessage(msg.seq, entity.RespActionSuccess, "remove member success")

	ClientManager.EnqueueMessage(msg.uid, resp)
	return nil
}

func (m *groupApi) AddGroupMember(msg *ApiMessage, request *entity.AddMemberRequest) error {

	g := GroupManager.GetGroup(request.Gid)

	for _, uid := range request.Uid {
		if err := dao.GroupDao.AddMember(request.Gid, uid, 1); err != nil {
			return err
		}
		client := ClientManager.GetClient(msg.uid)
		if client != nil {
			client.AddGroup(g)
		}
	}

	resp := entity.NewSimpleMessage(msg.seq, entity.RespActionSuccess, "add member success")
	ClientManager.EnqueueMessage(msg.uid, resp)

	return nil
}

func (m *groupApi) ExitGroup(msg *ApiMessage, request *entity.ExitGroupRequest) error {

	g := GroupManager.GetGroup(request.Gid)
	g.Unsubscribe(msg.uid)

	err := dao.GroupDao.RemoveMember(request.Gid, msg.uid)
	if err != nil {
		return err
	}
	resp := entity.NewSimpleMessage(msg.seq, entity.RespActionSuccess, "exit group success")
	ClientManager.EnqueueMessage(msg.uid, resp)
	return err
}

func (m *groupApi) JoinGroup(msg *ApiMessage, request *entity.JoinGroupRequest) error {

	g := GroupManager.GetGroup(request.Gid)

	if err := dao.GroupDao.AddMember(request.Gid, msg.uid, 1); err != nil {
		return err
	}

	chat, err := dao.ChatDao.NewChat(msg.uid, request.Gid, 2)

	if err != nil {
		return err
	}

	client := ClientManager.GetClient(msg.uid)
	if client == nil {
		return errors.New("client state exception")
	}
	client.AddGroup(g)
	resp := entity.NewMessage(msg.seq, entity.RespActionSuccess)
	if err = resp.SetData(chat); err != nil {
		return err
	}
	ClientManager.EnqueueMessage(msg.uid, resp)

	return nil
}
