package builder

import (
	"fmt"

	"github.com/louislouislouislouis/oasnake/app/pkg/compiler"
	"github.com/louislouislouislouis/oasnake/app/pkg/generator"
)

type Builder struct {
	generator *generator.Generator
	compiler  compiler.Compiler
	config    *BuiderConfig
}

func NewBuilder(cfg *BuiderConfig) *Builder {
	var c compiler.Compiler
	if cfg.NeedToCompile() {
		if cfg.CompilerConfig.CompileWithDocker {
			c, _ = compiler.NewCompiler(compiler.DockerCompilerType, cfg.CompilerConfig)
		}
		if cfg.CompilerConfig.CompileWithGo {
			c, _ = compiler.NewCompiler(compiler.GoCompilerType, cfg.CompilerConfig)
		}
	}
	return &Builder{
		generator: generator.NewGenerator(cfg.GeneratorConfig),
		compiler:  c,
		config:    cfg,
	}
}

func (b *Builder) validateAndPrepareConfig() error {
	if b.compiler != nil {
		return fmt.Errorf("compiler is not set")
	}

	// Sanitize Generator Config
	b.generator.Config.WithCompilerFile = b.config.NeedToCompile()
	b.generator.Config.OutputDirectory = b.config.OutputDirectory

	// Sanitize Compiler Config
	if b.config.NeedToCompile() {
		// Sanitize Binary Name
		if b.compiler.GetConfig().BinaryName == "" {
			b.compiler.GetConfig().BinaryName = b.generator.Config.CommandName
		}

		b.compiler.GetConfig().OutputDirectory = b.config.OutputDirectory

		// Verify Target Arguments are present
		if b.compiler.GetConfig().TargetOs == "" || b.compiler.GetConfig().TargetArch == "" {
			return fmt.Errorf("a compilation has been requested, but no target OS or architecture are set for compilation")
		}

	}

	return nil
}

func (b *Builder) Build() error {
	err := b.validateAndPrepareConfig()
	if err != nil {
		return err
	}
	rootUsage, err := b.generator.Generate()
	if err != nil {
		return err
	}

	// The binary name is a result of the rootUsage Compilation
	if b.config.NeedToCompile() {
		b.compiler.GetConfig().BinaryName = rootUsage
		return b.compiler.Compile()
	}

	return nil
}
