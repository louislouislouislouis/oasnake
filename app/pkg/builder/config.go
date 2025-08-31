package builder

import (
	"fmt"

	"github.com/louislouislouislouis/oasnake/app/pkg/compiler"
	"github.com/louislouislouislouis/oasnake/app/pkg/generator"
)

type BuiderConfig struct {
	GeneratorConfig *generator.GeneratorConfig
	CompilerConfig  *compiler.CompilerConfig
	OutputDirectory string
}

func (cfg *BuiderConfig) NeedToInstall() bool {
	// TODO: implement
	return false
}

func (cfg *BuiderConfig) NeedToCompile() bool {
	return (cfg.CompilerConfig.CompileWithGo || cfg.CompilerConfig.CompileWithDocker)
}

func (cfg *BuiderConfig) GetCompiler() (compiler.Compiler, error) {
	if !cfg.NeedToCompile() {
		return nil, fmt.Errorf("compilation not requested in config, no compiler available")
	}
	if cfg.CompilerConfig.CompileWithDocker {
		return compiler.NewCompiler(compiler.DockerCompilerType, cfg.CompilerConfig)
	}
	return compiler.NewCompiler(compiler.GoCompilerType, cfg.CompilerConfig)
}

func NewBuilderConfig() *BuiderConfig {
	return &BuiderConfig{
		generator.NewGeneratorConfig(),
		compiler.NewCompilerConfig(),
		"",
	}
}
