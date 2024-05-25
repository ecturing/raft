package Entries

// rpc请求结构体
type RequestVoteArgs struct {
	Term uint //携带的任期号
}

type RequestVoteReply struct {
	VoteGranted bool //是否投票
}

type DispatchLogArgs struct {
	Log string
}
type DispatchLogReply struct{}

type NetMeta struct {
	Ip   string
	Port string
}
type RegisterArgs struct {
	NodeId uint64
	IP     string
	Port   string
}
type RegisterReply struct{}

type HeartbeatArgs struct {
	Term     uint
	LogIndex int
}
type HeartbeatReply struct {
	LogBehind bool
}

type LeaderArgs struct {
	IP   string
	Port string
}
type LeaderReply struct {
	Empty bool
	IP    string
	Port  string
}
