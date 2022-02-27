package go_webserver

// -----------------------------------------------------------------------------

type loggerBridge struct {
	cb loggerCallback
}

type loggerCallback func (format string, args ...interface{})

// -----------------------------------------------------------------------------

func newLoggerBridge(cb loggerCallback) *loggerBridge {
	return &loggerBridge{
		cb:  cb,
	}
}

func (l *loggerBridge) Printf(format string, args ...interface{}) {
	l.cb(format, args...)
}
