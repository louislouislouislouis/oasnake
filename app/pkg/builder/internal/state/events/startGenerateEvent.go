package events

type StartGenerateEvent struct {
	Filename string
}

func (e StartGenerateEvent) Type() EventType { return StartGenerateCode }
