package im

import (
	"go_im/im/dao"
	"go_im/im/entity"
)

type Group struct {
	*mutex

	Gid   int64
	Cid   int64
	group *dao.Group

	members map[int64]*dao.GroupMember
}

func NewGroup(gid int64, group *dao.Group, cid int64, member []*dao.GroupMember) *Group {
	ret := new(Group)
	ret.mutex = NewMutex()
	ret.members = map[int64]*dao.GroupMember{}
	ret.Gid = gid
	ret.Cid = cid
	ret.group = group
	for _, m := range member {
		ret.members[m.Uid] = m
	}
	return ret
}

func (g *Group) PutMember(member *dao.GroupMember) {
	g.members[member.Uid] = member
}

func (g *Group) RemoveMember(uid int64) {
	delete(g.members, uid)
}

func (g *Group) HasMember(uid int64) bool {
	_, ok := g.members[uid]
	return ok
}

func (g *Group) IsMemberOnline(uid int64) bool {
	return false
}

func (g *Group) GetOnlineMember() []*dao.GroupMember {
	defer g.LockUtilReturn()()

	var online []*dao.GroupMember
	for id, member := range g.members {
		if ClientManager.IsOnline(id) {
			online = append(online, member)
		}
	}
	return online
}

func (g *Group) GetMembers() []*dao.GroupMember {
	members := make([]*dao.GroupMember, 0, len(g.members))
	for _, v := range g.members {
		members = append(members, v)
	}
	return members
}

func (g *Group) SendMessage(uid int64, message *entity.Message) {
	defer g.LockUtilReturn()()
	logger.D("Group.SendMessage: %s", message)

	for id := range g.members {
		ClientManager.EnqueueMessage(id, message)
	}
}
