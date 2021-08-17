package im

type Int64Set struct {
	m map[int64]interface{}
}

func NewInt64Set() *Int64Set {
	return &Int64Set{m: make(map[int64]interface{})}
}

func (i *Int64Set) Add(v int64) {
	if i.Contain(v) {
		return
	}
	i.m[v] = nil
}

func (i *Int64Set) Remove(v int64) {
	_, ok := i.m[v]
	if ok {
		delete(i.m, v)
	}
}

func (i *Int64Set) Clear() {
	i.m = make(map[int64]interface{})
}

func (i *Int64Set) ForEach(it func(value int64)) {
	for k := range i.m {
		it(k)
	}
}

func (i *Int64Set) Size() int {
	return len(i.m)
}

func (i *Int64Set) Contain(v int64) bool {
	_, ok := i.m[v]
	return ok
}
