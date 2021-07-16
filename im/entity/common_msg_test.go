package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestActionValue(t *testing.T) {

	t.Logf("\n%b\n%b\n%b", MaskActionApi, ActionUserLogin, RespActionFailed)
	assert.Equal(t, ActionUserLogin, MaskActionApi|1<<0)
}
