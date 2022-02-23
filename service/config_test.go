package service

import (
	"github.com/BurntSushi/toml"
	"testing"
)

func TestGetConfig(t *testing.T) {
	c := Configs{}
	_, err := toml.DecodeFile("example_config.toml", &c)
	if err != nil {
		t.Error(err)
	}
	t.Log(c.Api.Client)
}
