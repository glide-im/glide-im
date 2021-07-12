package im

import "go_im/im/entity"

type Group struct {
	*mutex

	Gid  int64
	Name string

	online []chan *entity.Message
}

func (g *Group) SendMessage(message *entity.Message) {
	defer g.LockUtilReturn()

	for i := range g.online {
		g.online[i] <- message
	}
}

func (g *Group) Subscribe(client *Client) {

}

func (g *Group) Unsubscribe(client *Client) {

}
