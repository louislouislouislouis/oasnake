package builder

import (
	"github.com/louislouislouislouis/oasnake/app/pkg/compiler"
	"github.com/louislouislouislouis/oasnake/app/pkg/generator"
)

type Builder struct {
	generator *generator.Generator
	compiler  *compiler.Compiler
}

func NewBuilder(cfg *BuiderConfig) *Builder {
	return &Builder{
		generator: generator.NewGenerator(cfg.GeneratorConfig),
		compiler:  compiler.NewCompiler(cfg.CompilerConfig),
	}
}

func (b *Builder) Build() error {
	b.generator.Config.ToInstall = b.compiler.Config.NeedToBeInstalled()
	err := b.generator.Generate()
	if err != nil {
		return err
	}
	return b.compiler.Compile()
}
