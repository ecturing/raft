package net

import "raft/share/rpcProto"

func (r *RpcCore) Vote(args rpcProto.RequestVoteArgs, reply *rpcProto.RequestVoteReply) error {
	//TODO implement me
	panic("implement me")
}

func (r *RpcCore) DisPatchLog(args rpcProto.RequestVoteArgs, reply *rpcProto.RequestVoteReply) error {
	//TODO implement me
	panic("implement me")
}

func (r *RpcCore) Register(args rpcProto.RegisterArgs, reply *rpcProto.RegisterReply) error {
	//TODO implement me
	panic("implement me")
}

func (r *RpcCore) HeartBeat(args rpcProto.HeartbeatArgs, reply *rpcProto.HeartbeatReply) error {
	//TODO implement me
	panic("implement me")
}
