package rpcProto

// rpc请求结构体
type RequestVoteArgs struct {
	Term uint //携带的任期号
}

type RequestVoteReply struct {
	VoteGranted bool //是否投票
}

type DispatchLogArgs struct{}
type DispatchLogReply struct{}

type NetMeta struct {
	ip   string
	port string
}
type RegisterArgs struct {
	NodeId uint64
	IP     string
	Port   string
}
type RegisterReply struct{}

type HeartbeatArgs struct {
	NodeId uint64
}
type HeartbeatReply struct{}
