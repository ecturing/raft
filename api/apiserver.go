package api

import (
	"context"
	"net/http"
	"raft/log"
	"time"
)

// apiServer仅在raft为leader时启动，状态转移时关闭
var (
	httpCtx, cf = context.WithCancel(context.Background())
)

func NewApiServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/sublog", submitLog)
	mux.HandleFunc("/getlogs", getLogs)
	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
}

func stopListen(ctx context.Context, s *http.Server) {
	defer func() {
		log.Logger.Println("stop goroutine is exiting...")
	}()

	log.Logger.Println("stop goroutine is listening...")
	select {
	case <-ctx.Done():
		log.Logger.Println("have a stop signal, stop the server...")
		httpCtxShutdownErr := s.Shutdown(ctx)
		if httpCtxShutdownErr != nil {
			// TODO 处理httpCtx错误
			return
		}
		return
	}
}

func stop() {
	cf()
	time.Sleep(1 * time.Second)
}

func apiServerStart() {
	server := NewApiServer()
	go server.ListenAndServe()
	go stopListen(httpCtx, server)
}

// 外部提交日志到集群的接口
func submitLog(r http.ResponseWriter, w *http.Request) {
	// entry := &LogsEntry{
	// 	metaData: w.FormValue("metaData"),
	// }
}

// 外部获取集群日志的接口,用于查询集群日志
func getLogs(r http.ResponseWriter, w *http.Request) {

}
