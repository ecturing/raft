package messageHandler

import "raft/pkg/rpc/net"

type MsgHandler interface {
	Vote(arg interface{}) interface{}
	Dispatch(arg interface{}) interface{}
	Register(arg interface{}) interface{}
	Heartbeat(arg interface{}) interface{}
}

type RpcMsgHandler struct {
	Rpc net.RpcCore
}

var api net.RaftRpcAPI

func (r RpcMsgHandler) Vote(arg interface{}) interface{} {
	//TODO implement me
	api.Vote()
	panic("implement me")
}

func (r RpcMsgHandler) Dispatch(arg interface{}) interface{} {
	//TODO implement me
	api.DisPatchLog()
	panic("implement me")
}

func (r RpcMsgHandler) Register(arg interface{}) interface{} {
	//TODO implement me
	api.Register()
	panic("implement me")
}

func (r RpcMsgHandler) Heartbeat(arg interface{}) interface{} {
	//TODO implement me
	api.HeartBeat()
	panic("implement me")
}
