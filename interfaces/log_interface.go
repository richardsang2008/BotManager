package interfaces

type LogInterface interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})
	Warning(args ...interface{})
	Panic(args ...interface{})
}
