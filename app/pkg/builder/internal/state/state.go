package state

type State int

const (
	Ready State = iota
	Generating
	WithCode
	Compiling
	WithBinary
	Installing
	Installed
	Success
	Failure
)

func (s State) String() string {
	switch s {
	case Ready:
		return "Ready"
	case Generating:
		return "Generating"
	case WithCode:
		return "WithCode"
	case Compiling:
		return "Compiling"
	case WithBinary:
		return "WithBinary"
	case Installing:
		return "Installing"
	case Installed:
		return "Installed"
	case Success:
		return "Success"
	case Failure:
		return "Failure"
	default:
		return "Unknown"
	}
}
