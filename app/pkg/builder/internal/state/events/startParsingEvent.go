package events

type StartParsingEvent struct{}

func (e StartParsingEvent) Type() EventType { return StartParsing }
