package net

import (
	"raft/log"
	"raft/pkg/LogEntries/logService"
	raftCore "raft/pkg/raft"
	"raft/pkg/rpc/Entries"
)

var (
	service = logService.NewLogService()
	raft    = raftCore.GetRaft()
)

// Vote 由Candidate节点调用，向本节点请求为其投票
func (r *RpcCore) Vote(args Entries.RequestVoteArgs, reply *Entries.RequestVoteReply) error {
	// TODO implement me
	if args.Term > raft.GetTerm() {
		reply.VoteGranted = true
	} else {
		reply.VoteGranted = false
	}
	return nil
}

// DisPatchLog 由主节点调用，传递日志到本节点
func (r *RpcCore) DisPatchLog(args Entries.DispatchLogArgs, reply *Entries.DispatchLogReply) error {
	return service.AppendLog(args.Log)
}

// Register 由新建节点调用，仅当本节点为Leader节点，接收注册新节点请求
func (r *RpcCore) Register(args Entries.RegisterArgs, reply *Entries.RegisterReply) error {
	raft.AppendNode(&Entries.NetMeta{
		Ip:   args.IP,
		Port: args.Port,
	})
	panic("implement me")
}

// HeartBeat 由主节点调用，传递心跳信息，主要进行主节点任期检查，防止该节点因停机未能及时更新集群任期,如果任期落后则日志也需要同步
func (r *RpcCore) HeartBeat(args Entries.HeartbeatArgs, reply *Entries.HeartbeatReply) error {
	if args.Term > raft.GetTerm() {
		log.Logger.Println("HeartBeat: term is maybe out of date,try to update term&log")
		// TODO try to implement update term and log
		raft.SetTerm(args.Term) // 更新任期
		reply.LogBehind = true
	} else {
		reply.LogBehind = false
	}
	return nil
}
