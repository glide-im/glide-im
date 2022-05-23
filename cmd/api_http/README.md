# API HTTP 服务入口

提供非主动推送部分的 IM 增删改查 API 接口

API HTTP 依赖 IM 服务, 例如添加好友, 邀请群成员, 需要通过 IM 服务主动推送消息

可以选择两种方式访问 IM 服务, ETCD 服务发现或者 地址+端口 直连, 如下 im 服务配置示例

```toml
[ApiHttp.IMService]
Addr = "0.0.0.0"
Port = 8080
# Etcd = []
# Name = "" 
```

如果两种方式都配置了, 则优先使用 ETCD 服务发现