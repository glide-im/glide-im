# GlideIM - High-performance distributed IM implemented in Golang

GlideIM is a **completely open source**, high-performance distributed IM service implemented by Golang, with Android APP examples, JAVA SDK, and Web-side examples, and it is continuously updated and iterated.

GlideIM supports single instance or distributed deployment. It supports WebSocket, TCP two connection protocols, built-in JSON, ProtoBuff two message exchange protocols, and supports adding other protocols, message encryption, etc. It also implements intelligent heartbeat keep-alive mechanism, dead connection detection , message delivery mechanism and other functions.

This project started in mid-2020, and all three children projects have been developed by me alone until now. It is also the first project that I spent the most energy on developing and learning. Since I have a lot of free time for work, I basically work in I developed this project, developed while learning, and consulted a lot of information. After a lot of thinking, the goal is a high-performance distributed IM. The first version of the microservice architecture was basically completed after three months of development, and the later is to do After adjusting the microservice architecture and further optimizing the details of the IM service, as of now (March 3, 2022), there are only some finishing work left for the second version I expected. this project will continue to update and iterate.

- Server source code: [ GlideIM - GitHub]( https://github.com/Glide-IM/Glide-IM)
- Android : [Glide-IM-Android - GitHub]( https://github.com/Glide-IM/Glide-IM-Android)
- Web: [Glide-IM-Web - GitHub]( https://github.com/Glide-IM/Glide-IM-Web)

## 1. Features

### 1.1 User Side Features

- No sense of disconnection and reconnection, message synchronization
- Login to register and stay logged in
- One-on-one chat, group chat
- message roaming, history
- Offline messages
- Log in with multiple devices, squeeze and go offline with the same device
- Multiple types of messages (picture, voice, etc., defined by the client)
- message retraction
- Contact management, group management

### 1.2 Development Side Features

- Support WebSocket, TCP, custom connection protocol
- Support JSON or Protobuff or custom data exchange protocol
- Support distributed deployment.
- Heartbeat keep alive, timeout disconnection, clean up dead links
- Message buffering, asynchronous processing, weak network optimization
- Message delivery mechanism, message resend, ACK
- Message deduplication, order guarantee, read diffusion

## 2. Project Architecture

In order to improve availability and overall stability , due to the limitation of stand-alone performance, a distributed architecture must be used, and the micro-service mode is convenient for maintenance. This project splits six core main modules (services) for the IM business part, and each service can be extended in any number horizontally, and the entire system can have a certain scalability. Each module is divided according to its business characteristics, and the logic and interface are separated. While ensuring the simplicity of the interface, it also has sufficient scalability.

![i](https://github.com/Glide-IM/Glide-IM/blob/master/_art/system_arch.png?raw=true)

### 2.1 The IM Services

- Gateway

The Gateway service is a chat service gateway that manages user connections. All user messages upstream and downstream are handled by this service. Gateway manages user connections, receives and parses messages, sends messages, determines whether the connection is alive, identifies user connections, and disconnects user connections.

The Gateway relies on the Messaging service, and the received client messages will be handled by Messaging. The Gateway provides three interfaces for logging in, logging out, and sending messages with the specified uid.

- Messaging

Messaging is responsible for the routing of different types of messages, such as group messages , API messages, and also processes some types of messages, such as ACK messages, single-chat messages , and heartbeat messages. According to the message type, it is forwarded to Dispatch (user messages), API (API messages) , Group ( group message ).

Messaging relies on API, GroupMessaging , Dispatch, to provide a message routing interface.

- Group

(GroupMessaging) group chat service, is mainly responsible for the delivery, storage, group message confirmation, and member management of multi-person chat messages. After users go online, they initialize the group chat list when synchronizing contacts, and enter the group chat according to the service where the group chat is located .

The Group service relies on the Dispatch service to provide three interfaces of group update, member update, and sending messages to the specified group.

- API

Currently, as the login authentication of persistent connections , all HTTP API interfaces can be accessed through persistent connection messages, which can be flexibly configured according to the specific situation, only need to configure relevant routes.

API relies on Dispatch, Group services. Provides an interface for handling API type messages.

- Dispatch

message routing service is used to route messages to the gateway where the user is located. When the user logs in, Dispatch is notified to update the gateway corresponding to the cached user. The cached information is stored in a fixed service through a consistent hash. When the message is sent, the cache is queried. According to the query The received gateway information is put into the NSQ queue, and each Gateway subscribes to its own message. The message here is not necessarily a user message, but also a control message that notifies the gateway to update the user's state, such as login and logout. Since the message queue is used for communication, So it is called a message, which is actually an interface for calling Gateway.

Dispatch does not directly depend on any service. Messages are sent to the Gateway through NSQ, providing two interfaces for updating routes and delivering messages.

- Broker

user gateway, group chat is a stateful service, and messages need to be routed to the path of the group accurately and quickly . The functions of Broker and Dispatch are roughly similar.

- NSQ

[NSQ](https://nsq.io/) is a message queue implemented by Golang, and all messages are routed through NSQ. Reasons for choosing NSQ compared to other MQs: decentralized distribution (direct connection between production and consumption), low latency, No ordering, high performance, simple binary protocol.

Deploy one on each producer nsqd .

### 2.2 HTTP API Service

- **atuh**: user authentication, login registration, etc.
- **user**: User information management
- **message**: message synchronization, pull and other interfaces, message ID
- **group**: group additions, deletions, revisions
- **other**: ...

The above division is just the division of modules in the project , not independent services , but it is easy to split them. The above interfaces can be accessed through HTTP services or long-term connections. The difference between long-term connections and HTTP access is that : 1. Long connection access needs to add `public request body`, 2. The definition of public response body is different.

### 2.3 Message Routing

#### Gateway message routing

In a distributed deployment environment , the gateway may deploy any number of instances, and users may connect to any of them. When we need to send a message to a user or disconnect a user, we need to find the user's location. Gateway, which needs to record the gateway where all online users are located. You may think of using Redis, or Redis coordinate with secondary cache, but the message throughput of the IM system is very large, and there are other reasons such as diffusion, Redis can easily become a performance bottleneck.

GlideIM uses a consistent hashing algorithm to distribute the gateway information of each user connection to different Dispatch services according to UID, so as to achieve decentralized caching, load balancing, and improve availability.

![i](https://github.com/Glide-IM/Glide-IM/blob/master/_art/msg_routing.png?raw=true)

As shown in the figure above , Dispatch assumes the role of message distribution in the whole process. When a Dispatch service goes down , all cached gateway information in the service will be lost. According to the consistent hash algorithm, the original request will be redirected to After the next Dispatch, obviously this service does not have the cached information of the service that has been down . We can cache the login information in Redis after the user logs in (for other scenarios such as querying the login device list), there is no information in the Dispatch memory. After finding it, go to Redis to look it up, and it can be cached in memory after one lookup. The newly online users will not be affected, but only the users whose Dispatch service where the gateway information is saved during the online period is down .

#### Group message routing

Groups and users need the same route for group services. Different groups may be distributed on different services, but groups do not switch services at will. Generally, only when the service where the group is located is abnormal, reloading group information can be reassigned. Broker reads more and writes less, so it is only necessary to cache the services of all groups in all brokers. Group messages can be forwarded through any broker. When the group service is offline and the group is loaded, all the brokers can be notified to update.

### 2.4 My Design Guidelines

In this project, I have been exploring and pursuing the solution of the contradiction between the cleanliness of the project and the complex system architecture, but since I do not have any relevant experience, there may be some non-mainstream styles or wrong redundant designs in the project, If you have any doubts or suggestions about this, please feel free to post your thoughts in issues.

In the face of a huge and complex problem , it is often convenient to divide and abstract the problem into objects , and divide it into small problems. OOP has always been an advantage in this regard, and this project uses more OOP ideas.

#### Using the interface

For some key points that may need to be extended , GlideIM uses the interface method to implement, for example, the Dao (data access object) layer is implemented as an interface. Aspects are all simple implementations, and optimizations such as performance are not considered, but the interface does not affect the later upgrade and replacement of this aspect.

#### Package by business

There are generally two ways to divide packages, package by layer and package by business. layer packaging means placing the same type of modules in the same package. For example, all API request processing logics(controller) are placed in one layer, and requests and responses The entity objects are placed in the same package. This method adds inconvenience to the management of modules. We need to modify the code in several different packages to change an interface. In this way, the different codes in a package are not related in any way. Similar to how librarians sort books by publisher.

However, subcontracting by function will not update a function or modify the code of different packages. The code related to a function is all in the same package. Each package of GlideIM follows the principle of single responsibility as much as possible, and minimizes the difference between multiple packages. Coupling.

#### Module interface dependencies

The division of project packages is basically the same as the module division of microservices . When this project was started at the beginning, it followed the method of dividing modules by function, and each module only exposed the interface provided to other modules, and the calls between modules only Through a specific interface, this design is convenient for maintenance and for others to view the code, it will not make the calls between modules messy, and it also provides great convenience for later microservices.

For example the interface of the Gateway module

```go
// interface provided to other modules
type Interface interface {
    ClientSignIn ( tempId int64, uid int64, device int64)
    ClientLogout ( uid int64, device int64)
    EnqueueMessage ( uid int64, device int64, message * message.Message )
}
// Gateway depends on the interface
type MessageHandleFunc func ( from int64, device int64, message * message.Message )
```

Other modules only need to know that the Gateway module provides client login and logout and adds messages to the queue of the client with the specified ID. Other modules do not need to know the specific implementation of the Gateway. We can easily replace this module with a microservice service.

### 2.5 Microservices

GlideIM uses [RPCX](https://rpcx.io/) as the basis for microservices. The out-of-the - box microservice solution made me choose it. RPCX has rich functions, superior performance, integrated service discovery, and multiple routes. Scheme, and failure mode, service discovery using [ETCD](http://www.etcd.cn/).

Inter-service communication is used The Protobuff + RPC method is the best combination in terms of performance.

### 2.6 Project Introduction

#### project directory

The following directory structure omits some unimportant directories

```text
├─ cmd // entry (run from here)
│ ├─ performance_test // Performance test code entry
│ ├─run // project program entry
│ │ ├─ api_http // HTTP service of API interface, provided for client access
│ │ ├─ api_rpc // RPC service of API interface, available to other services
│ │ ├─broker // Group routing service
│ │ ├─dispatch // Gateway routing service
│ │ ├─getaway // Gateway service
│ │ ├─group // group service
│ │ ├─messaging // IM message routing service
│ │ └─singleton // Single instance operation (both IM and HTTP API interfaces are started here)
│ └─script // Deployment script
│ ├─ etcd // Start etcd script
│ ├─ glide_im // project deployment script
│ └─ nsq // nsq deployment script
├─ config // Configuration entry
├─ doc // project documentation
├─ im // IM core logic entry
│ ├─ api ////////////// API interface
│ │ ├─ apidep // API external dependencies
│ │ ├─auth // Login authentication
│ │ ├─comm // public
│ │ ├─groups // Group management
│ │ ├─ http_srv // api http service startup logic
│ │ ├─msg // message
│ │ ├─router // Routing abstraction for accessing interfaces through persistent connections
│ │ └─user // User related
│ ├─client ///////////// User connection management related
│ ├─conn // Long connection basic abstraction
│ ├─ dao // data access layer, database related
│ ├─group // Group chat , and group chat messages
│ ├─message // IM message definition
│ ├─messaging // IM message routing
│ └─statistics // data statistics, for testing
├─ pkg ///////////// package, public dependency management
│ ├─ db // database
│ ├─hash // hash algorithm implementation
│ ├─logger // log print
│ ├─ lru // lru cache implementation
│ ├─ mq_nsq // nsq package
│ ├─ rpc // rpc package, based on rpcx
│ └─ timingwheel // Timer, time wheel algorithm implementation
├─ protobuf //////////// protobuf message definition
│ ├─gen // Compiled file
│ ├─ im // im message definition
│ └─ rpc // rpc communication message definition
├─ service /////////// Microservice
│ ├─ api_service // api Microservice implementation
│ ├─broker // Group routing broker service
│ ├─dispatch // Gateway message routing service
│ ├─gateway // Gateway service implementation
│ ├─ group_messaging // group service
│ ├─ messaging_service // im message routing service
├─ sql // Database table structure SQL
```

#### Project viewing guide

The core logic of IM is in the root directory ` im `. Except for microservices , the main business logic of IM is implemented in this directory. The division of packages under ` im ` can be roughly regarded as the division of the following microservices, ` im /conn The ` package contains the logic of the long connection server startup and connection object interface. The new connection will be managed by the ` im /client` package. This package is roughly used to manage the connection, read and write and parse the message, and the message received from the link will be Handed over to ` im /messaging` package, this package tools message types to different modules, for example: authentication message is processed by ` api` module, group message is distributed to `group`, single chat message is under ` im /client` There is basically an ` interface.go ` file under the package involving business logic under ` im ` , which defines the dependencies of this package and the interface provided to the outside world.

microservices is under the root directory `service` package, and the service logic related to IM core business is under ` im` . The default implementation of the interface is replaced. For example, the ` messaging_service` service (message routing) replaces the implementation of the `Interface` interface under the ` im /messaging` package with the `Server` of `messaging_service`, and the other `messaging` depends on The package is the `Client` under the corresponding service. Check the `run.go` under each service package to find the service startup code, which includes its dependencies and implementation settings.

For this mode of stripping core logic and microservices, it is because the division of microservices was finalized at the beginning ( also after several major changes), or for the convenience of implementing and clarifying part of the core logic of IM, so Put the two in two different places, but in the process of my practice, the changes to the service and the changes to the core logic (except the module interface changes) under `im` are not related to each other, which gives me a great deal For the convenience, I still have a certain degree of freedom to add services, such as adding a `broker` in the middle of the `group` service, or special processing can be performed for a method in the interface of a module in a specific ` im` .

But the convenience mentioned above is only because I was not familiar with microservices at the time, and the project development process changed a lot, and it was slightly more convenient to separate the changes. The price of this is that the separation of the two is not conducive to code viewing, which is actually a ` im` defines interfaces and default implementations, and the upper layer ` service` defines the implementation of interfaces, but the dependencies between these interfaces are still to be found in ` im` , and these two packages may be considered later. combined together.

#### Dependencies

- [ BurntSushi / toml ](https://github.com/BurntSushi/toml): This is an excellent configuration file format, which I personally prefer
- [gin](https://github.com/gin-gonic/gin): Excellent HTTP Web Framework
- [ protobuff ](https://github.com/gogo/protobuf): Google's binary data transfer protocol
- [gorilla/ websocket ](https://github.com/gorilla/websocket): The most star WebSocket library in Golang
- [ nsq ](https://github.com/nsqio/go-nsq): simple, high performance, distributed MQ
- [ rpcx ](https://github.com/smallnest/rpcx): High-performance, feature-rich microservices framework
- [gorm](gorm.io/gorm): ORM
- [go- redis / redis ](https://github.com/go-redis/redis): Redis client
- [ants](github.com/panjf2000/ants/v2): coroutine pool

#### Build and run

` cmd /run` in the root directory is the program entry, and the package name indicates its service/mode. For example, the ` api_http ` package is the interface for starting the API with the HTTP service (provided for the client to call in HTTP mode), and the ` api_rpc ` is the Start the API RPC service (provided for other service calls), for quick debugging, there is also a single instance mode `singleton`, this entry will start the IM long connection service and the HTTP API service at the same time, which is convenient for debugging the core logic of ` im `, or Debug client.

**Environment Dependent**

- single instance
  - redis
  - mysql
- Microservices
  - nsq
  - etcd _
  - contains all dependencies of a single instance

**Config Files**

Modify related configuration in `config.toml` in `singleton` package for single instance mode.

In the microservice mode, you need to copy the `service/ example_config.toml` file to the service entry, and modify the relevant configuration according to the environment.

If you cannot run the code in the IDE due to dependencies or other reasons , you can download the compiled executable in `singleton` mode at [here](https://github.com/Glide-IM/GlideIM/releases).

### 2.7 Existing Issues

#### Protocol related

- Client protocol selection

Currently, in order to facilitate the use of the client json for communication is also because the browser's support for binary protocols is not friendly, but the back-end implements both, and the dynamic selection of the protocol is not implemented, or to distinguish between websocket and tcp gateway, distinguish the two , browser and mobile There are differences in the end environment , but no processing is done.

- Protobuf and Json compatibility

Microservice usage protobuf protocol, while client messages may use json protocol, and some use the same struct compiled and generated by ` protoc` , there are some compatibility problems, which have not been dealt with yet.

- Message decoding performance

using json.Marshal in go is extremely poor, and it takes up a lot of time in the entire message flow process. It has not been optimized yet. There are currently a variety of third-party solutions to choose from.

- Protocol version

The message protocol may be upgraded , and new and old clients use different protocol versions for compatibility.

#### Database related

- Message ID generation

Redis Incr is currently used to generate incremental IDs, and there are performance problems. Later, solutions such as Leaf can be considered.

- Database, table, query optimization

No optimization has been made in this regard , and the dao layer simply implements the CRUD function.

#### Microservice related

- Configuration management

At present, the startup service is started by loading the configuration from the local configuration file , which is inconvenient to deploy and manage. Consider using the configuration center later.

- Package structure

At present , the microservice and IM logic are layered design, which is inconvenient to maintain and view, and needs to be adjusted later.

#### Group chat related

- Messages storm

problem of proliferation of group messages , especially if a member of a large group sends messages thousands of times, and the number of groups is slightly larger or the messages are slightly more frequent, the number of messages is very easy to get out of control. In the microservice mode, the user messages of the system gateway need to be processed. Merge and pack into one message.

- Group chat service failure

Group chat is a stateful service . How to not affect the distribution of group messages when the service where the group is located goes down has not yet been dealt with. For recovery, a group chat service monitoring service can be added, and if a group chat service is monitored, it will be disconnected immediately. recovery, but this is for failback only and not for failover.

## 3. Design details

because GlideIM has a lot of design details, limited by time and space, only a few more important links are listed here for a brief explanation. For details or other unexplained places, you can view the source code or join the discussion group to discuss together.

### 3.1 Message Types

The IM message type refers to the chat protocol message type. Note that the message type of the chat content is distinguished. The IM message type is related to the business logic of both front and back ends. The chat message type only needs to be processed in the front end.

- IM message type
  - Chat messages: resend, retry, withdraw, group chat , single chat
  - ACK message: server acknowledges receipt, recipient acknowledges notification, recipient acknowledges delivery
  - Heartbeat message: client heartbeat, server heartbeat
  - API messages: token authentication, logout, etc.
  - Notification messages: new contact, kickout, multi-login, etc.
- Chat content message type
  - picture
  - text

The message type is defined and processed by the client . Binary messages such as voice and pictures are uploaded to the server and then sent to the URL. For example, interactive messages such as red envelopes can also be used in this way.

Defining the agreed message type between clients not only facilitates the later addition of message types , but also facilitates back-end maintenance. The back-end does not need to know the type of message content.

### 3.2 Delivery Mechanism

Although TCP is a reliable transmission, it is not foolproof in the process of message delivery. For example, if the receiver is online but the network is not good, and the client is reconnecting, the message may not be delivered to the recipient in time. Design a delivery mechanism to improve delivery. The rate is necessary. GlideIM adopts the delivery mechanism of two confirmations (server and receiver) for a single message.

A sends a message to B. If B is online, the server will reply to A. A server confirms receipt of the message to tell A that the server has received it, and it has been stored in the warehouse. Then the server sends it to B, and B sends a message when it receives it. Confirm the receipt of the message to the server, and the server will send a confirmation delivery message to A after receiving it. If B is not online, it will directly send it to A to confirm the delivery. At this time, A knows that the message must be received by B, In this process, if A does not receive confirmation from B , it will retry multiple times.

Doing retransmission on the client side avoids the complexity of the server-side logic , while doing it on the client-side greatly simplifies the logic. With this delivery mechanism, whether the message is lost or the network environment is bad, the delivery rate can be improved And the message timeliness rate, and if we use UDP, this mechanism can also ensure the delivery rate of the message.

GlideIM message ACK mechanism in the case of several message loss situations is as follows:

![ Message ACK](https://github.com/Glide-IM/Glide-IM/blob/master/doc/img/message_ack.png?raw=true)

The above is only GlideIM is part of the guarantee of delivery, in fact, there are other measures such as the receiving end judging whether the message is missing according to Seq.

### 3.3 Message Routing

In a distributed environment, when users connect to different gateways, we must know the gateway where the user is located in order to deliver messages accurately. Therefore, the gateway information must be cached when the user connects to the gateway, so that the message can be delivered accurately. Delivered to the gateway where the user is located.

GlideIM is to cache the gateway information in Redis when the user logs in. Of course, we cannot query Redis every time we send a message. We can cache a copy in memory, and query Redis if there is no memory. After the user logs in to the chat service, according to the contact list In different gateways, notify the gateway to update its own routing information.

In order to prevent the client from switching the gateway frequently , the client access gateway is returned by the server when the user logs in. This can also be considered for formulating load balancing strategies, multi-terminal login connections to the same gateway , etc.

### 3.4 Keep Alive Mechanism

Although TCP has KeepAlive , its default time is too long, and it cannot judge whether the service is available or whether the client is available. KeepAlive only ensures that the TCP is in a connected state. If an error occurs at the application layer, the server connection is still smooth.

Heartbeat is divided into server heartbeat and client heartbeat . In most cases, the network on the client side is not smooth, such as entering the elevator, the screen is in power saving mode, and the server generally only has a network failure when it is down . Therefore, GlideIM performs active heartbeat on the client. If the server does not receive the heartbeat within the specified time, the server starts to send the heartbeat packet, and the connection is disconnected if the client does not receive the heartbeat for the specified number of times.

heartbeat of the client is 30s, but it does not send a heartbeat every 30s, but if the client does not actively send any message within 30s, a heartbeat is performed, and the server also judges according to this rule, which reduces the number of heartbeats and also to ensure survival.

### 3.5 Message Protocol

The message protocol needs to consider the size of the encoded message , readability, encoding speed, supported languages, etc. You can choose binary protocol and text protocol, binary such as Protobuff , or custom, text protocol such as JSON, XML.

GlideIM implements both Protobuff and JSON message protocols. The client can freely choose which protocol to use. The test results show that Protobuff is at least 10 times faster than JSON . When using JSON, message parsing takes up a large part of the entire process. time.

### 3.6 Message deduplication and sorting

The message deduplication of GlideIM relies on the global message ID. The message ID uses Redis Incr for the time being, and Meituan Leaf will be used later. The generation strategy of replacing the message ID in GlideIM is very simple. The message ID is obtained when the client sends a message, and the ID is attached when sending it. To achieve deduplication, and ordering of sender messages.

In the case of one-to-one chat , GlideIM only guarantees the order of the sender, and all messages in a session do not have to guarantee the order, which is relative to the cost and benefits. In group chat, all messages are guaranteed to be in order. In order, the group message will be sent with the continuous incrementing Seq of the current group, the received group messages will be sorted according to the Seq, and if the discontinuity is found, the discontinuous part of the message will be pulled.

## Four. Performance test

GlideIM messages face high concurrent throughput stress test:

![P](https://raw.githubusercontent.com/Glide-IM/Glide-IM/master/_art/msg_io_no_db.png)

### 4.1 Test Results

single instance deployment mode , 4H8G, 100Mbps broadband theoretically supports 20w active users chatting online at the same time, at this time bandwidth is the bottleneck of performance.

### 4.2 Test process

Server configuration

```
Windows 10
AMD R5 3600 6 cores 12 threads
16GB RAM
100Mbps network interface
```

#### Case 1 Test Process

A machine runs the server, B simulates the process of client connection, login and message sending, and runs the database

Simultaneously simulate 2000 clients, send a message every 60ms-200ms, each client sends 600 messages, a total of 1200_000 messages.

average load of the network is 90% (100Mbps), the throughput of about 30_000 messages per second, 15k upstream and downstream messages per second, the delivery rate is 100%, and the delay of all messages is <=20ms, refer to the result 1 picture.

Test situation :

```
Network 100%
2000 connections
5-20 messages/sec/connection
30k messages/ tps
100% delivery rate
```

#### Case 2 Test Process

This case is only to test the operation of the program under high concurrency conditions , and the reference value is limited to the shortcomings of the program itself and the logic limit of the code. In actual situations, we cannot ignore the influence of factors such as network speed.

Due to the need to remove the network rate limit , the service and the simulated client are both running on the same device and are not affected by network factors, simulating 10,000 links, each connection sends a message 50ms-100ms, a total of 800 messages are sent.

In this case, the message throughput limit is 28w. At this time, the CPU performance limit has been reached, and the occupancy rate is about 98%. The simulated cpu occupancy rate of the client is higher than that of the server, but it does not matter.

Test Case pprof CPU usage profiling data: [ cpu.out]( https://github.com/Glide-IM/Glide-IM/raw/master/_art/cpu_pprof_msg_io_no_db.out)

Test situation :

```
CPU 98%
Memory about 3 GB
10000 connections peak, start a connection every 10ms
10-20 messages/second/connection, a total of 800 messages are sent
280k messages/ tps peak
100% delivery rate
```

#### Test environment restrictions

Since there are only two computers in the test , there is a limit on the number of ports (deployed to the server, there is no such limit for user connections), in order to increase the size of the MySQL connection pool and improve concurrency, it is necessary to use fewer connections and send each connection more intensively message to test

maximum speed of the 100M network interface is only 12.5MB/s, and a message containing 20 Chinese characters requires at least 40B, plus other occupancy of the header, assuming 100B, then the actual maximum message concurrency of the 100M network interface is 12.5 x 1024 x 1024 / 100 = 13w strips, while removing
ACK, and then the number of upstream and downstream messages is half, and the actual number of 6w/s messages may be the theoretical limit of GlideIM in a 100Mbps environment.

#### Test code

Stand-alone performance test

Test code path : [/cmd/performance_test /]( https://github.com/Glide-IM/Glide-IM/tree/master/cmd/performance_test)

1. Start the server

```
go test -v -run= TestServerPerf
```

1. Start user simulation

```
go test -v -run= TestRunClient
```

### 4.3 Specifications for performance reference

#### Number of connections

Many IM projects like to use words like "million connections". Many people may be affected by the "million concurrency" of conventional HTTP services, and think that one million connections and one million concurrency have consistent performance reference values. Developers To attract attention to use the word million links, and indeed support millions of links.

In fact , the limit on the number of TCP (or WS) connections does not lie in the program itself, the system's maximum open file descriptors and memory limit the maximum number of connections, a TCP connection takes about 4-10KB of memory, a 16GB server theory Up to max support about 16*1024*
1024 / 5 ≈ 335w connection.

The number of active users is more valuable than the number of connections , for example, it supports 100w active connections, and each active connection sends a message every 10s.

#### Message throughput

Compared with the number of connections , the message throughput has more reference value, and the message throughput needs to be combined with the network rate for reference. The limited points of message throughput are mainly as follows:

1. The performance of the message delivery guarantee mechanism (that is, guaranteed delivery without redundancy)
2. The performance of the message data transmission protocol (a single message should be as small as possible)
3. Data packet control for non-user messages (such as heartbeat packets, content synchronization packets, etc.)
4. Network speed and quality (external factors)
5. Design flaws in the message distribution of the program itself

#### Other indicators

considering the performance of an IM may only be to reveal its own performance shortcomings. In most cases, horizontal comparison of other projects cannot get correct comparison results, and a single indicator cannot fully measure the pros and cons of the entire project. We should combine its business logic, design thinking, and learn its advantages to gain something.

1. Message delivery rate
2. Message Latency and Message Order Accuracy
3. Link keep alive and dead link (heartbeat)

### 4.4 Performance Metrics Estimation

1. Number of connections

Suppose the memory is set to M GB

```text
Theoretical number of connections = M * 1024m * 1024k / (4k/ tcp )

Conservatively estimated number of connections = M * 1024m * 1024k / (10k/ tcp )
```

2. Message throughput estimation

Let the network interface rate be S Mbps, and the average message size is K bytes (the total size when writing to a TCP connection).

```text
Message throughput = S/8*1024k*1024b/K
```

3. Estimate the number of active users

Let the throughput T, and let each user send a message every N seconds on average.

- Under the acknowledgment delivery mechanism (for the sending of each message, assuming that both the sender and the receiver are online , the server client needs ACK, and a total of 5 messages are required for upstream and downstream)

```text
Number of active users = T / 5 * N
```

- Only confirm delivery to the server

```text
Number of active users = T / 3 * N
```