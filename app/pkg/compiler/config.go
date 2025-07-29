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

type CompilerType int

const (
	DockerCompilerType CompilerType = iota
	GoCompilerType
)
