package im

import (
	"go_im/im/dao"
	"go_im/im/entity"
)

type groupApi struct{}

func (m *groupApi) GetGroupMember(c *Client, seq int64, gid uint64) error {
	members := GroupManager.GetGroup(gid).member

	msg := entity.NewMessage(seq, entity.ActionSuccess)
	if err := msg.SetData(members); err != nil {
		return err
	}
	c.EnqueueMessage(msg)
	return nil
}

func (m *groupApi) GetGroupInfo(c *Client, seq int64, gid uint64) {
	dao.GroupDao.GetGroup(gid)

}

func (m *groupApi) RemoveMember(c *Client, seq int64, gid uint64, uid int64) {

}

func (m *groupApi) AddMember(c *Client, seq int64, gid uint64, uid int64) {

}
