package mate_server

import (
	"context"
	"tangthinker.work/DFES/gateway"
	mateServerPB "tangthinker.work/DFES/mate-server/proto"
	"tangthinker.work/DFES/utils"
)

const (
	DefaultRaftDir             = "./raft/"
	DefaultRetainSnapshotCount = 20
	DefaultRaftAddr            = "127.0.0.1:70001"
	DefaultServerName          = "shanliao-mate-server-1"
	DefaultFragmentSize        = 16 // 16MB
	DefaultFragmentReplicaSize = 3
	DefaultDataClientCacheSize = 20
)

var mateServer = NewMateServer(DefaultRaftDir, DefaultRetainSnapshotCount, DefaultRaftAddr, DefaultServerName, DefaultFragmentSize, DefaultFragmentReplicaSize)

var dataClientCache = utils.NewActionCache(20)

func init() {
	mateServer.registryCenter = utils.NewRegistryClient(gateway.DefaultRegistryServerAddr)
	dataClientCache.RegisterGetFunc(func(key interface{}) interface{} {
		addr := key.(string)
		return utils.NewDataServerClient(addr)
	})
}

func InitRaft(firstNodeOrSingleMode bool) {
	mateServer.InitRaft(firstNodeOrSingleMode)
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

func SetServerName(serverName string) {
	mateServer.ServerName = serverName
	mateServer.idGenerator.ResetPrefix(serverName)
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
