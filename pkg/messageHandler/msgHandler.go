package messageHandler

import (
	"raft/pkg/rpc/Entries"
	"raft/pkg/rpc/net"
)

type MsgHandler interface {
	Vote(arg Entries.RequestVoteArgs) Entries.RequestVoteReply     //向外方法，向外发出投票信息
	Dispatch(arg Entries.DispatchLogArgs) Entries.RequestVoteReply //向外方法，向外分发日志
	Register(arg Entries.RegisterArgs) Entries.RegisterReply       //向外方法，向主节点提供新节点信息
	Heartbeat(arg Entries.HeartbeatArgs) Entries.HeartbeatReply    //向外方法，向其他节点发送心跳信息
}

func (r RpcMsgHandler) Vote(arg Entries.RequestVoteArgs) Entries.RequestVoteReply {
	//TODO implement me
	panic("implement me")
}

func (r RpcMsgHandler) Dispatch(arg Entries.DispatchLogArgs) Entries.RequestVoteReply {
	//TODO implement me
	panic("implement me")
}

func (r RpcMsgHandler) Register(arg Entries.RegisterArgs) Entries.RegisterReply {
	//TODO implement me
	net.Call(1, methodName+"Register()", arg, nil)
	return Entries.RegisterReply{}
}

func (r RpcMsgHandler) Heartbeat(arg Entries.HeartbeatArgs) Entries.HeartbeatReply {
	//TODO implement me
	panic("implement me")
}

type RpcMsgHandler struct{}

const methodName = "RpcCore"

func NewMsgHandler() MsgHandler {
	return RpcMsgHandler{}
}
