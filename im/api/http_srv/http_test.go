package http_srv

import (
	"bytes"
	"encoding/json"
	"go_im/im/api/auth"
	"go_im/im/dao"
	"go_im/pkg/db"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestRunHTTPServer(t *testing.T) {
	db.Init()
	dao.Init()
	Run("0.0.0.0", 8081)
}

func TestName(t *testing.T) {

	request := auth.RegisterRequest{
		Account:  "c",
		Password: "2",
	}
	j, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}

	cp := CommonParam{
		Uid:    1,
		Device: 1,
		Data:   string(j),
	}
	cj, _ := json.Marshal(cp)

	resp, err := http.Post("http://localhost:8081/api/auth/register", "application/json", bytes.NewBuffer(cj))
	if err != nil {
		panic(err)
	}
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	t.Log(string(all))
}
