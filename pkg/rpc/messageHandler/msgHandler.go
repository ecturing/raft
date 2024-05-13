package messageHandler

import (
	"context"
	"fmt"
	"raft/log"
	"raft/pkg/rpc/Entries"
	"raft/share"
	"time"
)

var (
	manager    = NewAsyncRpcClientManager()
	baseCtx    = context.Background()
	rpcCtx, cf = context.WithTimeout(baseCtx, 10*time.Second)
)

type MsgHandler interface {
	Vote(arg Entries.RequestVoteArgs) Entries.RequestVoteReply     //向外方法，向外发出投票信息
	Dispatch(arg Entries.DispatchLogArgs) Entries.DispatchLogReply //向外方法，向外分发日志
	Register(arg Entries.RegisterArgs) Entries.RegisterReply       //向外方法，向主节点提供新节点信息
	Heartbeat(arg Entries.HeartbeatArgs) Entries.HeartbeatReply    //向外方法，向其他节点发送心跳信息
}

func (r RpcMsgHandler) Vote(arg Entries.RequestVoteArgs) Entries.RequestVoteReply {
	//TODO implement me
	voteNum := make([]bool, 5)
	var reply Entries.RequestVoteReply

	for _, meta := range r.nodeList {
		manager.AsyncCall(rpcCtx, meta, methodName+"Vote()", arg, reply, func(resp interface{}, err error) {
			if resp.(Entries.RequestVoteReply).VoteGranted {
				voteNum = append(voteNum, true)
			} else {
				voteNum = append(voteNum, false)
			}

		})
	}
	num := share.CountVote(voteNum, func(v bool) bool {
		return v == true
	})
	fmt.Println("num:", num)
	if num > len(r.nodeList)/2 {
		reply.VoteGranted = true
	} else {
		reply.VoteGranted = false
	}
	return reply
}

func (r RpcMsgHandler) Dispatch(arg Entries.DispatchLogArgs) Entries.DispatchLogReply {
	//TODO implement me
	for _, meta := range r.nodeList {
		//异步任务，防止阻塞主线程
		manager.AsyncCall(rpcCtx, meta, methodName+"Dispatch()", arg, Entries.DispatchLogReply{}, func(resp interface{}, err error) {

		})
	}
	return Entries.DispatchLogReply{}
}

func (r RpcMsgHandler) Register(arg Entries.RegisterArgs) Entries.RegisterReply {
	//TODO implement me

	return Entries.RegisterReply{}
}

func (r RpcMsgHandler) Heartbeat(arg Entries.HeartbeatArgs) Entries.HeartbeatReply {
	singleCallReply := Entries.HeartbeatReply{}
	for _, meta := range r.nodeList {
		manager.AsyncCall(rpcCtx, meta, methodName+"Heartbeat()", arg, &singleCallReply, func(resp interface{}, err error) {
			if err != nil {
				log.Logger.Printf("Call failed, meta: %s\n", meta)
			} else {
				if singleCallReply.LogBehind {
					//TODO try to process missing logs when current node logs behind the leader
					log.Logger.Printf("Log behind, meta: %s\n", meta)
				}
			}
		})
	}
	panic("implement me")
}

type RpcMsgHandler struct {
	nodeList []*Entries.NetMeta
}

const methodName = "RpcCore"

func NewMsgHandler(list []*Entries.NetMeta) MsgHandler {
	return RpcMsgHandler{
		nodeList: list,
	}
}
