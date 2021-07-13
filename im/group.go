package im

import "go_im/im/entity"

type Group struct {
	*mutex

	Gid  int64
	Name string
	Mute bool

	memberOnline []int64
	memberCh     map[int64]chan *entity.Message
}

func (g *Group) Online() []int64 {
	return g.memberOnline
}

func (g *Group) SendMessage(message *entity.Message) {
	defer g.LockUtilReturn()()

	for i := range g.memberCh {
		g.memberCh[i] <- message
	}
}

func (g *Group) Subscribe(client *Client) {
	defer g.LockUtilReturn()()
	g.memberOnline = append(g.memberOnline, client.uid)
	g.memberCh[client.uid] = client.messages
}

func (g *Group) Unsubscribe(client *Client) {
	defer g.LockUtilReturn()()
	// delete ele memberOnline
	delete(g.memberCh, client.uid)
}

type GroupMap struct {
	*mutex
	groups map[uint64]*Group
}

func NewGroupMap() *GroupMap {
	ret := new(GroupMap)
	ret.groups = make(map[uint64]*Group)
	return ret
}

func (g *GroupMap) Get(gid uint64) *Group {
	defer g.LockUtilReturn()()
	return g.groups[gid]
}

func (g *GroupMap) Put(gid uint64, group *Group) {
	defer g.LockUtilReturn()()
	g.groups[gid] = group
}

func (g *GroupMap) Delete(gid uint64) {
	defer g.LockUtilReturn()()
	delete(g.groups, gid)
}
