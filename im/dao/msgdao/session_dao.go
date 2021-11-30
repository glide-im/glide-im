package msgdao

import (
	"go_im/im/dao/common"
	"go_im/pkg/db"
	"strconv"
	"time"
)

var SessionDaoImpl SessionDao = &sessionDaoImpl{}

type sessionDaoImpl struct{}

func getSessionId(from int64, to int64) string {
	lg, sm := from, to
	if lg < sm {
		lg, sm = sm, lg
	}
	return strconv.FormatInt(lg, 10) + "_" + strconv.FormatInt(sm, 10)
}

func (s *sessionDaoImpl) GetSession(uid int64, uid2 int64) (*Session, error) {
	sid := getSessionId(uid, uid2)
	var se Session
	query := db.DB.Model(&Session{}).Where("session_id = ?", sid).Find(&se)
	return &se, query.Error
}

func (s *sessionDaoImpl) CreateSession(uid1 int64, uid2 int64, updateAt int64) (*Session, error) {
	sid := getSessionId(uid1, uid2)
	se := &Session{
		SessionId: sid,
		Uid:       uid1,
		Uid2:      uid2,
		LastMID:   0,
		UpdateAt:  updateAt,
		CreateAt:  time.Now().Unix(),
	}
	create := db.DB.Create(se)
	if err := common.MustUpdate(create); err != nil {
		return nil, err
	}
	return se, nil
}

func (s *sessionDaoImpl) UpdateOrInitSession(uid1 int64, uid2 int64, update int64) error {
	sid := getSessionId(uid1, uid2)

	query := db.DB.Model(&Session{}).
		Where("`session_id` = ?", sid).
		Update("update_at", update)
	if query.Error != nil {
		return query.Error
	}
	if query.RowsAffected == 0 {
		create := db.DB.Create(&Session{
			SessionId: sid,
			Uid:       uid1,
			Uid2:      uid2,
			LastMID:   0,
			UpdateAt:  update,
			CreateAt:  time.Now().Unix(),
		})
		if err := common.MustUpdate(create); err != nil {
			return err
		}
	}
	return nil
}

func (s *sessionDaoImpl) GetRecentSession(uid int64, updateAfter int64) ([]*Session, error) {
	var se []*Session
	query := db.DB.Model(&Session{}).Where("(`uid` = ? OR `uid2` = ?) AND `update_at` > ?", uid, uid, updateAfter).Find(&se)
	if query.Error != nil {
		return nil, query.Error
	}
	return se, nil
}
