package logEntries

import rlog "raft/log"

type Logs interface {
	persist()
	AppendLog(log string) error
	GetIndex() int
}

type LogsEntry struct {
	Index    int      //当前日志最新索引
	Term     uint     //当前日志所属任期号
	MetaData []string //当前日志元数据
}

// persist : 持久化日志
func (l *LogsEntry) persist() {

}

func (l *LogsEntry) AppendLog(log string) error {
	defer func() {
		if err := recover(); err != nil {
			rlog.RLogger.Panicln("append log error")
		}
	}()
	l.MetaData = append(l.MetaData, log)
	l.Index++
	return nil
}

func (l *LogsEntry) GetIndex() int {
	return l.Index
}

func NewLogsEntries() *LogsEntry {
	return &LogsEntry{
		Index:    0,
		Term:     0,
		MetaData: make([]string, 0),
	}
}
