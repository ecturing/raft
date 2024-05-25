package share

const (
	NodePort          = "8000"
	BootStrapNodeIP   = "127.0.0.1"
	BootStrapNodePort = ":8080"
	ApiServerPort     = ":8080"
	RPCServerPort     = ":1234"
)

var (
	HeartBeatReceived = make(chan struct{})
)
