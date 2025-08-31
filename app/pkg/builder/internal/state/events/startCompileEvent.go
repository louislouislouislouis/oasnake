package events

type StartCompileEvent struct {
	RootUsage string
}

func (e StartCompileEvent) Type() EventType { return StartCompile }
