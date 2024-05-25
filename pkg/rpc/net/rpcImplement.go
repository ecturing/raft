package net

import (
	"raft/log"
	"raft/pkg/raft"
	"raft/pkg/rpc/Entries"
	"raft/share"
)

var (
	handler raft.RaftHandler
)

func RPCRegister(raft raft.RaftHandler) {
	handler = raft
}

// ReceiveVote 由Candidate节点调用，向本节点请求为其投票
func (r *RpcCore) ReceiveVote(args Entries.RequestVoteArgs, reply *Entries.RequestVoteReply) error {
	log.RLogger.Println("receive Vote Request")
	// TODO implement me
	if args.Term > handler.GetTerm() {
		reply.VoteGranted = true
	} else {
		reply.VoteGranted = false
	}
	return nil
}

// ReceiveLogs  由主节点调用，传递日志到本节点
func (r *RpcCore) ReceiveLogs(args Entries.DispatchLogArgs, reply *Entries.DispatchLogReply) error {
	return handler.GetLogs().AppendLog(args.Log)
}

// ReceiveNewNodeInfo  由新建节点调用，仅当本节点为Leader节点，接收注册新节点请求
func (r *RpcCore) ReceiveNewNodeInfo(args Entries.RegisterArgs, reply *Entries.RegisterReply) error {
	log.RLogger.Println("receive new node information:", args.IP, args.Port)
	//raft.AppendNode(&Entries.NetMeta{
	//	Ip:   args.IP,
	//	Port: args.Port,
	//})
	return nil
}

// ReceiveHeartBeat 由主节点调用，传递心跳信息，主要进行主节点任期检查，防止该节点因停机未能及时更新集群任期,如果任期落后则日志也需要同步
func (r *RpcCore) ReceiveHeartBeat(args Entries.HeartbeatArgs, reply *Entries.HeartbeatReply) error {
	if args.Term > handler.GetTerm() {
		log.RLogger.Println("ReceiveHeartBeat: term is maybe out of date,try to update term&log")
		// TODO try to implement update term and log
		handler.SetTerm(args.Term) // 更新任期
		reply.LogBehind = true
	} else {
		reply.LogBehind = false
		share.HeartBeatReceived <- struct{}{}
	}
	return nil
}
