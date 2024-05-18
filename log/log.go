package log

import (
	"log"
	"os"
)

var (
	RLogger = log.New(os.Stdout, "[Raft]", log.LstdFlags)
)
