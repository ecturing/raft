package logService

type LogService interface {
	AppendLog(log string) []string
	Persist()
}

type service struct{}

func NewLogService() LogService {
	return service{}
}

func (s service) AppendLog(log string) []string {
	//TODO implement me
	panic("implement me")
}

func (s service) Persist() {
	//TODO implement me
	panic("implement me")
}
