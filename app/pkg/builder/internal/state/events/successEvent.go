package events

type SuccessEvent struct{}

func (e SuccessEvent) Type() EventType { return End }
