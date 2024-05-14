package logService

type LogService interface {
	AppendLog(log string) error
	Persist() error
}

type service struct{}

func NewLogService() LogService {
	return service{}
}

func (s service) AppendLog(log string) error {
	//TODO implement me
	panic("implement me")
}

func (s service) Persist() error {
	//TODO implement me
	panic("implement me")
}
