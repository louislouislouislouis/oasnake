/* Package compiler handle generation of a binary file*/
package compiler

import "fmt"

type Compiler interface {
	Compile() error
	GetConfig() *CompilerConfig
}

// NewCompiler creates a new compiler based on the given type
func NewCompiler(compilerType CompilerType, cfg *CompilerConfig) (Compiler, error) {
	switch compilerType {
	case DockerCompilerType:
		return NewDockerCompiler(cfg), nil
	case GoCompilerType:
		return NewGoCompiler(cfg), nil
	default:
		return nil, fmt.Errorf("unknown compiler type: %d", compilerType)
	}
}
