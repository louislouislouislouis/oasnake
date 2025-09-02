/* Package builder handle the builder for generating code */
package builder

import (
	"runtime"

	"github.com/louislouislouislouis/oasnake/app/pkg/builder/internal/state"
	"github.com/louislouislouislouis/oasnake/app/pkg/builder/internal/state/events"
	"github.com/louislouislouislouis/oasnake/app/pkg/compiler"
	"github.com/louislouislouislouis/oasnake/app/pkg/generator"
	"github.com/louislouislouislouis/oasnake/app/pkg/parser"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
	"github.com/rs/zerolog/log"
)

type Builder struct {
	generator *generator.Generator
	compiler  compiler.Compiler
	parser    *parser.Parser
	config    *BuiderConfig
	sm        *state.StateManager
}

var parserCodeGenConf = &codegen.Configuration{
	PackageName: "client",
	Generate: codegen.GenerateOptions{
		Client: true,
		Models: true,
	},
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
	parser := parser.NewParser(*cfg.ParserConfig)

	return &Builder{
		generator: generator,
		compiler:  c,
		config:    cfg,
		parser:    parser,
		sm: state.NewStateManager(
			map[state.State]state.StateFunc{
				state.Parsing: func(event events.Event) events.Event {
					rootCmd, spec, err := parser.ParseAndGetOpts()
					if err != nil {
						return events.ErrorEvent{Error: err}
					}
					return events.FinishParsingEvent{
						RootCmd: rootCmd,
						Spec:    spec,
					}
				},
				state.Generating: func(event events.Event) events.Event {
					finishParsingEvent := event.(events.FinishParsingEvent)
					rootUsage, err := generator.Generate(finishParsingEvent.RootCmd, finishParsingEvent.Spec)
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

		// Set default target OS and architecture if not provided
		if b.compiler.GetConfig().TargetOs == "" {
			b.compiler.GetConfig().TargetOs = runtime.GOOS
		}
		if b.compiler.GetConfig().TargetArch == "" {
			b.compiler.GetConfig().TargetArch = runtime.GOARCH
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
		events.StartParsingEvent{},
	)
}
