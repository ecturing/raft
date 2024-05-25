package messageHandler

import (
	"context"
	rlog "raft/log"
	"raft/pkg/rpc/Entries"
	"raft/share"
	"time"
)

const (
	methodName    = "RpcCore."
	bootstrap     = "BootStrap."
	HEARTBEATTIME = 5 * time.Second
)

var (
	baseCtx = context.Background()

	/**
	//奇葩bug，context自创立开始即开始倒计时，与任期滴答器倒计时几乎同步，在触发leader转换时，刚准备rpc信息，该ctx已超时
	*/
	rpcCtx, cf = context.WithTimeout(baseCtx, 15*time.Second)
)

type MsgHandler interface {
	SendVote(arg Entries.RequestVoteArgs) Entries.RequestVoteReply    //向外方法，向外发出投票信息，仅当节点为Candidate时调用
	DispatchLog(arg Entries.DispatchLogArgs) Entries.DispatchLogReply //向外方法，向外分发日志，仅当节点为主节点时调用
	Register(arg Entries.RegisterArgs) Entries.RegisterReply          //向外方法，向主节点提供新节点信息，仅当节点为新节点且为follower时调用
	SendHeartbeat(ctx context.Context, arg Entries.HeartbeatArgs)     //向外方法，向其他节点发送心跳信息，仅当节点为主节点时调用
}

type RpcMsgHandler struct {
	nodeList []*Entries.NetMeta
	heabeat  *time.Ticker
}

func NewMsgHandler(list []*Entries.NetMeta) MsgHandler {
	return RpcMsgHandler{
		nodeList: list,
		heabeat:  time.NewTicker(HEARTBEATTIME),
	}
}

// getInfoOfLeader 内部方法，向BootStrap节点获取现在集群中leader节点信息
func getInfoOfLeader() (meta *Entries.NetMeta, empty bool) {
	bs := &Entries.NetMeta{
		Ip:   share.BootStrapNodeIP,
		Port: share.BootStrapNodePort,
	}
	result, err := SyncCallFunc[Entries.LeaderReply](bs, bootstrap+"GetLeaderInfo", Entries.LeaderArgs{})
	if err != nil {
		rlog.RLogger.Panicln("getInfoOfLeader error: please check result and err", result, err)
		return nil, true
	}
	if result.Empty {
		rlog.RLogger.Println("Leader empty,please waiting for Leader node")
		return nil, true
	}
	return &Entries.NetMeta{
		Ip:   result.IP,
		Port: result.Port,
	}, false
}

// SetLeaderInfo 向BootStrap节点设置leader节点信息
func SetLeaderInfo(meta *Entries.NetMeta) {
	bs := &Entries.NetMeta{
		Ip:   share.BootStrapNodeIP,
		Port: share.BootStrapNodePort,
	}
	AsyncCallFunc[Entries.LeaderReply](rpcCtx, bs, bootstrap+"ReceiveLeaderInfo", Entries.LeaderArgs{
		IP:   meta.Ip,
		Port: meta.Port,
	},
		func(resp any, err error) {
			if err != nil {
				rlog.RLogger.Println("SetLeaderInfo error:", err)
			} else {
				if resp.(*Entries.LeaderReply).Empty {
					rlog.RLogger.Println("SetLeaderInfo error:", err)
				} else {

					rlog.RLogger.Println("SetLeaderInfo success,ACK msg:", resp.(*Entries.LeaderReply).IP, resp.(*Entries.LeaderReply).Port)
				}
			}
		},
	)
}

func (r RpcMsgHandler) SendVote(arg Entries.RequestVoteArgs) Entries.RequestVoteReply {
	voteNum := make([]bool, 5)
	var reply Entries.RequestVoteReply
	rlog.RLogger.Println("SendVote to all nodes,Node Numbers:", len(r.nodeList))
	if len(r.nodeList) == 0 {
		reply.VoteGranted = true
		return reply
	}
	for _, meta := range r.nodeList {
		AsyncCallFunc[Entries.RequestVoteReply](rpcCtx, meta, methodName+"ReceiveVote", arg, func(resp any, err error) {
			if err != nil {
				cf()
				// TODO 错误处理
			} else {
				if resp.(Entries.RequestVoteReply).VoteGranted {
					voteNum = append(voteNum, true)
				} else {
					voteNum = append(voteNum, false)
				}
			}
		})
	}
	num := share.CountVote(voteNum, func(v bool) bool {
		return v == true
	})
	if num > len(r.nodeList)/2 {
		reply.VoteGranted = true
	} else {
		reply.VoteGranted = false
	}
	return reply
}

func (r RpcMsgHandler) DispatchLog(arg Entries.DispatchLogArgs) Entries.DispatchLogReply {
	//TODO implement me
	for _, meta := range r.nodeList {
		//异步任务，防止阻塞主线程
		AsyncCallFunc[Entries.DispatchLogReply](rpcCtx, meta, methodName+"ReceiveLogs", arg, func(resp any, err error) {

		})
	}
	return Entries.DispatchLogReply{}
}

func (r RpcMsgHandler) Register(arg Entries.RegisterArgs) Entries.RegisterReply {
	leader, empty := getInfoOfLeader()
	if empty {
		return Entries.RegisterReply{}
	}
	//TODO 尝试获取leader节点信息后向leader节点发送新节点注册信息，并传播给其他节点
	netInfo := &Entries.NetMeta{
		Ip:   leader.Ip,
		Port: leader.Port,
	}
	AsyncCallFunc[Entries.RegisterReply](rpcCtx,
		netInfo, //固定BootStrap节点用于实现节点注册信息
		methodName+"ReceiveNewNodeInfo",
		arg,
		func(resp any, err error) {
			if err != nil {
				cf()
				rlog.RLogger.Printf("Call failed, meta: %s\n,err:%v", netInfo, err)
			} else {
				rlog.RLogger.Printf("register success, meta: %s\n", netInfo)
			}
		})
	return Entries.RegisterReply{}
}

func sendHeartbeat(arg Entries.HeartbeatArgs, list []*Entries.NetMeta) Entries.HeartbeatReply {
	for _, meta := range list {
		AsyncCallFunc[Entries.HeartbeatReply](rpcCtx, meta, methodName+"ReceiveHeartBeat", arg, func(resp any, rpcCallErr error) {
			if rpcCallErr != nil {
				cf()
				rlog.RLogger.Printf("Call failed, meta: %s\n", meta)
			} else {
				if resp.(Entries.HeartbeatReply).LogBehind {
					//TODO try to process missing logs when current node logs behind the leader
					rlog.RLogger.Printf("Log behind, meta: %s\n", meta)
				}
			}
		})
	}
	return Entries.HeartbeatReply{}
}

// SendHeartbeat 阻塞方法，协程调用
func (r RpcMsgHandler) SendHeartbeat(ctx context.Context, arg Entries.HeartbeatArgs) {
	for {
		select {
		case <-r.heabeat.C:
			rlog.RLogger.Println("HeartBeat Tick,Send HeartBeat")
			sendHeartbeat(arg, r.nodeList)
		case <-ctx.Done():
			return
		}
	}
}
