# 公共

## 请求

客户端所有请求需要包含以下内容, 请求具体接口的数据放在 Data 中, 通过 Action 区别不同接口.

- 定义

| 字段 | Ver      | Seq              | Action    | Data                               |
| ------ | ---------- | ------------------ | ----------- | ------------------------------------ |
| 类型 | number   | number           | string    | object                             |
| 描述 | 接口版本 | 客户端消息序列号 | 动作      | 具体动作的数据                     |
| 例子 | 1        | 1                | api.login | {"Account":"abc","Password":"abc"} |

- 例子

```json
{
  "Ver": 1,
  "Seq": 1,
  "Action": "api.login",
  "Data": {}
}
```

## 响应

以下内容每次请求都会返回, 具体响应数据放在 Data 中, 通过 Action 确定请求是否成功, failed 表示失败, success 表示成功

- 定义

| 字段 | Ver        | Seq                  | Action           | Data                           |
| ------ | ------------ | ---------------------- | ------------------ | -------------------------------- |
| 类型 | number     | number               | string           | object                         |
| 描述 | 消息版本号 | 对于客户端请求序列号 | 结果             | 响应结果                       |
| 例子 | 1          | 1                    | success / failed | {"Token":"asdfjwe","Uid": 123} |

- 例子

```json
{
  "Ver": 1,
  "Seq": 1,
  "Action": "success",
  "Data": {
    "Uid": 123,
    "Token": "JKLjflkqwfjfad"
  }
}
```

# 鉴权

## 登录

### 请求 (api.login)

- 定义

| 字段 | Device   | Account | Password |
| ------ | ---------- | --------- | ---------- |
| 类型 | number   | string  | string   |
| 描述 | 设备标识 | 账号    | 密码     |

- 例子

```json
{
  "Ver": 1,
  "Seq": 1,
  "Action": "api.login",
  "Data": {
    "Device": 1,
    "Account": "abc",
    "Password": "abc"
  }
}
```

### 响应

- 定义

| 字段 | Token    | Uid      |
| ------ | ---------- | ---------- |
| 类型 | string   | number   |
| 描述 | 登录凭证 | 用户标识 |

- 例子

```json
{
  "Ver": 1,
  "Seq": 1,
  "Action": "success",
  "Data": {
    "Uid": 1234,
    "Token": "JOQIfJKANFJKAq"
  }
}
```

## 注册

### 请求 (api.register)

- 定义

| 字段 | Account | Password |
| ------ | --------- | ---------- |
| 类型 | string  | string   |
| 描述 | 账号    | 密码     |

- 例子

> 省略了公共响应内容, 后续例子也如此

```json
{
  "Account": "abc",
  "Password": "abc"
}
```

### 响应

- 定义

无响应体

- 例子

```json
{
  "Ver": 1,
  "Seq": 1,
  "Action": "success",
  "Data": null
}
```

## 退出登录

### 请求 (api.logout)

无请求体

### 响应

无响应体

# 消息

## 客户端获取消息 ID

客户端每次发送消息前需要获取一个消息 ID

> 暂无, 客户端随机生成确保唯一即可

## 客户端发送消息

### 三种不同发送消息的 Action

| Action | message.chat | message.chat.retry | message.chat.resend |
| -------- | -------------- | -------------------- | --------------------- |
| 描述   | 首次发送这条消息         | 重试(没有收到服务端确认)              | 重试(没有收到接收方确认)                |

### 请求体

- 定义

| 字段 | mid    | from      | to       | c_seq     | type     | content | c_time   |
| ------ | -------- | ----------- | ---------- | ----------- | ---------- | --------- | ---------- |
| 类型 | int64  | int64     | int64    | int       | int      | string  | int64    |
| 说明 | 消息id | 发送者uid | 接收者id | 客户端seq | 消息类型 | 内容    | 发送时间 |

- 例子

```json
{
  "Ver": 0,
  "Seq": 1,
  "Action": "message.chat",
  "Data": {
    "Mid": 1001,
    "CSeq": 1,
    "From_": 1234,
    "To": 5678,
    "Type": 1,
    "Content": "HelloWorld",
    "CTime": 1637472935
  }
}
```

## 服务端 ACK 消息 (ack.message)

服务端收到客户端发送的消息后, 回复客户端一条 ACK 消息, 表示服务端已收到

- 定义

| 字段 | Mid    |
| ------ | -------- |
| 类型 | number |
| 描述 | 消息ID |

- 例子

```json
{
  "Ver": 0,
  "Seq": 1,
  "Action": "ack.message",
  "Data": {
    "Mid": 1234
  }
}
```
> 客户端收到这条消息时, 则确认该消息发送已成功一半

## 客户端 ACK 消息 (ack.request)

客户端收到一条消息后, 需要回复服务端一条 ack.request 消息, 表示我已收到这条消息

- 定义

| 字段 | Mid    | From     |
| ------ | -------- | ---------- |
| 类型 | number | number   |
| 描述 | 消息ID | 发送者ID |

- 例子

```json
{
  "Ver": 0,
  "Seq": 1,
  "Action": "ack.request",
  "Data": {
    "Mid": 1234,
    "From": 333
  }
}
```
> 客户端收到任何一条聊天消息后, 都需要向服务端发送一条确认收到消息

## 服务端确认收到消息 (ack.notify)

服务端收到一条消息接收方的 ack.request 消息后, 将确认消息发送给消息发送者方, 表示接收方已收到消息

- 定义

| 字段 | Mid    | From     |
| ------ | -------- | ---------- |
| 类型 | number | number   |
| 描述 | 消息ID | 发送者ID |

- 例子

```json
{
  "Ver": 0,
  "Seq": 1,
  "Action": "ack.notify",
  "Data": {
    "Mid": 1234,
    "From": 333
  }
}
```

> 客户端发送一条消息后, 在收到 ack.notify 之前都认为接受方没有收到该消息, 需要按策略重发

### 客户端收到新消息 (message.chat)

- 定义

| 字段 | mid |from      | to       |  c_seq     |type     | content | c_time     |
| ------ |------ | ----------- | ------------ |  ----------- | ---------- | --------- | ------------ |
| 类型 |  int64|int64     | int64                  | int      | int      | string  | int64      |
| 说明 | 消息id| 发送者uid | 接收者id | 客户端seq  | 消息类型 | 内容    | 发送时间 |

- 例子

```json
{
  "Mid": 5125345,
  "CSeq": 2,
  "From": 123,
  "To": 456,
  "Type": 1,
  "Content": "Hello World",
  "CTime": 1637474399
}
```

> 客户端可能会收到相同 Mid 的消息(消息重发), 相同的则不处理