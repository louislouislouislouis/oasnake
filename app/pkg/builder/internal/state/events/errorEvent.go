package events

type ErrorEvent struct {
	Error error
}

func (e ErrorEvent) Type() EventType { return Error }
