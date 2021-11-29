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
	return strconv.FormatInt(from, 10) + "_" + strconv.FormatInt(to, 10)
}

func (s *sessionDaoImpl) CreateSession(from int64, to int64) error {
	sid := getSessionId(from, to)
	create := db.DB.Create(&Session{
		SessionId: sid,
		Uid:       from,
		To:        to,
		ReadAt:    0,
		LastMID:   0,
		UpdateAt:  0,
		CreateAt:  time.Now().Unix(),
	})
	return common.MustUpdate(create)
}

func (s *sessionDaoImpl) UpdateOrInitSession(from int64, to int64, update int64) error {
	sid := getSessionId(from, to)

	query := db.DB.Model(&Session{}).
		Where("`session_id` = ?", sid).
		Update("update_at", update)
	if query.Error != nil {
		return query.Error
	}
	if query.RowsAffected == 0 {
		create := db.DB.Create(&Session{
			SessionId: sid,
			Uid:       from,
			To:        to,
			ReadAt:    0,
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

func (s *sessionDaoImpl) GetRecentSession(updateAfter int64) ([]*Session, error) {
	var se []*Session
	query := db.DB.Model(&Session{}).Where("`update_at` > ?", updateAfter).Find(&se)
	if query.Error != nil {
		return nil, query.Error
	}
	return se, nil
}
