package raft

import (
	"context"
	"net/http"
	"os"
	"raft/log"
	"raft/pkg/api"
	"raft/pkg/messageHandler"
	"raft/pkg/rpc/Entries"
	"raft/pkg/rpc/net"
	"time"
)

const (
	RESETTIME = 10 * time.Second
)

var (
	//NodeIP 对于docker环境在启动容器时会填入容器的IP地址
	NodeIP     = os.Getenv("IP")
	msgHandler = messageHandler.NewMsgHandler()
	ctx, cf    = context.WithCancel(context.Background())
	raft       = NewCore(ctx, cf)
)

// ----------------------------------------RaftCore部分----------------------------------------
type RaftCore struct {
	nodeId   uint64                    //节点ID
	ticker   *time.Ticker              //内置状态定时器
	fsm      State                     //状态机
	term     uint                      //任期号
	NodeList []*Entries.NetMeta        //节点列表
	logs     []*Logs                   //Raft算法需要同步的日志
	msg      messageHandler.MsgHandler //消息处理器
	apiCore  *http.Server              //api服务，仅在Leader节点开启
	Ctx      context.Context           //系统上下文
	cf       context.CancelFunc        //取消函数
}

// NewCore 创建一个Raft核心/**
func NewCore(ctx context.Context, cf context.CancelFunc) *RaftCore {
	return &RaftCore{nodeId: 0,
		ticker:   time.NewTicker(RESETTIME),
		fsm:      NewStateMachine(),
		term:     0,
		NodeList: make([]*Entries.NetMeta, 0),
		logs:     make([]*Logs, 0),
		msg:      msgHandler,
		apiCore:  api.NewApiServer(),
		Ctx:      ctx,
		cf:       cf,
	}
}

func (c *RaftCore) Start() {
	subCtx := context.WithoutCancel(c.Ctx)
	go c.termTick(subCtx) //注册节点，bootstrap节点会返回所有当前集群节点信息
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
			c.fsm.SetLeader()
		case <-net.HeartBeatReceived:
			c.ticker.Reset(RESETTIME)
		case <-ctx.Done(): //根上下文取消，停止定时器
			c.ticker.Stop()
			return
		}
	}
}

func (c *RaftCore) Register() {
	arg := Entries.RegisterArgs{
		NodeId: c.nodeId,
		IP:     NodeIP,
		Port:   net.NodePort,
	}
	msgHandler.Register(arg)
}

func (c *RaftCore) GetState() State {
	return c.fsm
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
