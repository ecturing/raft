package share

const (
	NodePort = "8000"
)

var (
	HeartBeatReceived = make(chan struct{})
)
