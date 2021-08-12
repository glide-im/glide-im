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

func (g *Group) HasMember(uid int64) bool {
	_, ok := g.members[uid]
	return ok
}

func (g *Group) IsMemberOnline(uid int64) bool {
	// TODO
	return false
}

func (g *Group) GetOnlineMember() []*dao.GroupMember {
	defer g.LockUtilReturn()()

	return []*dao.GroupMember{}
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

	for id, _ := range g.members {
		ClientManager.EnqueueMessage(id, message)
	}
}

func (g *Group) Unsubscribe(uid int64) {
	defer g.LockUtilReturn()()

}

//////////////////////////////////////////////////////////////////////////////////

type Int64Set struct {
	m map[int64]interface{}
}

func (i *Int64Set) Add(v int64) {
	if i.Contain(v) {
		return
	}
	i.m[v] = nil
}

func (i *Int64Set) Remove(v int64) {
	_, ok := i.m[v]
	if ok {
		delete(i.m, v)
	}
}

func (i *Int64Set) Size() int {
	return len(i.m)
}

func (i *Int64Set) Contain(v int64) bool {
	_, ok := i.m[v]
	return ok
}
