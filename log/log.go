package log

import (
	"log"
	"os"
)

var (
	Logger = log.New(os.Stdout, "[Raft]", log.LstdFlags)
)
