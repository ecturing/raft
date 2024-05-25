package rpc

import "raft/pkg/rpc/Entries"

type RaftRpcAPI interface {
	ReceiveVote(args Entries.RequestVoteArgs, reply *Entries.RequestVoteReply) error  // Vote 投票函数
	ReceiveLogs(args Entries.DispatchLogArgs, reply *Entries.DispatchLogReply) error  // DisPatchLog 分发Log函数
	ReceiveNewNodeInfo(args Entries.RegisterArgs, reply *Entries.RegisterReply) error // 接收新加入节点信息
	ReceiveHeartBeat(args Entries.HeartbeatArgs, reply *Entries.HeartbeatReply) error // HeartBeat 心跳函数
}
