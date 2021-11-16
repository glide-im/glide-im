## 注意事项

API 处理函数使用 Context.Response 方法需要注意, 该方法必须在返回该请求结果的时候调用, 且必须调用

```go
package api

import (
	"go_im/im/api/apidep"
	"go_im/im/api/groups"
	"go_im/im/api/router"
	"go_im/im/message"
)

func AddGroupMember(ctx *router.Context, req *groups.AddMemberRequest) error {
	//... do something
	notifyGroupMemberAdded := message.NewMessage(1, "", "")
	// 通知更新群成员, 不可使用 Response
	apidep.ClientManager.EnqueueMessage(ctx.Uid, ctx.Device, notifyGroupMemberAdded)
	// 响应当前请求, 必须 Response
	ctx.Response(message.NewMessage(ctx.seq, "success", ""))
	return nil
}

```

TODO

* [ ] 完成各 api 的更新
* [x] http 服务 api
* [x] http 鉴权