package net

import (
	"net"
	"net/http"
	"net/rpc"
	"raft/log"
	"raft/pkg/rpc/Entries"
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
	Vote(args Entries.RequestVoteArgs, reply *Entries.RequestVoteReply) error        // Vote 投票函数
	DisPatchLog(args Entries.RequestVoteArgs, reply *Entries.RequestVoteReply) error // DisPatchLog 分发Log函数
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
func Call(nodes []*Entries.NetMeta, method string, args any, reply any) {
	for _, meta := range nodes {
		rpcClient, err := rpc.DialHTTP("tcp", meta.Ip+":"+meta.Port)
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
}

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
