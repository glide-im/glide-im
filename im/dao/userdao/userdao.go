package userdao

var Dao = UserDao{
	UserInfoDaoInterface: UserInfoDao,
	ContactsDaoInterface: ContactsDao,
}

type Cache interface {
	GetUserLoginState(uid int64) ([]*LoginState, error)
	AuthToken(token string) (int64, error)
}

type UserInfoDaoInterface interface {
	AddUser(u *User) error
	DelUser(uid int64) error
	HasUser(uid int64) (bool, error)

	UpdateNickname(uid int64, nickname string) error
	UpdateAvatar(uid int64, avatar string) error
	UpdatePassword(uid int64, password string) error
	GetPassword(uid int64) (string, error)

	GetUserInfo(uid int64) (*User, error)
	GetUserSimpleInfo(uid ...int64) ([]*User, error)
}

type ContactsDaoInterface interface {
	AddContacts(uid int64, id int64, type_ int8) error
	DelContacts(uid int64, id int64, type_ int8) error
	GetContacts(uid int64) ([]*Contacts, error)
}

type UserDao struct {
	Cache
	UserInfoDaoInterface
	ContactsDaoInterface
}
