package mate_server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb/v2"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	dataServerPB "tangthinker.work/DFES/data-server/proto"
	"tangthinker.work/DFES/gateway"
	gatewayPB "tangthinker.work/DFES/gateway/proto"
	idGenerator "tangthinker.work/DFES/id-generator"
	"tangthinker.work/DFES/utils"
	"time"
)

type MateServer struct {
	raftDir             string
	retainSnapshotCount int
	raftAddr            string
	ServerName          string // as ServerID in raft
	FragmentSize        int64  // uint MB
	FragmentReplicaSize int64

	mutex *sync.RWMutex

	raft *raft.Raft

	FileMates map[string]FileMate // file dataId -> FileMate

	idGenerator *idGenerator.SequenceIdGenerator

	registryCenter gatewayPB.RegistryClient

	leaderRpcAddr string

	localRpcAddr string
}

const (
	DefaultRaftTimeout = 5 * time.Second // 5s
)

func (ms *MateServer) IsLeader() bool {
	if addr, _ := ms.raft.LeaderWithID(); addr == "" {
		ms.raft.LeadershipTransfer()
	}
	return ms.raft.State() == raft.Leader
}

func (ms *MateServer) GetLeaderInfo() (raft.ServerAddress, raft.ServerID) {
	leaderAddr, leaderServerName := ms.raft.LeaderWithID()
	return leaderAddr, leaderServerName
}

func (ms *MateServer) Push(ctx context.Context, data []byte) (string, error) {
	// 首先获得所有data-node
	// 然后算出分片分布
	// push 并存储fragment信息 使用Apply
	getRes, err := ms.registryCenter.GetProvideServices(ctx, &gatewayPB.GetProvideInfo{
		ServiceType: gateway.DateService,
	})
	if err != nil {
		log.Println("get data nodes rpc error:", err)
		return "", err
	}
	if !getRes.GetResult {
		log.Println("get data nodes logic error:", getRes.GetResult)
		return "", errors.New("get data nodes logic error")
	}
	fragmentCnt := (int64((len(data))/1024) / 1024) / ms.FragmentSize // MB
	floatCnt := (float64(len(data)) / 1024 / 1024) / float64(ms.FragmentSize)
	hasRest := strings.Split(fmt.Sprintf("%.1f", floatCnt), ".")[1] == "0"
	if hasRest {
		fragmentCnt++
	}
	var fileMate FileMate
	fileMate.SourceHashCode = utils.Hash(data)
	fileMate.FragmentCnt = fragmentCnt
	fileMate.Fragments = make(map[int64]*Fragment)
	for i := int64(0); i < fragmentCnt; i++ {
		log.Println("push fragment", i, " total:", fragmentCnt)
		var fragment Fragment
		fragment.Replicas = make([]FragmentUint, 0)
		left := i * (ms.FragmentSize * 1024 * 1024)
		right := left + (ms.FragmentSize * 1024 * 1024)
		if right >= int64(len(data)) {
			right = int64(len(data))
		}
		idxArr := findMachina(len(getRes.GetProvideServices()), ms.FragmentReplicaSize)
		for k, idx := range idxArr {
			log.Println("push fragment idx", k, " to node")
			target := getRes.GetProvideServices()[idx]
			targetAddr := target.ServiceAddress.Host + ":" + target.ServiceAddress.Port
			targetDataNode := dataClientCache.Get(targetAddr).(dataServerPB.DataServiceClient)
			pushRes, err := targetDataNode.Push(context.Background(), &dataServerPB.PushRequest{
				FragmentData: data[left:right],
			})
			if err != nil {
				log.Println("push fragment[", left, ":", right, "] rpc error")
			}
			fragment.Replicas = append(fragment.Replicas, FragmentUint{
				FragmentId:   pushRes.FragmentId,
				DataNodeAddr: targetAddr,
			})
			log.Println("push fragment idx", k, " to node successful ")
		}
		log.Println("push fragment", i, " successful,  total:", fragmentCnt)
		fileMate.Fragments[i] = &fragment
	}
	fileMateId := ms.idGenerator.Next()
	return fileMateId, ms.applyPush(fileMateId, fileMate)
}

func (ms *MateServer) applyPush(fileMateId string, mate FileMate) error {
	comm := &command{
		FileId:   fileMateId,
		FileMate: mate,
		Op:       opPush,
	}
	b, err := json.Marshal(comm)
	if err != nil {
		return err
	}

	f := ms.raft.Apply(b, DefaultRaftTimeout)
	return f.Error()
}

func findMachina(size int, fragmentReplicaSize int64) []int {
	var ret []int
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i := int64(0); i < fragmentReplicaSize; i++ {
		ret = append(ret, r.Intn(size))
	}
	return ret
}

func (ms *MateServer) PushStream(ctx context.Context, stream io.Reader) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (ms *MateServer) Get(ctx context.Context, id string) ([]byte, error) {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()
	fileMate := ms.FileMates[id]
	var ret []byte
	for i := int64(0); i < fileMate.FragmentCnt; i++ {
		current := fileMate.Fragments[i]
		for idx, dataNode := range current.Replicas {
			node := dataClientCache.Get(dataNode.DataNodeAddr).(dataServerPB.DataServiceClient)
			getResp, err := node.Get(ctx, &dataServerPB.GetRequest{
				FragmentId: dataNode.FragmentId,
			})
			if err != nil {
				log.Println("get fragment from node", dataNode.DataNodeAddr, " rpc error", err)
				if idx == len(current.Replicas)-1 { // 已经到最后一个了，最后一个也出现了错误，文件无法复原
					return nil, err
				}
				continue
			}
			ret = append(ret, getResp.FragmentData...)
			break
		}
	}
	if fileMate.SourceHashCode != utils.Hash(ret) {
		return nil, errors.New("GetData hash not equal source hash code")
	}
	return ret, nil
}

func (ms *MateServer) GetStream(ctx context.Context, id string) (io.Reader, error) {
	//TODO implement me
	panic("implement me")
}

func (ms *MateServer) Delete(ctx context.Context, id string) (bool, error) {
	// 先删掉元信息mate，再删除dataNode上的真实文件
	ms.mutex.RLock()
	fileMate := ms.FileMates[id]
	ms.mutex.RUnlock()
	err := ms.applyDelete(id)
	if err != nil {
		return false, err
	}
	for i := int64(0); i < fileMate.FragmentCnt; i++ {
		current := fileMate.Fragments[i]
		for _, dataNode := range current.Replicas {
			node := dataClientCache.Get(dataNode.DataNodeAddr).(dataServerPB.DataServiceClient)
			delResp, err := node.Delete(ctx, &dataServerPB.DeleteRequest{
				FragmentId: dataNode.FragmentId,
			})
			if err != nil || !delResp.DeleteResult {
				log.Println("delete fragment from node", dataNode.DataNodeAddr, " rpc error", err)
				continue
			}
		}
	}
	return true, nil
}

func (ms *MateServer) applyDelete(mateId string) error {
	comm := &command{
		FileId: mateId,
		Op:     opDelete,
	}
	b, err := json.Marshal(comm)
	if err != nil {
		return err
	}

	f := ms.raft.Apply(b, DefaultRaftTimeout)
	return f.Error()
}

func (ms *MateServer) applyLeaderChange(leaderAddr string) error {
	comm := &command{
		Op:         opLeaderChange,
		LeaderAddr: leaderAddr,
	}
	b, err := json.Marshal(comm)
	if err != nil {
		return err
	}

	f := ms.raft.Apply(b, DefaultRaftTimeout)
	return f.Error()
}

func (ms *MateServer) Join(serverName string, addr string) error {
	log.Printf("received join request for remote node %s at %s", serverName, addr)
	configFuture := ms.raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		log.Printf("failed to get raft configuration: %v", err)
		return err
	}
	for _, srv := range configFuture.Configuration().Servers {
		// If a node already exists with either the joining node's ID or address,
		// that node may need to be removed from the config first.
		if srv.ID == raft.ServerID(serverName) || srv.Address == raft.ServerAddress(addr) {
			// However if *both* the ID and the address are the same, then nothing -- not even
			// a join operation -- is needed.
			if srv.Address == raft.ServerAddress(addr) && srv.ID == raft.ServerID(serverName) {
				log.Printf("node %s at %s already member of cluster, ignoring join request", serverName, addr)
				return nil
			}
			future := ms.raft.RemoveServer(srv.ID, 0, 0)
			if err := future.Error(); err != nil {
				return fmt.Errorf("error removing existing node %s at %s: %s", serverName, addr, err)
			}
		}
	}
	f := ms.raft.AddVoter(raft.ServerID(serverName), raft.ServerAddress(addr), 0, 0)
	if f.Error() != nil {
		return f.Error()
	}
	log.Printf("node %s at %s joined successfully", serverName, addr)
	return nil
}

func NewMateServer(raftDir string, retainSnapshotCount int, raftAddr string,
	serverName string, fragmentSize int64, fragmentReplicaSize int64) *MateServer {
	return &MateServer{
		raftDir:             raftDir,
		retainSnapshotCount: retainSnapshotCount,
		raftAddr:            raftAddr,
		ServerName:          serverName,
		mutex:               &sync.RWMutex{},
		FileMates:           make(map[string]FileMate),
		idGenerator:         idGenerator.NewSequenceIdGenerator(serverName),
		FragmentSize:        fragmentSize,
		FragmentReplicaSize: fragmentReplicaSize,
	}
}

func (ms *MateServer) InitRaft(firstNodeOrSingleMode bool) {
	err := utils.CreateFileIfNotExist(ms.raftDir, "raft-log.db")
	if err != nil {
		log.Fatal(err)
	}
	logStore, err := raftboltdb.NewBoltStore(filepath.Join(ms.raftDir, "raft-log.db"))
	if err != nil {
		log.Fatal(err)
	}
	err = utils.CreateFileIfNotExist(ms.raftDir, "raft-stable.db")
	if err != nil {
		log.Fatalln(err)
	}
	stableStore, err := raftboltdb.NewBoltStore(filepath.Join(ms.raftDir, "raft-stable.db"))
	if err != nil {
		log.Fatal(err)
	}
	snapshots, err := raft.NewFileSnapshotStore(ms.raftDir, ms.retainSnapshotCount, os.Stderr)
	if err != nil {
		log.Fatal(err)
	}
	addr, err := net.ResolveTCPAddr("tcp", ms.raftAddr)
	if err != nil {
		log.Fatal(err)
	}
	transport, err := raft.NewTCPTransport(ms.raftAddr, addr, 10, 5*time.Second, os.Stderr)
	if err != nil {
		log.Fatal(err)
	}
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(ms.ServerName)
	config.HeartbeatTimeout = 1 * time.Second
	raftNode, err := raft.NewRaft(config, (*fsm)(ms), logStore, stableStore, snapshots, transport)
	if err != nil {
		log.Fatal(err)
	}
	ms.raft = raftNode
	if firstNodeOrSingleMode {
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      raft.ServerID(ms.ServerName),
					Address: raft.ServerAddress(ms.raftAddr),
				},
			},
		}
		raftNode.BootstrapCluster(configuration)
	}

}

type fsm MateServer

func (f *fsm) Apply(l *raft.Log) interface{} {
	var c command
	if err := json.Unmarshal(l.Data, &c); err != nil {
		panic(fmt.Sprintf("failed to unmarshal command: %s", err.Error()))
	}

	switch c.Op {
	case opPush:
		return f.push(c.FileId, c.FileMate)
	case opDelete:
		return f.delete(c.FileId)
	case opLeaderChange:
		return f.leaderChange(c.LeaderAddr)
	default:
		panic(fmt.Sprintf("unrecognized command op: %s", c.Op))
	}
}

func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	s := make(map[string]FileMate)
	for k, v := range f.FileMates {
		s[k] = v
	}
	return &fsmSnapshots{
		snapshot: s,
	}, nil
}

func (f *fsm) Restore(snapshot io.ReadCloser) error {
	o := make(map[string]FileMate)
	if err := json.NewDecoder(snapshot).Decode(&o); err != nil {
		return err
	}

	// Set the state from the snapshot, no lock required according to
	// Hashicorp docs.
	f.FileMates = o
	return nil
}

func (f *fsm) push(fileId string, mate FileMate) interface{} {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.FileMates[fileId] = mate
	return fileId
}

func (f *fsm) delete(fileId string) interface{} {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	delete(f.FileMates, fileId)
	return nil
}

func (f *fsm) leaderChange(addr string) interface{} {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.leaderRpcAddr = addr
	return nil
}

type fsmSnapshots struct {
	snapshot map[string]FileMate
}

func (f *fsmSnapshots) Persist(sink raft.SnapshotSink) error {
	err := func() error {
		// Encode data.
		b, err := json.Marshal(f.snapshot)
		if err != nil {
			return err
		}

		// Write data to sink.
		if _, err := sink.Write(b); err != nil {
			return err
		}

		// Close the sink.
		return sink.Close()
	}()

	if err != nil {
		sink.Cancel()
	}

	return err
}

func (f *fsmSnapshots) Release() {
}

type command struct {
	Op         string   `json:"op,omitempty"`
	FileMate   FileMate `json:"file-mate,omitempty"`
	FileId     string   `json:"file-id,omitempty"`
	LeaderAddr string   `json:"leader-addr,omitempty"`
}

const (
	opPush         = "push"
	opDelete       = "delete"
	opLeaderChange = "leaderChange"
)