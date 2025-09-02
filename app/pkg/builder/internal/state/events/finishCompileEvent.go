package events

type FinishCompileEvent struct{}

func (e FinishCompileEvent) Type() EventType { return FinishCompile }
