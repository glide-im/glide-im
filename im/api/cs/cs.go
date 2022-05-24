package cs

import route "github.com/glide-im/glideim/im/api/router"

type CsApi struct {
}

func (*CsApi) GetRecentChatMessage(ctx *route.Context) error {

	// TODO 2022-4-26
	ctx.ReturnSuccess(&Waiter{
		Uid:      0,
		Nickname: "CustomerService",
		Avatar:   "",
	})
	return nil
}
