package compiler

type CompilerConfig struct {
	BinaryName        string
	InstallWithGo     bool
	InstallWithDocker bool
}

func NewCompilerConfig() *CompilerConfig {
	return &CompilerConfig{}
}

type InstallationConfig struct {
	BinaryName        string
	InstallWithGo     bool
	InstallWithDocker bool
}

func (cfg *CompilerConfig) NeedToBeInstalled() bool {
	return (cfg.BinaryName != "" || cfg.InstallWithGo || cfg.InstallWithDocker)
}
