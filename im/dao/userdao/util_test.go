package userdao

import "testing"

func Test_genToken(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Log(genToken(15))
	}
}
