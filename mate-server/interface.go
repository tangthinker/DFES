package mate_server

import "github.com/shanliao420/DFES/api"

type MateService interface {
	api.Api
	Join(serverName string, addr string) error
}

type FileMate struct {
	FragmentCnt    int64
	Fragments      map[int64]*Fragment
	SourceHashCode string
}

type Fragment struct {
	Replicas []FragmentUint
}

type FragmentUint struct {
	DataNodeAddr string
	FragmentId   string
}
