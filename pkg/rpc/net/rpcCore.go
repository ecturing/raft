package net

import (
	"net"
	"net/http"
	"net/rpc"
	"raft/log"
	"raft/share/rpcProto"
	"time"
)

const NodePort = ":1234"

var (
	HeartBeatReceived = make(chan struct{})
	heartbeatTime     = 1 * time.Second
)

type RpcCore struct {
	tick *time.Ticker //心跳定时器
}
type RaftRpcAPI interface {
	Vote(args rpcProto.RequestVoteArgs, reply *rpcProto.RequestVoteReply) error
	DisPatchLog(args rpcProto.RequestVoteArgs, reply *rpcProto.RequestVoteReply) error
	Register(args rpcProto.RegisterArgs, reply *rpcProto.RegisterReply) error
	HeartBeat(args rpcProto.HeartbeatArgs, reply *rpcProto.HeartbeatReply) error
}

// Run 阻塞方法，注意协程启动
func (r *RpcCore) Run() {
	log.Logger.Println("Raft Node RPC is listening on", NodePort)
	rpcRegisterErr := rpc.RegisterName("RpcCore", &RpcCore{tick: time.NewTicker(heartbeatTime)})
	if rpcRegisterErr != nil {
		log.Logger.Panicln("rpcRegisterErr:", rpcRegisterErr)
		return
	}
	rpc.HandleHTTP()
	// Initialize the rpcServer
	rpcListener, netListenerErr := net.Listen("tcp", ":1234")
	if netListenerErr != nil {
		panic(netListenerErr)
	}
	rpcHttpErr := http.Serve(rpcListener, nil)
	if rpcHttpErr != nil {
		log.Logger.Panicln("rpcHttpErr:", rpcHttpErr)
		return
	}
}

// Call 调用其他节点的rpc，调用方法需要请求到主节点
func Call(method string, args any, reply any) {
	rpcClient, err := rpc.DialHTTP("tcp", "address")
	log.Logger.Printf("Dialing to %s\n", "address")
	if err != nil {
		rpcErrHandler(err)
	}
	rpcCallErr := rpcClient.Call(method, args, reply)
	if rpcCallErr != nil {
		rpcErrHandler(rpcCallErr)
		return
	}
}

// rpcErrHandler rpc错误处理机制
func rpcErrHandler(err error) {
	switch err.(type) {
	case *net.OpError:
	case *net.DNSError:

	}
}
