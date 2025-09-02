package events

import "fmt"

type EventType int

const (
	StartParsing EventType = iota
	FinishParsing
	StartGenerateCode
	FinishGenerateCode
	End
	StartCompile
	FinishCompile
	StartInstall
	FinishInstall
	Error
)

// String returns the string representation of the EventType
func (e EventType) String() string {
	switch e {
	case StartParsing:
		return "StartParsing"
	case FinishParsing:
		return "FinishParsing"
	case StartGenerateCode:
		return "StartGenerateCode"
	case FinishGenerateCode:
		return "FinishGenerateCode"
	case End:
		return "End"
	case StartCompile:
		return "StartCompile"
	case FinishCompile:
		return "FinishCompile"
	case StartInstall:
		return "StartInstall"
	case FinishInstall:
		return "FinishInstall"
	case Error:
		return "Error"
	default:
		return fmt.Sprintf("EventType(%d)", int(e))
	}
}
