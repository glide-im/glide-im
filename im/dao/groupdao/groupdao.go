package groupdao

var Dao = GroupDao{
	GroupMemberDao: GroupMemberDaoImpl{},
	GroupInfoDao:   &GroupInfoDaoImpl{},
}

type GroupInfoDao interface {
	CreateGroup(name string, flag int) (*GroupModel, error)
	GetGroup(gid int64) (*GroupModel, error)
	GetAllGroup() ([]*GroupModel, error)
	GetGroups(gid ...int64) ([]*GroupModel, error)
	UpdateGroupName(gid int64, name string) error
	UpdateGroupAvatar(gid int64, avatar string) error
	UpdateGroupMute(gid int64, mute bool) error
	GetGroupMute(gid int64) (bool, error)
	UpdateGroupFlag(gid int64, flag int) error
	GetGroupFlag(gid int64) (int, error)
	HasGroup(gid int64) (bool, error)
}

type GroupMemberDao interface {
	HasMember(gid int64, uid int64) (bool, error)
	GetMembers(gid int64) ([]*GroupMemberModel, error)
	AddMember(gid int64, uid int64, typ int64, defaultFlag int64) error
	AddMembers(gid int64, flag int64, typ int64, uid ...int64) error
	RemoveMember(gid int64, uid int64) error
	UpdateMemberFlag(gid int64, uid int64, flag int) error
	GetMemberFlag(gid int64, uid int64) (int64, error)
	UpdateMemberType(gid int64, uid int64, flag int) error
	GetMemberType(gid int64, uid int64) (int64, error)
	GetMember(gid int64, uid int64) (*GroupMemberModel, error)
}

type GroupDao struct {
	GroupInfoDao
	GroupMemberDao
}
