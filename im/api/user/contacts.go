package user

import (
	"go_im/im/api/apidep"
	"go_im/im/api/comm"
	"go_im/im/api/router"
	"go_im/im/dao/userdao"
	"go_im/im/message"
)

func (a *UserApi) DeleteContact(ctx *route.Context, request *DeleteContactsRequest) error {
	// TODO 2021-11-29
	return nil
}

func (a *UserApi) UpdateContactRemark(ctx *route.Context, request *UpdateRemarkRequest) error {
	// TODO 2021-11-29
	return nil
}

func (a *UserApi) ContactApproval(ctx *route.Context, request *ContactApproval) error {
	// TODO 2021-11-29
	return nil
}

func (a *UserApi) AddContact(ctx *route.Context, request *AddContacts) error {

	if ctx.Uid == request.Uid {
		return errAddSelf
	}
	hasUser, err := userdao.UserInfoDao.HasUser(request.Uid)
	if err != nil {
		return comm.NewDbErr(err)
	}
	if !hasUser {
		return errUserNotExist
	}

	isC, err := userdao.ContactsDao.HasContacts(ctx.Uid, request.Uid, 1)
	if err != nil {
		return comm.NewDbErr(err)
	}
	if isC {
		return errAlreadyContacts
	}
	// TODO 2021-11-29 use transaction
	err = userdao.ContactsDao.AddContacts(ctx.Uid, request.Uid, contactsTypeUser)
	if err != nil {
		return comm.NewDbErr(err)
	}
	err = userdao.ContactsDao.AddContacts(request.Uid, ctx.Uid, contactsTypeUser)
	if err != nil {
		return comm.NewDbErr(err)
	}

	ctx.Response(message.NewMessage(ctx.Seq, comm.ActionSuccess, ""))

	m := NewContactResponse{
		Uid:  request.Uid,
		Type: contactsTypeUser,
	}
	apidep.SendMessage(request.Uid, 0, message.NewMessage(-1, comm.ActionNewContact, m))
	return nil
}

//goland:noinspection GoPreferNilSlice
func (a *UserApi) GetContactList(ctx *route.Context) error {

	contacts, err := userdao.ContactsDao.GetContacts(ctx.Uid)
	if err != nil {
		return comm.NewDbErr(err)
	}

	resp := []ContactResponse{}
	for _, contact := range contacts {
		if contact.Type == contactsTypeGroup {
			_ = apidep.GroupManager.MemberOnline(contact.Id, ctx.Uid)
		}
		resp = append(resp, ContactResponse{
			Id:     contact.Id,
			Type:   contact.Type,
			Remark: contact.Remark,
		})
	}

	ctx.Response(message.NewMessage(ctx.Seq, comm.ActionSuccess, resp))
	return nil
}
