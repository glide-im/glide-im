package userdao

import "time"

var Dao = UserDao{
	Cache:                UserCacheDao{},
	UserInfoDaoInterface: UserInfoDao,
	ContactsDaoInterface: ContactsDao,
}

type Cache interface {
	//GetUserSignState(uid int64) ([]*LoginState, error)
	//IsUserSignIn(uid int64, device int64) (bool, error)
	//DelToken(token string) error
	//DelAllToken(uid int64) error
	//GetTokenInfo(token string) (int64, int64, error)
	//SetSignInToken(uid int64, device int64, token string, expiredAt time.Duration) error

	DelAuthToken(uid int64, device int64) error
	SetTokenVersion(uid int64, device int64, version int64, expiredAt time.Duration) error
	GetTokenVersion(uid int64, device int64) (int64, error)
}

type UserInfoDaoInterface interface {
	AddUser(u *User) error
	DelUser(uid int64) error
	HasUser(uid int64) (bool, error)

	UpdateNickname(uid int64, nickname string) error
	UpdateAvatar(uid int64, avatar string) error
	UpdatePassword(uid int64, password string) error
	GetPassword(uid int64) (string, error)

	GetUidInfoByLogin(account string, password string) (int64, error)
	GetUser(uid int64) (*User, error)
	GetUserSimpleInfo(uid ...int64) ([]*User, error)
}

type ContactsDaoInterface interface {
	HasContacts(uid int64, id int64, type_ int8) (bool, error)
	AddContacts(uid int64, id int64, type_ int8) error
	DelContacts(uid int64, id int64, type_ int8) error
	GetContacts(uid int64) ([]*Contacts, error)
	GetContactsByType(uid int64, type_ int) ([]*Contacts, error)
}

type UserDao struct {
	Cache
	UserInfoDaoInterface
	ContactsDaoInterface
}
