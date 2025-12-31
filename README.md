## Glide

Glide 是一个轻量、快速、稳定、可拓展的消息服务。定位为单体 IM 服务实现，提供长连接网关、消息路由、群聊/订阅等能力，并通过 RPC 对外暴露管理接口。

### 项目定位与边界

- 提供长连接接入（TCP/WS）、客户端会话管理、消息路由与转发等核心能力。
- 支持群聊/订阅场景，面向业务侧提供 RPC 管理接口。
- 仅负责消息转发与存储等流程相关能力，用户体系、关系、群管理等由配套 HTTP API 服务负责。

### 架构概览

- **消息流转路径**：客户端连接（`pkg/conn`）→ 网关（`pkg/gate`）→ 路由层（`pkg/messaging`）→ 具体处理器（聊天/订阅/心跳等）→ 存储或分发。
- **管理路径**：业务 HTTP API 服务通过 RPC 调用本服务（`im_service/server`），执行踢人、推送、订阅管理等操作。
- **内部事件**：上线/下线通过 `internal.online` / `internal.offline` 触发内部路由与订阅联动。

### 模块详解

- `cmd/im_service`: 服务启动入口与整体装配。
- `pkg/conn`: TCP/WS 连接抽象与服务端实现。
- `pkg/gate`: 长连接网关，负责客户端生命周期管理、鉴权、设置/解绑 id、推送与断连。
- `pkg/messaging`: 消息路由层，按消息类型分发到具体处理器（聊天、心跳、订阅、离线等）。
- `pkg/subscription` + `pkg/subscription/subscription_impl`: 频道/群聊订阅、权限校验与消息分发。
- `pkg/messages`: 公共消息协议、动作定义与编解码器，`GlideMessage` 为基础消息实体。
- `pkg/rpc` + `im_service/client|server`: 对外 RPC 管理接口与客户端封装。
- `pkg/store` + `internal/message_store_db`: 消息存储与持久化实现。
- `internal/world_channel`: 世界频道相关实现（示例场景）。
- `config`: 配置加载与默认配置。
- `pkg/hash`、`pkg/timingwheel`: 一致性哈希与定时轮等基础组件。

### 协议与动作

**消息结构（`GlideMessage`）**

- `ver`: 协议版本。
- `seq`: 消息序号，用于确认/重发等流程。
- `action`: 动作类型（字符串）。
- `from` / `to`: 发送者与接收者标识。
- `data`: 消息体（JSON），接收端在 `Data.Deserialize` 时进行反序列化。
- `ticket` / `sign`: 鉴权与签名相关字段（聊天/群聊消息会校验）。
- `extra`: 扩展字段。

**常用动作分类（`pkg/messages/actions.go`）**

- 系统/心跳：`hello`、`heartbeat`、`notify.*`
- 鉴权：`authenticate`
- 单聊/群聊：`message.chat`、`message.chat.resend`、`message.group`、`message.group.notify`
- ACK：`ack.request`、`ack.message`、`ack.group.msg`、`ack.notify`、`ack.offline`
- 订阅/接口：`api.group.members`、`api.state.sub`、`api.success`、`api.failed`
- 内部：`internal.online`、`internal.offline`

### 连接与鉴权流程

- **连接建立**：客户端连接后会收到 `Action=hello` 的欢迎消息，其中包含临时连接 id。
- **ID 结构**：客户端 id 由 `gateway_uid_device` 组成，临时 id 使用 `tmp@` 前缀。
- **登录鉴权**：客户端发送 `Action=authenticate`，`data` 内包含 `EncryptedCredential`。该凭证由业务服务生成，使用 `CommonConf.SecretKey` 进行 AES-CBC 加密，包含 `user_id/device_id/connection_id/secrets` 等信息。
- **鉴权时效**：凭证时间戳超过约 1500 秒会被拒绝（`credential expired`）。
- **重复登录处理**：若同一 `uid/device` 已在线，会触发踢下线逻辑（`notify.kickout`）。
- **消息鉴权**：对 `message.chat` / `message.group` / `message.chat.resend` 会校验 `ticket`，需携带由业务侧基于 `MessageDeliverSecret` 生成的签名票据，不合法则返回 `notify.forbidden`。

### 消息路由与处理

- `pkg/messaging` 初始化默认 Handler（聊天、群聊、心跳、订阅、内部在线/离线等）。
- 未匹配的动作会返回 `notify.unknown.action`，并记录日志。
- 离线处理由 `messaging/offline_handler.go` 触发，结合存储配置决定是否入库与后续推送。

### 订阅/群聊能力

- RPC 支持 `Subscribe/Unsubscribe/UpdateSubscriber`，以及 `Create/Update/RemoveChannel` 与 `Publish`。
- 频道信息包含 `muted/blocked/closed` 等状态字段，可用于权限与分发控制。
- 群聊消息通过订阅层统一分发，`message.group.notify` 用于通知类事件。

### 存储与离线

- `CommonConf.StoreMessageHistory=true` 时写入 MySQL；`StoreOfflineMessage=true` 时写入 Redis。
- 关闭历史与离线后无需配置 MySQL/Redis。
- `pkg/store` 提供可替换的存储接口，包含 Kafka 生产/消费实现示例。

### 配置与运行

- 配置示例：`config/config.toml`。
- 入口：`cmd/im_service/main.go`。
- 关键配置：
  - `CommonConf`: 消息存储开关与服务秘钥。
  - `WsServer`: WebSocket 监听地址与端口。
  - `IMRpcServer`: RPC 监听地址与端口。
  - `MySql` / `Redis`: 历史与离线消息持久化。
  - `Kafka`: 可选消息流转配置。

### 示例与调用

- RPC 客户端示例：`example/client/rpc_client_example.go`。
- API 调用示例：`example/api/api_example.go`。
- 文档补充：`docs/documents.md` 与 `docs/消息鉴权.md`。

**相关项目**

[Dart SDK](https://github.com/glide-im/glide_dart_sdk)

[Glide Flutter](https://github.com/glide-im/glide-chat)

[业务 HTTP API接口](https://github.com/glide-im/api)
