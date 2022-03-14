package dispatch

import (
	"math/rand"
	"strconv"
	"testing"
)

func TestConsistentHash_Add(t *testing.T) {
	c := NewConsistentHash()
	c.Add("A")
	c.Add("B")
	c.Add("C")
	c.Add("D")
	c.Add("E")
	c.Add("F")
	c.Add("G")

	//for _, n := range c.nodes {
	//	t.Log(n.val, n.hash, n.virtual)
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
		nd := c.Get(s)
		r := rates[nd.val]
		rates[nd.val] = r + 1
	}

	for k, v := range rates {
		t.Log(k, v, int(float64(v)/float64(count)*float64(100)))
	}
}

func TestConsistentHash_Remove(t *testing.T) {
	c := NewConsistentHash()
	c.Add("A")
	c.Add("B")
	c.Add("C")
	c.Add("D")
	c.Add("E")
	c.Add("F")
	//for _, n := range c.nodes {
	//	t.Log(n.val, n.hash, n.virtual)
	//}
	e := c.Remove("A")
	if e != nil {
		t.Error(e)
	}
	//t.Log("=====================")
	//for _, n := range c.nodes {
	//	t.Log(n.val, n.hash, n.virtual)
	//}
}
