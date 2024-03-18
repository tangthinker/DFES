Distributed File Encryption System
==================================

Distributed file encryption system implement by golang.

DFES is a distributed system using:
1. protocol buff
2. gRPC
3. Raft
4. asymmetric/symmetric encryption
5. LRU cache

feature to implement.

Feature
-------
1. distributed system
2. mate-server through raft protocol implement CP architecture
3. data-server through Quorum NRW implement AP architecture
4. using registry center manager all server information
5. using gRPC implement the communication among services

Quick start
-----------
1. Using go mod tidy to install go dependencies.
    ```shell
    go mod tidy
    ```
2. Start leader Mate-Server in a terminal window
    ```shell
    go run cmd/mate-server/main.go -port 7001 -server-name mate-node-1 -raft-addr "127.0.0.1:9001" -leader-addr ""
    ```
3. Start follower Mater-Server in other terminal window (optional)
    ```shell
    go run cmd/mate-server/main.go -port 7002 -server-name mate-node-2 -raft-addr "127.0.0.1:9002" -leader-addr "127.0.0.1:7001"
    go run cmd/mate-server/main.go -port 7003 -server-name mate-node-3 -raft-addr "127.0.0.1:9003" -leader-addr "127.0.0.1:7001"
    ```
4. Start Data-Server in other terminal window
    ```shell
    go run cmd/data-server/main.go -port 8001 -server-name "data-node-1"
    go run cmd/data-server/main.go -port 8002 -server-name "data-node-2"
    go run cmd/data-server/main.go -port 8003 -server-name "data-node-3"
    ```
5. Test your server, we give the example in cmd/test-*, have fun.

Structure
---------
Continuously supplementing...
