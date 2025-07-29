package builder

import (
	"github.com/louislouislouislouis/oasnake/app/pkg/compiler"
	"github.com/louislouislouislouis/oasnake/app/pkg/generator"
)

type BuiderConfig struct {
	GeneratorConfig *generator.GeneratorConfig
	CompilerConfig  *compiler.CompilerConfig
	OutputDirectory string
}

func (cfg *BuiderConfig) NeedToCompile() bool {
	return (cfg.CompilerConfig.CompileWithGo || cfg.CompilerConfig.CompileWithDocker)
}

func NewBuilderConfig() *BuiderConfig {
	return &BuiderConfig{
		generator.NewGeneratorConfig(),
		compiler.NewCompilerConfig(),
		"",
	}
}
