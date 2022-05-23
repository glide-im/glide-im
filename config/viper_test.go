package config

import (
	"github.com/spf13/viper"
	"testing"
)

func TestLoad(t *testing.T) {
	err := Load()

	if err != nil {
		t.Errorf("Error loading config: %s", err)
	}

	t.Log(viper.GetString("ApiHttp.Port"))
}
