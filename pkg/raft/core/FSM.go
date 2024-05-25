package core

import (
	"context"
	"raft/log"
	"raft/pkg/rpc/Entries"
	"raft/pkg/rpc/messageHandler"
	"raft/share"
)

const (
	// State
	Leader = iota
	Follower
	Candidate
)

var (
	//heartBeatCtx不能是根上下文，因为heartBeatCtx会因为节点状态的改变而取消，根上下文除非程序结束，否则不会取消
	heartBeatCtx, hbCf = context.WithCancel(raft.Ctx)
)

type NodeState int

type StateMachine struct {
	state NodeState
}

type State interface {
	SetLeader()    //SetLeader() is a method that tries to transition to the leader state
	SetFollower()  //SetFollower() is a method that tries to transition to the follower state
	SetCandidate() //SetCandidate() is a method that tries to transition to the candidate state
}

func NewStateMachine() *StateMachine {
	return &StateMachine{
		state: Follower,
	}
}

// SetLeader Leader状态转移，Candidate->Leader
//转换为Leader状态后：
//1.启动api服务

func (sm *StateMachine) SetLeader() {
	// TODO  实现转移Leader函数
	reply := msgHandler.SendVote(Entries.RequestVoteArgs{Term: raft.term})
	if reply.VoteGranted {
		// TODO 实现成为Leader后所有包括API服务的启动等等，若有错误，则回滚
		sm.state = Leader
		messageHandler.SetLeaderInfo(&Entries.NetMeta{
			Ip:   NodeIP,
			Port: share.NodePort,
		}) //向Bootstrap节点注册该节点已成功成为Leader
		err := raft.apiCore.Start()
		raft.tickStop() //主节点不需要选举定时器
		go raft.msg.SendHeartbeat(heartBeatCtx, Entries.HeartbeatArgs{
			Term:     raft.term,
			LogIndex: raft.logs.GetIndex(),
		})
		if err != nil {
			return
		}
		log.RLogger.Println("Node State:Leader")
	} else {
		sm.SetFollower()
		log.RLogger.Println("try to Leader failed, Node State:Follower")
	}

}

// SetFollower Follower状态转移，Candidate->Follower
func (sm *StateMachine) SetFollower() {
	// Transition to follower state
	//TODO 实现转移Follower函数
	//subctx, _ := context.WithCancel(heartBeatCtx)
	sm.state = Follower
	hbCf()
}

// SetCandidate Candidate状态转移，Follower->Candidate
func (sm *StateMachine) SetCandidate() {
	// Transition to candidate state
	//TODO 实现转移Candidate函数
	sm.state = Candidate
	sm.SetLeader()
}
