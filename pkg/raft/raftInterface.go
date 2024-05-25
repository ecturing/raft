package raft

import "raft/pkg/logEntries"

type RaftHandler interface {
	GetLogs() logEntries.Logs
	SetTerm(term uint)
	GetTerm() uint
}
