package group

import (
	"go_im/im/dao"
	"go_im/pkg/logger"
)

func LoadAllGroup() map[int64]*Group {
	res := map[int64]*Group{}
	groups, err := dao.GroupDao.GetAllGroup()
	if err != nil {
		logger.E("Init group error", err)
		return res
	}
	for _, g := range groups {
		group, err := initGroup(g)
		if err != nil {
			logger.E("Init group error", err)
			continue
		}
		res[g.Gid] = group
	}
	return res
}

func LoadGroup(gid int64) (*Group, error) {
	dbGroup, err := dao.GroupDao.GetGroup(gid)
	if err != nil {
		logger.E("load group error", err)
		return nil, err
	}
	group, err := initGroup(dbGroup)
	return group, err
}

func initGroup(dbGroup *dao.Group) (*Group, error) {
	group := newGroup(dbGroup.Gid, dbGroup.ChatId)
	group.mute = dbGroup.Mute

	if dbGroup.ChatId <= 0 {
		chat, err := dao.ChatDao.CreateChat(dao.ChatTypeGroup, 0, dbGroup.Gid)
		if err != nil {
			return nil, err
		}
		err = dao.GroupDao.UpdateGroupChatId(dbGroup.Gid, chat.Cid)
		if err != nil {
			return nil, err
		}
		group.cid = chat.Cid
	} else {
		// todo restore msg sequence
		//group.msgSequence = 0
	}

	members, err := dao.GroupDao.GetMembers(dbGroup.Gid)
	if err != nil {
		return nil, err
	}

	for _, member := range members {
		group.PutMember(member.Uid, member.Flag)
	}
	return group, nil
}
