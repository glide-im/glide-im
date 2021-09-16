package rpc

import (
	"context"
	"github.com/smallnest/rpcx/share"
)

type Ctx struct {
	context.Context
}

func NewCtx() *Ctx {
	return &Ctx{
		context.Background(),
	}
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
