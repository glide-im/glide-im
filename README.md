# Glide-IM

![i](_art/logo.png)

[![Go Report Card](https://goreportcard.com/badge/github.com/dengzii/go_im)](https://goreportcard.com/report/github.com/dengzii/go_im)

一款高可靠, 高性能的 IM 服务, 支持分布式集群部署, 单机(单实例)部署.

## 简介

GlideIM 是一款完全开源的 IM 服务, 支持单实例, 微服务等多种方式部署. 支持 WebSocket, TCP 两种连接协议, 内置 JSON,
ProtoBuff 两种消息交换协议, 并支持添加协议解析过程, 消息加密等. GlideIM 还实现了智能心跳保活机制, 死链接检测, 消息 ACK 机制等功能.

## 立即体验

[安卓体验 APP 下载](https://github.com/Glide-IM/Glide-IM-Android/releases)

[GlideIM Web (TODO)](https://github.com/Glide-IM/Glide-IM-Web)

## 一键部署

```shell
wget https://raw.githubusercontent.com/Glide-IM/Glide-IM/master/cmd/script/glide_im/fast_deploy.sh && chmod +x fast_deploy.sh && ./fast_deploy.sh 
```

## 相关仓库

[前端示例项目源码](https://github.com/Glide-IM/im_web)

[Android 端示例项目源码及 SDK](https://github.com/Glide-IM/Glide-IM-Android)

[Java SDK](https://github.com/Glide-IM/Glide-IM-Java-SDK)

## 性能

单机支持约 20w(4万消息吞吐量) 活跃用户同时聊天(100Mbps), [查看测试数据](https://github.com/Glide-IM/Glide-IM/blob/master/doc/performance_test.md)

## 系统架构

介绍文章: [GlideIM - Golang 实现的高性能的分布式 IM](https://github.com/Glide-IM/Glide-IM/blob/master/doc/arch.md)

![i](_art/system_arch.png)

## 讨论群

[![QQ Group 793204140](http://pub.idqqimg.com/wpa/images/group.png)](https://qm.qq.com/cgi-bin/qm/qr?k=PJvSdCQXtAXyBGuOyP-T2CPu9eVNmzls&jump_from=webapi)

## 特别鸣谢

[![JetBrains](https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.svg)]( https://jb.gg/OpenSourceSupport)

## License

参见 [LICENSE](LICENSE)