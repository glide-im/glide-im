package dispatch

import (
	"testing"
)

func TestReflectMethodName(t *testing.T) {

	s := newSelector()
	nodes := map[string]string{
		"node_AAA": "",
		"node_BBB": "",
		"node_CCC": "",
	}
	s.UpdateServer(nodes)

	for i := 0; i < 10; i++ {
		ctx := contextOfUidHashRoute(int64(i))
		n := s.Select(ctx, "", "", nil)
		t.Log(i, "=>", n)
	}

	nodes["node_e"] = ""
	nodes["node_f"] = ""
	s.UpdateServer(nodes)
	t.Log("=====================")

	for i := 0; i < 10; i++ {
		ctx := contextOfUidHashRoute(int64(i))
		n := s.Select(ctx, "", "", nil)
		t.Log(i, "=>", n)
	}

	delete(nodes, "node_AAA")
	s.UpdateServer(nodes)
	t.Log("=====================")

	for i := 0; i < 10; i++ {
		ctx := contextOfUidHashRoute(int64(i))
		n := s.Select(ctx, "", "", nil)
		t.Log(i, "=>", n)
	}
}
