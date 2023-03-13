package logger

// GinLoggerWriter ...
type GinLoggerWriter struct {
	Logger AipLogger
	IsErr  bool
}

// Write Write
func (loggerWriter *GinLoggerWriter) Write(p []byte) (int, error) {
	if loggerWriter.IsErr {
		loggerWriter.Logger.Error(string(p))
	} else {
		loggerWriter.Logger.Info(string(p))
	}
	return 0, nil
}
