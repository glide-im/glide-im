package group

import (
	"go_im/im/dao/groupdao"
	"go_im/pkg/logger"
)

func LoadAllGroup() map[int64]*Group {
	res := map[int64]*Group{}
	groups, err := groupdao.GroupDao2.GetAllGroup()
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
	dbGroup, err := groupdao.GroupDao2.GetGroup(gid)
	if err != nil {
		logger.E("load group error", err)
		return nil, err
	}
	group, err := initGroup(dbGroup)
	return group, err
}

func initGroup(dbGroup *groupdao.Group) (*Group, error) {
	group := newGroup(dbGroup.Gid)
	group.mute = dbGroup.Mute

	members, err := groupdao.GroupDao2.GetMembers(dbGroup.Gid)
	if err != nil {
		return nil, err
	}

	for _, member := range members {
		group.PutMember(member.Uid, newMemberInfo())
	}
	return group, nil
}
