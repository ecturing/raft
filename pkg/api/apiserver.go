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

type APIServer interface {
	Start() error
	Stop()
}

type Api struct {
	server *http.Server
}

// Start 启动API服务
/*
 可能潜在bug，select可能需要等待，先这样写
*/
func (a Api) Start() error {
	//TODO implement me
	errChan := make(chan error, 1)
	go func() {
		err := a.server.ListenAndServe()
		if err != nil {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		return err
	default:

	}
	go stopListen(httpCtx, a.server)
	return nil
}

func (a Api) Stop() {
	//TODO implement me
	cf()
	time.Sleep(1 * time.Second)
	panic("implement me")
}

func NewApiServer() Api {
	mux := http.NewServeMux()
	mux.HandleFunc("/sublog", submitLog)
	mux.HandleFunc("/getlogs", getLogs)
	return Api{
		server: &http.Server{
			Addr:    ":8080",
			Handler: mux,
		},
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

// 外部提交日志到集群的接口
func submitLog(r http.ResponseWriter, w *http.Request) {
	// entry := &LogsEntry{
	// 	metaData: w.FormValue("metaData"),
	// }
}

// 外部获取集群日志的接口,用于查询集群日志
func getLogs(r http.ResponseWriter, w *http.Request) {

}
