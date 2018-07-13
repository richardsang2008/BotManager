//go:generate go-enum -f=enums.go
package model
// x ENUM(
//DEBUG,
//INFO,
//ERROR,
//WARNING,
//PANIC
// )
type LogLevel int


/*const (
	DEBUG LogLevel = 1 + iota
	INFO
	ERROR
	WARNING
	PANIC
)*/
//x ENUM(
// METERS,
//MILES
//)
type MeansureUnit int
/*const (
	Meters MeansureUnit = 1 + iota
	Miles
)*/
//x ENUM(
//USER,
//MODERATOR,
//	ADMINISTRATOR,
//	SUPER,
//	OWNER
//)
type Rights int
/*const (
	USER Rights = 1 + iota
	MODERATOR
	ADMINISTRATOR
	SUPER
	OWNER
)*/