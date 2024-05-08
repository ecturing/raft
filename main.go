package main

import (
	"context"
	"raft/pkg/raft"
)

func main() {
	ctx, cf := context.WithCancel(context.Background())
	raft := raft.NewCore(ctx, cf)
	raft.Register()
	raft.Start()
}
