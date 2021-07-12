package im

import "go_im/im/entity"

type Chain interface {
	Process(message *entity.Message) *entity.Message
}

type Interceptor interface {
	Intercept(message *entity.Message, chain Chain) *entity.Message
}

type interceptorChain struct {
	interceptors []Interceptor
	index        int
}

func NewInterceptorChain() *interceptorChain {
	ic := new(interceptorChain)
	return ic
}

func (c *interceptorChain) Process(message *entity.Message) *entity.Message {
	if len(c.interceptors) < c.index+1 {
		return nil
	}
	interceptor := c.interceptors[c.index]
	ret := interceptor.Intercept(message, c)
	c.index++
	return ret
}

func (c *interceptorChain) Add(interceptor Interceptor) {
	c.interceptors = append(c.interceptors, interceptor)
}

type ChatMessageInterceptor struct {
}

func (c *ChatMessageInterceptor) Intercept(message *entity.Message, chain Chain) *entity.Message {
	return chain.Process(message)
}
