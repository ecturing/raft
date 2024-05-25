package messageHandler

import (
	"context"
	"math/rand"
	"net/rpc"
	"raft/log"
	"raft/pkg/rpc/Entries"
	"sync"
)

var manager = newAsyncRpcClientManager()

type rpcClientManager struct {
	mu    sync.Mutex
	calls map[uint64]RpcCallback // 使用请求ID映射到其对应的回调函数
}

type SyncRpcClientManager struct{}

type RpcCallback func(resp any, err error)

// newAsyncRpcClientManager 创建一个新的异步RPC客户端管理器
func newAsyncRpcClientManager() *rpcClientManager {
	return &rpcClientManager{
		calls: make(map[uint64]RpcCallback),
	}
}
func generateRequestId() uint64 {
	return uint64(rand.Int63())

}

// asyncCall 发起异步RPC调用，并注册回调处理响应或错误
func (m *rpcClientManager) asyncCall(ctx context.Context, meta *Entries.NetMeta, serviceMethod string, args any, resp any, callback RpcCallback) uint64 {
	reqID := generateRequestId()
	rpcClient, err := rpc.DialHTTP("tcp", meta.Ip+meta.Port)
	if err != nil {
		log.RLogger.Panicln("RPC dial failed,", err)
	}
	call := rpcClient.Go(serviceMethod, args, resp, nil)
	// 使用Go方法发起异步调用
	m.mu.Lock()
	m.calls[reqID] = callback
	m.mu.Unlock()

	go func() {
		select {
		case <-call.Done:
			// 调用完成，执行回调处理结果
			callback(resp, call.Error) // 假设resp已经被赋值
			// 移除已完成的请求
			m.mu.Lock()
			delete(m.calls, reqID)
			m.mu.Unlock()
			log.RLogger.Println("async RPC call done")
			return //完成任务，协程结束
		case <-ctx.Done():
			// 上下文取消，手动移除并处理（可选）
			m.mu.Lock()
			delete(m.calls, reqID)
			m.mu.Unlock()
			log.RLogger.Println("async call goroutine stop:", serviceMethod)
			return
		}
	}()

	return reqID
}

func (m *rpcClientManager) syncCall(meta *Entries.NetMeta, serviceMethod string, args any, reply any) error {
	client, e := rpc.DialHTTP("tcp", meta.Ip+meta.Port)
	if e != nil {
		log.RLogger.Panicln("RPC dial failed,", e)
	}
	return client.Call(serviceMethod, args, reply)
}

// SyncCallFunc 同步调用
func SyncCallFunc[T any](meta *Entries.NetMeta, serviceMethod string, args any) (r T, err error) {
	/*
		有个相当奇怪的问题，虽然泛型指定了指针类型,采用下面代码仍然会报错:DecodeValue of unassignable value，
		错误是指数据解码到一个未分配空间的变量上，但是如果对返回值reply取地址再进行计算，就不会报错了
		var reply T
		e:=manager.syncCall(meta, serviceMethod, args,reply)
		return reply,e
	*/
	err = manager.syncCall(meta, serviceMethod, args, &r)
	return
}

func AsyncCallFunc[T any](ctx context.Context, meta *Entries.NetMeta, serviceMethod string, args any, callback RpcCallback) (r T) {
	manager.asyncCall(ctx, meta, serviceMethod, args, &r, callback)
	return
}
