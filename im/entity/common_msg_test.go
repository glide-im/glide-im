package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestActionValue(t *testing.T) {

	t.Logf("\n%b\n%b\n%b", MaskActionApi, ActionUserLogin, RespActionFailed)
	assert.Equal(t, ActionUserLogin, MaskActionApi|1<<0)
}

func TestName(t *testing.T) {

	f := func(i ...int64) {
		t.Log(len(i))
		t.Log(i)
	}

	f(1)
	f(1, 3, 4)
}
