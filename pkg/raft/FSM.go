package raft

import (
	"context"
	"raft/log"
	"raft/pkg/rpc/Entries"
)

const (
	// State
	Leader = iota
	Follower
	Candidate
)

var (
	//heartBeatCtx不能是根上下文，因为heartBeatCtx会因为节点状态的改变而取消，根上下文除非程序结束，否则不会取消
	heartBeatCtx, _ = context.WithCancel(raft.Ctx)
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
func (sm *StateMachine) SetLeader() {
	// TODO  实现转移Leader函数
	reply := msgHandler.Vote(Entries.RequestVoteArgs{Term: raft.term})
	if reply.VoteGranted {
		// TODO 实现成为Leader后所有包括API服务的启动等等，若有错误，则回滚
		sm.state = Leader
		log.Logger.Println("Node State:Leader")
	}

}

// SetFollower Follower状态转移，Candidate->Follower
func (sm *StateMachine) SetFollower() {
	// Transition to follower state
	//TODO 实现转移Follower函数
	//subctx, _ := context.WithCancel(heartBeatCtx)
	sm.state = Follower
}

// SetCandidate Candidate状态转移，Follower->Candidate
func (sm *StateMachine) SetCandidate() {
	// Transition to candidate state
	//TODO 实现转移Candidate函数
	sm.state = Candidate
	sm.SetLeader()
}
