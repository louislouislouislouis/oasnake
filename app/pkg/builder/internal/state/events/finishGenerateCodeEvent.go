package events

type FinishGenerateCodeEvent struct {
	RootUsage string
}

func (e FinishGenerateCodeEvent) Type() EventType { return FinishGenerateCode }
