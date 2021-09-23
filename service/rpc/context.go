package rpc

import (
	"context"
	"github.com/smallnest/rpcx/share"
)

type Ctx struct {
	context.Context
}

func NewCtxFrom(c context.Context) *Ctx {
	return &Ctx{c}
}

func NewCtx() *Ctx {
	return NewCtxFrom(context.Background())
}

func (c *Ctx) PutReqExtra(k string, v string) *Ctx {
	mate := c.Context.Value(share.ReqMetaDataKey)
	if mate == nil {
		mate = map[string]string{}
		c.Context = context.WithValue(c.Context, share.ReqMetaDataKey, mate)
	}
	m := c.Context.Value(share.ReqMetaDataKey).(map[string]string)
	m[k] = v
	return c
}

func (c *Ctx) PutResExtra(k string, v string) *Ctx {
	mate := c.Context.Value(share.ResMetaDataKey)
	if mate == nil {
		mate = map[string]string{}
		c.Context = context.WithValue(c.Context, share.ResMetaDataKey, mate)
	}
	m := c.Context.Value(share.ResMetaDataKey).(map[string]string)
	m[k] = v
	return c
}

func (c *Ctx) GetReqExtra(k string) (string, bool) {
	mate := c.Context.Value(share.ReqMetaDataKey)
	if mate == nil {
		return "", false
	}
	m := c.Context.Value(share.ReqMetaDataKey).(map[string]string)
	v, ok := m[k]
	return v, ok
}

func (c *Ctx) GetResExtra(k string) (string, bool) {
	mate := c.Context.Value(share.ResMetaDataKey)
	if mate == nil {
		return "", false
	}
	m := c.Context.Value(share.ResMetaDataKey).(map[string]string)
	v, ok := m[k]
	return v, ok
}
