package err

var (
	rpcClientErr = raftError{1, "rpcErr"}
	rpcServerErr = raftError{2, "rpcErr"}
)

type raftError struct {
	code int
	msg  string
}

func (r raftError) Error() string {
	//TODO implement me
	panic("implement me")
}
