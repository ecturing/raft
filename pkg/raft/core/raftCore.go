package core

import (
	"context"
	"os"
	"raft/log"
	"raft/pkg/api"
	"raft/pkg/logEntries"
	"raft/pkg/rpc/Entries"
	"raft/pkg/rpc/messageHandler"
	"raft/pkg/rpc/net"
	"raft/share"
	"strconv"
	"time"
)

const (
	RESETTIME = 10 * time.Second
)

var (
	//NodeIP 对于docker环境在启动容器时会填入容器的IP地址
	raft            = NewCore(rootCtx, rootCf)
	NodeIP          = os.Getenv("IP")
	msgHandler      = messageHandler.NewMsgHandler(raft.nodeList) //消息处理器服务注册
	rootCtx, rootCf = context.WithCancel(context.Background())
	apiCore         = api.NewApiServer()
)

// ----------------------------------------RaftCore部分----------------------------------------
type RaftCore struct {
	nodeId   uint64                    //节点ID
	ticker   *time.Ticker              //内置状态定时器
	fsm      State                     //状态机
	term     uint                      //任期号
	nodeList []*Entries.NetMeta        //节点列表
	logs     logEntries.Logs           //Raft算法需要同步的日志
	msg      messageHandler.MsgHandler //消息处理器
	apiCore  api.APIServer             //api服务，仅在Leader节点开启
	Ctx      context.Context           //系统上下文
	cf       context.CancelFunc        //取消函数
}

type hander struct {
}

// NewCore 创建一个Raft核心/**
func NewCore(ctx context.Context, cf context.CancelFunc) *RaftCore {
	return &RaftCore{
		nodeId:   0,
		ticker:   time.NewTicker(RESETTIME),
		fsm:      NewStateMachine(),
		term:     0,
		nodeList: make([]*Entries.NetMeta, 0),
		logs:     logEntries.NewLogsEntries(),
		msg:      nil,
		apiCore:  apiCore,
		Ctx:      ctx,
		cf:       cf,
	}
}

func (c *RaftCore) start() {
	log.RLogger.SetPrefix("[raft:" + strconv.Itoa(int(c.nodeId)) + "]")
	log.RLogger.Println("Raft Node is started")
	net.RPCRegister(hander{})
	net.RPCStart()
	c.register()
	c.termTick()
}

func Start() {
	raft.start()
}

func (c *RaftCore) Stop() {
	c.cf() //取消系统上下文
	time.Sleep(1 * time.Second)
}

// 阻塞方法，选举定时器
func (c *RaftCore) termTick() {
	log.RLogger.Println("Raft Node is ticking")
	defer log.RLogger.Println("Raft Node tick is stopped")
	for {
		select {
		case <-c.ticker.C:
			log.RLogger.Println("current Node termTick timeout try to be Leader")
			c.fsm.SetLeader()
		case <-share.HeartBeatReceived:
			c.ticker.Reset(RESETTIME)
		}
	}
}

func (c *RaftCore) tickStop() {
	c.ticker.Stop()
}

func (c *RaftCore) register() {
	arg := Entries.RegisterArgs{
		NodeId: c.nodeId,
		IP:     NodeIP,
		Port:   share.NodePort,
	}
	msgHandler.Register(arg)
	c.msg = msgHandler
}

func (c hander) GetTerm() uint {
	return raft.term
}
func (c hander) SetTerm(term uint) {
	raft.term = term
}
func (c hander) GetLogs() logEntries.Logs {
	return raft.logs
}
func GetRaft() *RaftCore {
	return raft
}
