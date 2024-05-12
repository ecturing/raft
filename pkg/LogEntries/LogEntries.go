package LogEntries

type Logs interface {
	persist()
}

type LogsEntry struct {
	index    int      //当前日志最新索引
	term     int      //当前日志所属任期号
	metaData []string //当前日志元数据
}

// persist : 持久化日志
func (l *LogsEntry) persist() {

}

func (l *LogsEntry) AppendLog(log string) {
	l.metaData = append(l.metaData, log)
	l.index++
}
