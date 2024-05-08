package raft

import (
	"context"
	"raft/pkg/rpc/net"
	"raft/share"
	"raft/share/rpcProto"
)

const (
	// State
	Leader = iota
	Follower
	Candidate
)

var (
	//heartBeatCtx不能是根上下文，因为heartBeatCtx会因为节点状态的改变而取消，根上下文除非程序结束，否则不会取消
	heartBeatCtx, _ = context.WithCancel(share.DataManger.Get("RaftCore").(RaftCore).Ctx)
)

type NodeState int

type StateMachine struct {
	state NodeState
}

type State interface {
	TryToLeader()    //TryToLeader() is a method that tries to transition to the leader state
	TryToFollower()  //TryToFollower() is a method that tries to transition to the follower state
	TryToCandidate() //TryToCandidate() is a method that tries to transition to the candidate state
}

func NewStateMachine() *StateMachine {
	return &StateMachine{
		state: Follower,
	}
}

func (sm *StateMachine) TryToLeader() {
	// Transition to leader state
	core := share.DataManger.Get("RaftCore").(RaftCore)
	Vote := rpcProto.RequestVoteReply{
		VoteGranted: false,
	}
	net.Call("RpcCore.Vote", &rpcProto.RequestVoteArgs{Term: core.Term}, &Vote)

	if Vote.VoteGranted {
		// Transition to leader state & term + 1
		core.Term++
		sm.state = Leader
	}
}

func (sm *StateMachine) TryToFollower() {
	// Transition to follower state
	//subctx, _ := context.WithCancel(heartBeatCtx)
	sm.state = Follower
}

func (sm *StateMachine) TryToCandidate() {
	// Transition to candidate state
	sm.state = Candidate
}
