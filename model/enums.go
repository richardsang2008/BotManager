package model

type LogLevel int

const (
	DEBUG LogLevel = 1 + iota
	INFO
	ERROR
	WARNING
	PANIC
)

type MeansureUnit int

const (
	Meters MeansureUnit = 1 + iota
	Miles
)
