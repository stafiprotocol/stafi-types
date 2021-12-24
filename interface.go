package types

type LoggerInterface interface {
	Info(msg string, ctx ...interface{})
	Warn(msg string, ctx ...interface{})
	Error(msg string, ctx ...interface{})
}

type SarpcInterface interface {
	RegCustomTypes(content []byte)
	GetSystemChain() (string, error)
}
