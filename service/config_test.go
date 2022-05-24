package service

import (
	"testing"
)

func TestGetConfig(t *testing.T) {
	cf, err := GetConfig()
	if err != nil {
		t.Error(err)
	}
	t.Log(cf)
}
