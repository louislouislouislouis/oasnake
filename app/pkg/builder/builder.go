/* Package builder handle the builder for generating code */
package builder

import (
	"fmt"

	"github.com/louislouislouislouis/oasnake/app/pkg/builder/internal/state"
	"github.com/louislouislouislouis/oasnake/app/pkg/builder/internal/state/events"
	"github.com/louislouislouislouis/oasnake/app/pkg/compiler"
	"github.com/louislouislouislouis/oasnake/app/pkg/generator"
	"github.com/rs/zerolog/log"
)

type Builder struct {
	generator *generator.Generator
	compiler  compiler.Compiler
	config    *BuiderConfig
	sm        *state.StateManager
}

func NewBuilder(cfg *BuiderConfig) (*Builder, error) {
	var c compiler.Compiler
	var err error

	if cfg.NeedToCompile() {
		c, err = cfg.GetCompiler()
	}

	if err != nil {
		return nil, err
	}

	generator := generator.NewGenerator(cfg.GeneratorConfig)

	return &Builder{
		generator: generator,
		compiler:  c,
		config:    cfg,
		sm: state.NewStateManager(
			map[state.State]state.StateFunc{
				state.Generating: func(event events.Event) events.Event {
					rootUsage, err := generator.Generate()
					if err != nil {
						return events.ErrorEvent{Error: err}
					}
					return events.FinishGenerateCodeEvent{
						RootUsage: rootUsage,
					}
				},
				state.WithCode: func(event events.Event) events.Event {
					rootUsage := event.(events.FinishGenerateCodeEvent).RootUsage
					log.Debug().Msgf("code is now generated with root usage: %s", rootUsage)
					if cfg.NeedToCompile() {
						return events.StartCompileEvent{RootUsage: rootUsage}
					}
					return events.SuccessEvent{}
				},
				state.Compiling: func(event events.Event) events.Event {
					rootUsage := event.(events.StartCompileEvent).RootUsage
					c.GetConfig().BinaryName = rootUsage
					err := c.Compile()
					if err != nil {
						return events.ErrorEvent{Error: err}
					}
					return events.FinishCompileEvent{}
				},
				state.WithBinary: func(event events.Event) events.Event {
					if cfg.NeedToInstall() {
						return events.StartInstallEvent{}
					}
					return events.SuccessEvent{}
				},
			},
		),
	}, nil
}

func (b *Builder) validateAndSanitizeConfig() error {
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
	if err := b.validateAndSanitizeConfig(); err != nil {
		return err
	}
	return b.start()
}

func (b *Builder) start() error {
	return b.sm.Accept(
		events.StartGenerateEvent{Filename: b.generator.Config.InputFilePath},
	)
}
