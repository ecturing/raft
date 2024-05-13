package messageHandler

import (
	"context"
	"log"
	"math/rand"
	"net/rpc"
	"raft/pkg/rpc/Entries"
	"sync"
)

type AsyncRpcClientManager struct {
	mu    sync.Mutex
	calls map[uint64]RpcCallback // 使用请求ID映射到其对应的回调函数
}

type RpcCallback func(resp interface{}, err error)

// NewAsyncRpcClientManager 创建一个新的异步RPC客户端管理器
func NewAsyncRpcClientManager() *AsyncRpcClientManager {
	return &AsyncRpcClientManager{
		calls: make(map[uint64]RpcCallback),
	}
}

func generateRequestId() uint64 {
	return uint64(rand.Int63())
}

// AsyncCall 发起异步RPC调用，并注册回调处理响应或错误
func (m *AsyncRpcClientManager) AsyncCall(ctx context.Context, meta *Entries.NetMeta, serviceMethod string, args any, resp any, callback RpcCallback) uint64 {
	reqID := generateRequestId()
	rpcClient, err := rpc.DialHTTP("tcp", meta.Ip+":"+meta.Port)
	if err != nil {
		panic("e")
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
			if err := call.Error; err != nil {
				log.Printf("RPC call failed: %v", err)
				callback(nil, err)
			} else {
				callback(resp, nil) // 假设resp已经被赋值
			}
			// 移除已完成的请求
			m.mu.Lock()
			delete(m.calls, reqID)
			m.mu.Unlock()
		case <-ctx.Done():
			// 上下文取消，手动移除并处理（可选）
			m.mu.Lock()
			delete(m.calls, reqID)
			m.mu.Unlock()
			log.Println("RPC call cancelled")
		}
	}()

	return reqID
}
