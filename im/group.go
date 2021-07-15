package im

import "go_im/im/entity"

type Group struct {
	*mutex

	Gid  uint64
	Name string

	member   []int64
	memberCh map[int64]chan *entity.Message
}

func NewGroup(gid uint64, name string, member []int64) *Group {
	ret := new(Group)
	ret.Gid = gid
	ret.Name = name
	ret.member = member
	return ret
}

func (g *Group) GetOnlineMember() []int64 {
	defer g.LockUtilReturn()()

	m := make([]int64, 1)
	for k := range g.memberCh {
		m = append(m, k)
	}
	return m
}

func (g *Group) SendMessage(uid int64, message *entity.Message) {
	defer g.LockUtilReturn()()

	for i := range g.memberCh {
		g.memberCh[i] <- message
	}
}

func (g *Group) Subscribe(client *Client) {
	defer g.LockUtilReturn()()

	g.memberCh[client.uid] = client.messages
}

func (g *Group) Unsubscribe(client *Client) {
	defer g.LockUtilReturn()()

	delete(g.memberCh, client.uid)
}

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
