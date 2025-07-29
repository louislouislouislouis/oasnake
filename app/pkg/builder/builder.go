package builder

import (
	"fmt"

	"github.com/louislouislouislouis/oasnake/app/pkg/compiler"
	"github.com/louislouislouislouis/oasnake/app/pkg/generator"
	"github.com/rs/zerolog/log"
)

type Builder struct {
	generator *generator.Generator
	compiler  compiler.Compiler
	config    *BuiderConfig
}

func NewBuilder(cfg *BuiderConfig) *Builder {
	return &Builder{
		generator: generator.NewGenerator(cfg.GeneratorConfig),
		compiler:  compiler.NewDockerCompiler(cfg.CompilerConfig),
		config:    cfg,
	}
}

func (b *Builder) validateAndPrepareConfig() {
	b.generator.Config.ToInstall = b.compiler.GetConfig().NeedToBeCompiled()
	b.generator.Config.OutputDirectory = b.config.OutputDirectory
	b.compiler.GetConfig().OutputDirectory = b.config.OutputDirectory

	if b.compiler.GetConfig().BinaryName == "" {
		b.compiler.GetConfig().BinaryName = b.generator.Config.CommandName
	}
}

func (b *Builder) Build() error {
	// TODO: implement
	if b.compiler.GetConfig().CompileWithGo {
		log.Error().Msg("Installation with Go is not implemented yet, please use Docker instead.")
		return fmt.Errorf("not implemented")
	}
	b.validateAndPrepareConfig()
	rootUsage, err := b.generator.Generate()
	if err != nil {
		return err
	}
	b.compiler.GetConfig().BinaryName = rootUsage
	if b.compiler.GetConfig().NeedToBeCompiled() {
		return b.compiler.Compile()
	}
	return nil
}
