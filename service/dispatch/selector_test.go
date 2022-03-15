package dispatch

import (
	"testing"
)

func TestReflectMethodName(t *testing.T) {

	s := newSelector()
	nodes := map[string]string{
		"node_a": "",
		"node_c": "",
		"node_d": "",
	}
	s.UpdateServer(nodes)

	for i := 0; i < 10; i++ {
		ctx := contextOfUidHashRoute(int64(i))
		n := s.Select(ctx, "", "", nil)
		t.Log(i, "=>", n)
	}

	nodes["node_e"] = ""
	s.UpdateServer(nodes)
	t.Log("=====================")

	for i := 0; i < 10; i++ {
		ctx := contextOfUidHashRoute(int64(i))
		n := s.Select(ctx, "", "", nil)
		t.Log(i, "=>", n)
	}

	delete(nodes, "node_a")
	s.UpdateServer(nodes)
	t.Log("=====================")

	for i := 0; i < 10; i++ {
		ctx := contextOfUidHashRoute(int64(i))
		n := s.Select(ctx, "", "", nil)
		t.Log(i, "=>", n)
	}
}
