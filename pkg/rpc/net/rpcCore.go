package net

import (
	"net"
	"net/http"
	"net/rpc"
	"raft/log"
	"raft/share"
	"time"
)

var (
	heartbeatTime = 1 * time.Second
)

type RpcCore struct {
	tick *time.Ticker //心跳定时器
}

// run 阻塞方法，注意协程启动
func (r *RpcCore) run() {
	log.RLogger.Println("Raft Node RPC is listening on", share.RPCServerPort)
	rpcRegisterErr := rpc.RegisterName("RpcCore", &RpcCore{tick: time.NewTicker(heartbeatTime)})
	if rpcRegisterErr != nil {
		log.RLogger.Panicln("rpcRegisterErr:", rpcRegisterErr)
		return
	}
	rpc.HandleHTTP()
	// Initialize the rpcServer
	rpcListener, netListenerErr := net.Listen("tcp", share.RPCServerPort)
	if netListenerErr != nil {
		panic(netListenerErr)
	}
	rpcHttpErr := http.Serve(rpcListener, nil)
	if rpcHttpErr != nil {
		log.RLogger.Panicln("rpcHttpErr:", rpcHttpErr)
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

func (r *RpcCore) setTimeTick(tick *time.Ticker) {
	r.tick = tick
}

func RPCStart() {
	rpcCore := &RpcCore{}
	rpcCore.setTimeTick(time.NewTicker(heartbeatTime))
	go rpcCore.run()
}
