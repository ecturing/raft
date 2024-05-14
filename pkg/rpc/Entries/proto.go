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
	LogIndex uint
}
type HeartbeatReply struct {
	LogBehind bool
}

type Logs interface {
	persist()
}

type LogsEntry struct {
	metaData string
}

func (l *LogsEntry) persist() {

}
