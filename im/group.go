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

	members  map[int64]*dao.GroupMember
	memberCh map[int64]chan *entity.Message
}

func NewGroup(gid int64, group *dao.Group, cid int64, member []*dao.GroupMember) *Group {
	ret := new(Group)
	ret.mutex = new(mutex)
	ret.members = map[int64]*dao.GroupMember{}
	ret.memberCh = map[int64]chan *entity.Message{}
	ret.Gid = gid
	ret.group = group
	for _, m := range member {
		ret.memberCh[m.Uid] = nil
		ret.members[m.Uid] = m
	}
	return ret
}

func (g *Group) PutMember(member *dao.GroupMember, c chan *entity.Message) {
	g.memberCh[member.Uid] = c
	g.members[member.Uid] = member
}

func (g *Group) HasMember(uid int64) bool {
	_, ok := g.members[uid]
	return ok
}

func (g *Group) IsMemberOnline(uid int64) bool {
	return g.memberCh[uid] != nil
}

func (g *Group) GetOnlineMember() []*dao.GroupMember {
	defer g.LockUtilReturn()()

	onlineMember := make([]*dao.GroupMember, 1)
	for k, v := range g.memberCh {
		m, exist := g.members[k]
		if v != nil && exist {
			onlineMember = append(onlineMember, m)
		}
	}
	return onlineMember
}

func (g *Group) SendMessage(uid int64, message *entity.Message) error {
	defer g.LockUtilReturn()()

	for _, v := range g.memberCh {
		if v != nil {
			v <- message
		}
	}
	return nil
}

func (g *Group) Subscribe(uid int64, mc chan *entity.Message) {
	defer g.LockUtilReturn()()
	g.memberCh[uid] = mc
}

func (g *Group) Unsubscribe(uid int64) {
	defer g.LockUtilReturn()()
	g.memberCh[uid] = nil
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
