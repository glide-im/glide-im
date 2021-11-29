package msgdao

import (
	"go_im/pkg/db"
	"testing"
	"time"
)

func init() {
	db.Init()
}

func TestSessionDaoImpl_CreateSession(t *testing.T) {
	err := SessionDaoImpl.CreateSession(1, 2)
	if err != nil {
		t.Error(err)
	}
}

func TestSessionDao_UpdateOrInitSession(t *testing.T) {
	err := SessionDaoImpl.UpdateOrInitSession(1, 2, time.Now().Unix())
	if err != nil {
		t.Error(err)
	}
}

func TestSessionDao_GetRecentSession(t *testing.T) {
	session, err := SessionDaoImpl.GetRecentSession(time.Now().Unix() - 1000)
	if err != nil {
		t.Error(err)
	}
	for _, s := range session {
		t.Log(s)
	}
}
