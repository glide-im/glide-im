package group

import (
	"go_im/im/client"
	"go_im/im/comm"
	"go_im/im/dao"
	"go_im/im/message"
	"go_im/pkg/logger"
)

type Group struct {
	*comm.Mutex

	Gid   int64
	Cid   int64
	Group *dao.Group

	members *groupMemberMap
}

func NewGroup(gid int64, group *dao.Group, cid int64, member []*dao.GroupMember) *Group {
	ret := new(Group)
	ret.Mutex = comm.NewMutex()
	ret.members = newGroupMemberMap()
	ret.Gid = gid
	ret.Cid = cid
	ret.Group = group
	for _, m := range member {
		ret.members.Put(m.Uid, m)
	}
	return ret
}

func (g *Group) PutMember(member *dao.GroupMember) {
	g.members.Put(member.Uid, member)
}

func (g *Group) RemoveMember(uid int64) {
	g.members.Delete(uid)
}

func (g *Group) HasMember(uid int64) bool {
	return g.members.Contain(uid)
}

func (g *Group) IsMemberOnline(uid int64) bool {
	return false
}

func (g *Group) GetOnlineMember() []*dao.GroupMember {
	var online []*dao.GroupMember
	for _, member := range g.members.members {
		// TODO 2021-9-9 17:12:58
		online = append(online, member)
	}
	return online
}

func (g *Group) GetMembers() []*dao.GroupMember {
	defer g.LockUtilReturn()()
	members := make([]*dao.GroupMember, 0, g.members.Size())
	for _, v := range g.members.members {
		members = append(members, v)
	}
	return members
}

func (g *Group) SendMessage(uid int64, message *message.Message) {
	logger.D("Group.SendMessage: %s", message)

	for id := range g.members.members {
		client.EnqueueMessage(id, message)
	}
}

////////////////////////////////////////////////////////////////////////////////

type groupMemberMap struct {
	*comm.Mutex
	members map[int64]*dao.GroupMember
}

func newGroupMemberMap() *groupMemberMap {
	ret := new(groupMemberMap)
	ret.Mutex = new(comm.Mutex)
	ret.members = make(map[int64]*dao.GroupMember)
	return ret
}

func (g *groupMemberMap) Size() int {
	return len(g.members)
}

func (g *groupMemberMap) Get(id int64) *dao.GroupMember {
	defer g.LockUtilReturn()()
	member, ok := g.members[id]
	if ok {
		return member
	}
	return nil
}

func (g *groupMemberMap) Contain(id int64) bool {
	_, ok := g.members[id]
	return ok
}

func (g *groupMemberMap) Put(id int64, member *dao.GroupMember) {
	defer g.LockUtilReturn()()
	g.members[id] = member
}

func (g *groupMemberMap) Delete(id int64) {
	defer g.LockUtilReturn()()
	delete(g.members, id)
}
