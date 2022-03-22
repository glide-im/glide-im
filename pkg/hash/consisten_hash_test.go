package hash

import (
	"math/rand"
	"strconv"
	"testing"
)

func TestConsistentHash_Add(t *testing.T) {
	c := NewConsistentHash()
	_ = c.Add("A")
	_ = c.Add("B")
	_ = c.Add("C")
	_ = c.Add("D")
	_ = c.Add("E")
	_ = c.Add("F")
	_ = c.Add("G")

	//for _, n := range hash.nodes {
	//	t.Log(n.Val, n.hash, n.virtual)
	//}

	rates := map[string]int{
		"A": 0,
		"B": 0,
		"C": 0,
		"D": 0,
		"E": 0,
		"F": 0,
		"G": 0,
	}

	count := 10000

	for i := 0; i < count; i++ {
		s := strconv.FormatInt(rand.Int63n(100000), 10)
		nd, _ := c.Get(s)
		r := rates[nd.Val]
		rates[nd.Val] = r + 1
	}

	for k, v := range rates {
		t.Log(k, v, int(float64(v)/float64(count)*float64(100)))
	}
}

func TestConsistentHash_Remove(t *testing.T) {
	c := NewConsistentHash()
	_ = c.Add("A")
	_ = c.Add("B")
	_ = c.Add("C")
	_ = c.Add("D")
	_ = c.Add("E")
	_ = c.Add("F")
	//for _, n := range hash.nodes {
	//	t.Log(n.Val, n.hash, n.virtual)
	//}
	e := c.Remove("A")
	if e != nil {
		t.Error(e)
	}
	//t.Log("=====================")
	//for _, n := range hash.nodes {
	//	t.Log(n.Val, n.hash, n.virtual)
	//}
}

func TestAdd(t *testing.T) {
	c := NewConsistentHash2(1)
	_ = c.Add("A")
	_ = c.Add("B")

	for i := 0; i < 5; i++ {
		s := strconv.FormatInt(int64(i), 10)
		n, _ := c.Get(s)
		t.Log(i, ":", n.Val)
	}

	t.Log("===================")
	_ = c.Add("C")
	for i := 0; i < 5; i++ {
		s := strconv.FormatInt(int64(i), 10)
		n, _ := c.Get(s)
		t.Log(i, ":", n.Val)
	}
}
