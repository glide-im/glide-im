package hash

import (
	"errors"
	"fmt"
	"sync"
)

const (
	duplicateVirtual = 100 // 1_000_000
	seed             = 0xabcd1234
)

var (
	errNodeExist = errors.New("node already exist")
)

type Node struct {
	Val     string
	hash    uint32
	virtual bool
	real    *Node
}

type Nodes struct {
	nd      Node
	virtual []Node
	hit     int64
}

func (n *Nodes) appendVirtual(node Node) {
	n.virtual = append(n.virtual, node)
}

type ConsistentHash struct {
	nodes   []Node
	nodeMap map[string]*Nodes
	virtual int

	mu sync.RWMutex
}

func NewConsistentHash() *ConsistentHash {
	return NewConsistentHash2(duplicateVirtual)
}

func NewConsistentHash2(virtual int) *ConsistentHash {
	hash := &ConsistentHash{
		nodes:   []Node{},
		nodeMap: map[string]*Nodes{},
		virtual: virtual,
		mu:      sync.RWMutex{},
	}
	return hash
}

// Remove node by id, include virtual node.
func (c *ConsistentHash) Remove(id string) error {
	nodes, ok := c.nodeMap[id]
	if !ok {
		return errors.New("node does not exist, id:" + id)
	}
	for _, vNd := range nodes.virtual {
		ndIndex, exist := c.findIndex(vNd.hash)
		if exist {
			ndIndex--
		} else {
			return errors.New("virtual node does not exist, id:" + vNd.Val)
		}
		c.mu.RLock()
		nd := c.nodes[ndIndex]
		c.mu.RUnlock()
		if nd.hash != vNd.hash {
			return errors.New("could not find virtual node, id:" + vNd.Val)
		} else {
			c.removeIndex(ndIndex)
		}
	}
	index, exist := c.findIndex(nodes.nd.hash)
	if !exist {
		return errors.New("real node not fund")
	}
	index--
	c.removeIndex(index)
	delete(c.nodeMap, id)
	return nil
}

func (c *ConsistentHash) Get(data string) (*Node, error) {
	hash := Hash([]byte(data), seed)
	index, _ := c.findIndex(hash)
	return c.get(index)
}

func (c *ConsistentHash) Add(id string) error {
	_, ok := c.nodeMap[id]
	if ok {
		return errors.New("node already exist, id=" + id)
	}
	hash := Hash([]byte(id), seed)
	nd := Node{
		Val:     id,
		hash:    hash,
		virtual: false,
		real:    nil,
	}
	c.nodeMap[id] = &Nodes{
		nd:      nd,
		virtual: []Node{},
	}
	c.addNode(nd)
	c.addVirtual(&nd, c.virtual)
	return nil
}

func (c *ConsistentHash) get(index int) (*Node, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.nodes) == 0 {
		return nil, errNodeExist
	}
	if index == len(c.nodes) {
		index = len(c.nodes) - 1
	}
	n := c.nodes[index]
	if n.virtual {
		return n.real, nil
	}
	return &n, nil
}

func (c *ConsistentHash) addVirtual(real *Node, duplicate int) {
	for i := 0; i < duplicate; i++ {
		vNodeID := fmt.Sprintf("%s_#%d", real.Val, i)
		hash := Hash([]byte(vNodeID), seed)
		vNode := Node{
			Val:     vNodeID,
			hash:    hash,
			virtual: true,
			real:    real,
		}
		c.addNode(vNode)
		nds := c.nodeMap[real.Val]
		nds.appendVirtual(vNode)
	}
}

func (c *ConsistentHash) addNode(nd Node) {

	index, _ := c.findIndex(nd.hash)

	c.mu.Lock()
	defer c.mu.Unlock()

	p1 := c.nodes[:index]
	p2 := c.nodes[index:]
	n := make([]Node, len(p1))
	copy(n, p1)
	n = append(n, nd)
	for _, i := range p2 {
		n = append(n, i)
	}
	c.nodes = n
}

func (c *ConsistentHash) removeIndex(index int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if index == len(c.nodes)-1 {
		c.nodes = c.nodes[:len(c.nodes)-1]
		return
	}

	p2 := c.nodes[index+1:]
	c.nodes = c.nodes[:index]
	for _, n := range p2 {
		c.nodes = append(c.nodes, n)
	}
}

func (c *ConsistentHash) findIndex(h uint32) (int, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	left := 0
	right := len(c.nodes)
	exist := false

LOOP:
	if left < right {
		middle := (left + right) / 2
		hash := c.nodes[middle].hash
		if hash < h {
			left = middle + 1
		} else if hash == h {
			left = middle + 1
			exist = true
		} else {
			right = middle
		}
		goto LOOP
	}
	return left, exist
}
