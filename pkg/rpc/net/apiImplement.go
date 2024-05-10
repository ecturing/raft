package net

import (
	"raft/pkg/rpc/Entries"
)

// Vote 由Candidate节点调用，向本节点请求为其投票
func (r *RpcCore) Vote(args Entries.RequestVoteArgs, reply *Entries.RequestVoteReply) error {
	// TODO implement me
	panic("implement me")
}

// DisPatchLog 由主节点调用，传递日志到本节点
func (r *RpcCore) DisPatchLog(args Entries.RequestVoteArgs, reply *Entries.RequestVoteReply) error {
	// TODO implement me
	panic("implement me")
}

// Register 由新建节点调用，仅当本节点为Leader节点，接收注册新节点请求
func (r *RpcCore) Register(args Entries.RegisterArgs, reply *Entries.RegisterReply) error {
	// TODO implement me
	panic("implement me")
}

// HeartBeat 由主节点调用，传递心跳信息
func (r *RpcCore) HeartBeat(args Entries.HeartbeatArgs, reply *Entries.HeartbeatReply) error {
	// TODO implement me
	panic("implement me")
}
