package share

var (
	DataManger = &dataManager{
		data: make(map[string]any),
	}
)

// VoteCheck : 检查投票资格
func VoteCheck(self uint, remote uint) bool {
	if self < remote {
		return true
	}
	return false //任期号小于等于自己的任期号，不具备当选leader的资格
}

type dataManager struct {
	data map[string]any
}

func (d *dataManager) Get(key string) any {
	return d.data[key]
}

func (d *dataManager) Set(key string, value any) {
	d.data[key] = value
}
