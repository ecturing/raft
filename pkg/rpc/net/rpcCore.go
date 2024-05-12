package net

import (
	"net"
	"net/http"
	"net/rpc"
	"raft/log"
	"raft/pkg/rpc/Entries"
	"time"
)

const NodePort = "127.0.0.1:1234"

var (
	heartbeatTime = 1 * time.Second
)

type RpcCore struct {
	tick *time.Ticker //心跳定时器
}
type RaftRpcAPI interface {
	Vote(args Entries.RequestVoteArgs, reply *Entries.RequestVoteReply) error        // Vote 投票函数
	DisPatchLog(args Entries.DispatchLogArgs, reply *Entries.DispatchLogReply) error // DisPatchLog 分发Log函数
	Register(args Entries.RegisterArgs, reply *Entries.RegisterReply) error          // Register 注册节点函数
	HeartBeat(args Entries.HeartbeatArgs, reply *Entries.HeartbeatReply) error       // HeartBeat 心跳函数
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
	rpcListener, netListenerErr := net.Listen("tcp", NodePort)
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

// rpcErrHandler rpc错误处理机制
func rpcErrHandler(err error) {
	switch err.(type) {
	case *net.OpError:
	case *net.DNSError:

	}
}

func (r *RpcCore) SetTimeTick(tick *time.Ticker) {
	r.tick = tick
}
