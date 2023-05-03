package keycloak_amg

type Logger struct{}

func (l *Logger) Errorf(format string, v ...interface{}) {}
func (l *Logger) Warnf(format string, v ...interface{})  {}
func (l *Logger) Debugf(format string, v ...interface{}) {}

func NewLogger() *Logger {
	return &Logger{}
}
