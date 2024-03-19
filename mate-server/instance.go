package mate_server

import (
	"context"
	"fmt"
	"github.com/shanliao420/DFES/gateway"
	id_generator "github.com/shanliao420/DFES/id-generator"
	mateServerPB "github.com/shanliao420/DFES/mate-server/proto"
	"github.com/shanliao420/DFES/utils"
	"log"
	"time"
)

const (
	DefaultRaftDir             = "./raft/"
	DefaultRetainSnapshotCount = 20
	DefaultRaftAddr            = "127.0.0.1:70001"
	DefaultServerName          = "shanliao-mate-server-1"
	DefaultFragmentSize        = 24 * 1024 * 1024 // 24MB
	DefaultFragmentReplicaSize = 3
	DefaultDataClientCacheSize = 20
)

var mateServer = NewMateServer(DefaultRaftDir, DefaultRetainSnapshotCount, DefaultRaftAddr, DefaultServerName, DefaultFragmentSize, DefaultFragmentReplicaSize)

var dataClientCache = utils.NewActionCache(DefaultDataClientCacheSize)

func Init() {
	mateServer.registryCenter = utils.NewRegistryClient(gateway.DefaultRegistryServerAddr)
	dataClientCache.RegisterGetFunc(func(key interface{}) interface{} {
		addr := key.(string)
		return utils.NewDataServerClient(addr)
	})
	serviceCnt, err := mateServer.registryCenter.GetHistoryAllServiceCnt(context.Background(), nil)
	if err != nil {
		log.Fatalln("init mate server in id node err:", err)
	}
	log.Println("init mate server id node -> ", serviceCnt.GetServiceCnt())
	mateServer.idGenerator = id_generator.NewSnowflakeIdGenerator(serviceCnt.ServiceCnt)
}

func InitRaft(firstNodeOrSingleMode bool) {
	mateServer.InitRaft(firstNodeOrSingleMode)
	go func() {
		for {
			addr, id := mateServer.raft.LeaderWithID()
			fmt.Println(mateServer.raft.State().String(), " leader addr:", addr, " leader id:", id)
			time.Sleep(1 * time.Second)
		}
	}()
	go func() {
		for {
			select {
			case <-mateServer.raft.LeaderCh():
				if mateServer.IsLeader() {
					_ = mateServer.applyLeaderChange(mateServer.localRpcAddr)
					log.Println("leader addr change:", mateServer.localRpcAddr)
				}
			}
		}
	}()
}

func Join(leaderAddr string) error {
	mateServer.leaderRpcAddr = leaderAddr
	leaderClient := utils.NewMateServerClient(leaderAddr)
	joinResp, err := leaderClient.Join(context.Background(), &mateServerPB.JoinRequest{
		ServerName: mateServer.ServerName,
		ServerAddr: mateServer.raftAddr,
	})
	if err != nil || !joinResp.JoinResult {
		return err
	}
	return nil
}

func SetRaftAddr(raftAddr string) {
	mateServer.raftAddr = raftAddr
}

func SetLocalRpcAddr(localRpcAddr string) {
	mateServer.localRpcAddr = localRpcAddr
}

func SetServerName(serverName string) {
	mateServer.ServerName = serverName
	SetRaftDir(mateServer.raftDir + serverName + "/")
}

func SetRetainSnapshotCount(retainSnapshotCount int) {
	mateServer.retainSnapshotCount = retainSnapshotCount
}

func SetRaftDir(raftDir string) {
	mateServer.raftDir = raftDir
}

func SetFragmentSize(fragmentSize int64) {
	mateServer.FragmentSize = fragmentSize
}

func SerFragmentReplicaSize(fragmentReplicaSize int64) {
	mateServer.FragmentReplicaSize = fragmentReplicaSize
}
