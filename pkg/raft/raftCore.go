package raft

import (
	"context"
	"net/http"
	"os"
	"raft/api"
	"raft/log"
	"raft/pkg/messageHandler"
	"raft/pkg/rpc/net"
	"raft/share"
	"raft/share/rpcProto"
	"time"
)

const (
	RESETTIME = 10 * time.Second
)

var (
	NumOFNodes []uint
	//NodeIP 对于docker环境在启动容器时会填入容器的IP地址
	NodeIP     = os.Getenv("IP")
	msgHandler messageHandler.RpcMsgHandler
)

// ----------------------------------------RaftCore部分----------------------------------------
type RaftCore struct {
	NodeId   uint64                    //节点ID
	ticker   *time.Ticker              //内置状态定时器
	FSM      State                     //状态机
	Term     uint                      //任期号
	NodeList []*rpcProto.NetMeta       //节点列表
	logs     []*Logs                   //Raft算法需要同步的日志
	Msg      messageHandler.MsgHandler //消息处理器
	apiCore  *http.Server              //api服务，仅在Leader节点开启
	Ctx      context.Context           //系统上下文
	cf       context.CancelFunc        //取消函数
}

// NewCore 创建一个Raft核心/**
func NewCore(ctx context.Context, cf context.CancelFunc) *RaftCore {
	return &RaftCore{NodeId: 0,
		ticker:   time.NewTicker(RESETTIME),
		FSM:      NewStateMachine(),
		Term:     0,
		NodeList: make([]*rpcProto.NetMeta, 0),
		logs:     make([]*Logs, 0),
		Msg:      msgHandler,
		apiCore:  api.NewApiServer(),
		Ctx:      ctx,
		cf:       cf,
	}
}

func (c *RaftCore) Start() {
	subCtx := context.WithoutCancel(c.Ctx)
	go c.termTick(subCtx)
	net.Call("BOOTSTRAP.Register", &rpcProto.RegisterArgs{
		NodeId: c.NodeId,
		IP:     NodeIP,
		Port:   net.NodePort,
	}, &rpcProto.RegisterReply{}) //注册节点，bootstrap节点会返回所有当前集群节点信息
	log.Logger.Println("Raft Node is started")
}

func (c *RaftCore) Stop() {
	c.cf() //取消系统上下文
	time.Sleep(1 * time.Second)
}

// 阻塞方法，选举定时器
func (c *RaftCore) termTick(ctx context.Context) {
	log.Logger.Println("Raft Node is ticking")
	defer log.Logger.Println("Raft Node tick is stopped")
	for {
		select {
		case <-c.ticker.C:
			log.Logger.Println("本节点时间到期，尝试选举")
			c.FSM.TryToLeader()
		case <-net.HeartBeatReceived:
			c.ticker.Reset(RESETTIME)
		case <-ctx.Done(): //根上下文取消，停止定时器
			c.ticker.Stop()
			return
		}
	}
}

func (c *RaftCore) Register() {
	share.DataManger.Set("RaftCore", c)
}

func (c *RaftCore) GetState() State {
	return c.FSM
}

// ----------------------------------------日志部分----------------------------------------

type Logs interface {
	persist()
}

type LogsEntry struct {
	metaData string
}

func (l *LogsEntry) persist() {

}
