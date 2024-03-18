Distributed File Encryption System
==================================

Distributed file encryption system implement by golang.

DFES is a distributed system using the following features to implement:
1. [Protocol buff](https://github.com/protocolbuffers/protobuf)
2. [gRPC](https://github.com/grpc/grpc)
3. [Raft](https://github.com/hashicorp/raft)
4. Asymmetric/Symmetric encryption
5. [LRU cache](https://github.com/hashicorp/golang-lru)

Feature
-------
1. distributed system.
2. mate-server through raft protocol implement CP architecture.
3. data-server through Quorum NRW implement AP architecture.
4. using registry center manager all server information.
5. using gRPC implement the communication among services.
6. using multi-fragment to implement high availability.

Quick start
-----------
1. Using go mod tidy to install go dependencies.
    ```shell
    go mod tidy
    ```
2. Start **Registry-Center** in a terminal window.
   ```shell
   go run cmd/registry/main.go 
   ```
3. Start **Leader** **Mate-Server** in a terminal window
    ```shell
    go run cmd/mate-server/main.go -port 7001 -server-name mate-node-1 -raft-addr "127.0.0.1:9001" -leader-addr ""
    ```
4. Start **Follower** **Mater-Server** in other terminal window (optional)
    ```shell
    go run cmd/mate-server/main.go -port 7002 -server-name mate-node-2 -raft-addr "127.0.0.1:9002" -leader-addr "127.0.0.1:7001"
    go run cmd/mate-server/main.go -port 7003 -server-name mate-node-3 -raft-addr "127.0.0.1:9003" -leader-addr "127.0.0.1:7001"
    ```
5. Start **Data-Server** in other terminal window
    ```shell
    go run cmd/data-server/main.go -port 8001 -server-name "data-node-1"
    go run cmd/data-server/main.go -port 8002 -server-name "data-node-2"
    go run cmd/data-server/main.go -port 8003 -server-name "data-node-3"
    ```
6. Test your server, you can use the given example in cmd/test-*, have fun.

Structure
---------
Continuously supplementing...
