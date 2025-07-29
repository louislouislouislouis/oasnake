package compiler

type CompilerConfig struct {
	OutputDirectory   string
	BinaryName        string
	CompileWithGo     bool
	CompileWithDocker bool
	TargetOs          string
	TargetArch        string
}

func NewCompilerConfig() *CompilerConfig {
	return &CompilerConfig{}
}

func (cfg *CompilerConfig) NeedToBeCompiled() bool {
	return (cfg.BinaryName != "" || cfg.CompileWithGo || cfg.CompileWithDocker)
}
