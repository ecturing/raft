package raft

import (
	"context"
	"os"
	"raft/log"
	"raft/pkg/LogEntries"
	"raft/pkg/api"
	"raft/pkg/rpc/Entries"
	"raft/pkg/rpc/messageHandler"
	"raft/share"
	"time"
)

const (
	RESETTIME = 10 * time.Second
)

var (
	//NodeIP 对于docker环境在启动容器时会填入容器的IP地址
	raft       = NewCore(ctx, cf)
	NodeIP     = os.Getenv("IP")
	msgHandler = messageHandler.NewMsgHandler(raft.nodeList) //消息处理器服务注册
	ctx, cf    = context.WithCancel(context.Background())
	apiCore    = api.NewApiServer()
)

// ----------------------------------------RaftCore部分----------------------------------------
type RaftCore struct {
	nodeId   uint64                    //节点ID
	ticker   *time.Ticker              //内置状态定时器
	fsm      State                     //状态机
	term     uint                      //任期号
	nodeList []*Entries.NetMeta        //节点列表
	logs     []LogEntries.Logs         //Raft算法需要同步的日志
	msg      messageHandler.MsgHandler //消息处理器
	apiCore  api.APIServer             //api服务，仅在Leader节点开启
	Ctx      context.Context           //系统上下文
	cf       context.CancelFunc        //取消函数
}

// NewCore 创建一个Raft核心/**
func NewCore(ctx context.Context, cf context.CancelFunc) *RaftCore {
	return &RaftCore{nodeId: 0,
		ticker:   time.NewTicker(RESETTIME),
		fsm:      NewStateMachine(),
		term:     0,
		nodeList: make([]*Entries.NetMeta, 0),
		logs:     make([]LogEntries.Logs, 0),
		msg:      nil,
		apiCore:  apiCore,
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
		case <-share.HeartBeatReceived: //这个需要解耦
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
		Port:   share.NodePort,
	}
	msgHandler.Register(arg)
	c.msg = msgHandler
}

func (c *RaftCore) getLogsIndex() uint {
	return 0
}

func (c *RaftCore) GetTerm() uint {
	return c.term
}
func (c *RaftCore) SetTerm(term uint) {
	c.term = term
}

func (c *RaftCore) AppendLog(newLog string) {
	nLog := LogEntries.LogsEntry{
		Index:    1,
		Term:     c.term,
		MetaData: newLog,
	}
	c.logs = append(c.logs, &nLog)
}

func (c *RaftCore) AppendNode(meta *Entries.NetMeta) {
	c.nodeList = append(c.nodeList, meta)
}

func GetRaft() *RaftCore {
	return raft
}
