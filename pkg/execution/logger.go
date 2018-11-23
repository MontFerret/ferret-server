package execution

type Logger struct {
	log LogWriter
	job Job
}

func NewLogger(writer LogWriter, job Job) *Logger {
	return &Logger{writer, job}
}

func (l *Logger) Write(data []byte) (int, error) {
	return l.log.Write(l.job, data)
}
