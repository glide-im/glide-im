package service

import (
	"testing"
)

func TestGetConfig(t *testing.T) {
	config, err := GetServiceConfig()
	if err != nil {
		t.Error(err)
	}
	t.Log(config)
}
