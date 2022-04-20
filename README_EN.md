# Glide-IM

![i](_art/logo.png)

[![Go Report Card](https://goreportcard.com/badge/github.com/dengzii/go_im)](https://goreportcard.com/report/github.com/dengzii/go_im)
[![Reddit](https://img.shields.io/badge/Reddit-Reddit-red)](https://www.reddit.com/r/golang/comments/u7nnwy/my_first_highperformace_distributed_project_an/)

[README - 中文](README.md)

## Summary

Glide is a high-availability, high-performance instant messaging service.

Feel free to post your development in the issue.

## Related Project

Front-End: [GlideIM-Web](https://github.com/Glide-IM/im_web)

Android-App:[GlideIM-Android](https://github.com/Glide-IM/Glide-IM-Android)

## Build and Run

### Env Dependencies

- Singleton Mode
  - Redis
  - MySQL

- Microservice Mode
  - ETCD: `./cmd/script/etcd`
  - NSQ: `./cmd/script/nsq`

### Application Entries

The directory `./cmd/run` is the program entries dir, directory names express the service or deployment mode.

- Microservices
  - api_rpc API services
  - dispatch message delivery service
  - gateway user connection gateway
  - group  the group cha messaging
  - messaging message router service
- HTTP API
  - api_http the http server entry for api
- Singleton 
  - singleton run all features in one instance

### Compile `.proto`

The compiled protobuf message definition files in dir: `/protobuf/gen`.

```shell
/protobuf/compile.sh
```

### Run in Singleton mode

For easy debugging, most cases we just need run project singleton mode, just list debug client, or handle IM core business, the ws server, message storage etc. 

```shell
go run ./cmd/run/singleton/main.go
```

### Executable bin file

If you are feel boredom with build project, or other cases, there is a compiled binary file in singleton mode.

 [Releases](https://github.com/Glide-IM/Glide-IM/releases)

## Performance

Single machine support about 20W (40k message throughput) active users chat simultaneously (100Mbps)  , [View Testing Data](https://github.com/Glide-IM/Glide-IM/blob/master/doc/performance_test.md)

## System architecture

Introduce article: [View Article](https://github.com/Glide-IM/Glide-IM/blob/master/doc/arch-en.md)

![i](_art/system_arch.png)

## Tanks

[![JetBrains](https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.svg)]( https://jb.gg/OpenSourceSupport)

## License

 [LICENSE](LICENSE)