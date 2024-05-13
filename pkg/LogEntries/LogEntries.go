package LogEntries

type Logs interface {
	persist()
}

type LogsEntry struct {
	Index    int    //当前日志最新索引
	Term     uint   //当前日志所属任期号
	MetaData string //当前日志元数据
}

// persist : 持久化日志
func (l *LogsEntry) persist() {

}

func (l *LogsEntry) AppendLog(log string) {
	l.MetaData = log
	l.Index++
}
