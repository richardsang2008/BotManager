package model
type LogLevel int
const(
	DEBUG LogLevel = 1+ iota
	INFO
	ERROR
	WARNING
	PANIC
)