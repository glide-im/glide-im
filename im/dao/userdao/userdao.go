package userdao

type Cache interface {
	GetUserLoginState()
}

type IUserDao interface {
	HasUser(uid int64) (bool, error)
	GetUser(uid int64) (*User, error)
}

func GetUserSignInState() {

}

func HasID() (bool, error) {
	return false, nil
}

func HasUser(uid int64) (bool, error) {
	return false, nil
}

func GetUser(uid int64) (*User, error) {
	return nil, nil
}
