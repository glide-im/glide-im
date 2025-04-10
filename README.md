## glide

**[立即体验](http://im.dengzii.com/)**

[项目介绍文章](https://juejin.cn/post/7088262766779170853)

**几个较为重要的模块**

- `gate`: 长连接消息网关抽象, 所有消息的入口, 提供管理网关中客户端的接口, 如设置 id, 退出, 推送消息等.
- `messaging`: 消息路由层, 处理来自 `gate` 的消息, 并根据消息类型进行转发给相应的消息处理器.
- `subscription`: 提供适用于群聊, 实时订阅等场景的接口.

**公共消息的定义**

`messages` 包提供了客户端和服务端通讯的消息实体, 类型及消息编解码器. `GlideMessage` 为最基本的公共消息实体. 

**相关项目**

[TypeScript WebApp](https://github.com/glide-im/glide_ts_sdk)

[业务 HTTP API接口](https://github.com/glide-im/api)
