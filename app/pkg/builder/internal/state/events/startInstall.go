package events

type StartInstallEvent struct{}

func (e StartInstallEvent) Type() EventType { return StartInstall }
