package log

import (
	"log"
	"os"
)

func init() {
	RLogger.SetFlags(log.LstdFlags)
	RLogger.SetOutput(os.Stdout)
}

var (
	RLogger log.Logger
)

func SetPrefix(prefix string) {
	RLogger.SetPrefix(prefix)
}
