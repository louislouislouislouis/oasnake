package cmd

import (
	"github.com/louislouislouislouis/oasnake/app/pkg/builder"
	"github.com/spf13/cobra"
)

func NewGenerateCommand() *cobra.Command {
	builderCfg := builder.NewBuilderConfig()
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a binary terminal CLI for REST",
		RunE: func(cmd *cobra.Command, args []string) error {
			myBuilder, err := builder.NewBuilder(builderCfg)
			if err != nil {
				handleError(err)
				return err
			}
			err = myBuilder.Build()
			if err != nil {
				handleError(err)
			}
			return nil
		},
	}

	// Parser Required Flags
	cmd.PersistentFlags().StringVarP(&builderCfg.ParserConfig.InputFilePath, "input", "i", "", "the input OpenAPI file path")
	cmd.MarkPersistentFlagRequired("input")

	// Generator required flags
	cmd.PersistentFlags().StringVarP(&builderCfg.GeneratorConfig.Module, "module", "m", "", "The module name for the generated code")
	cmd.MarkPersistentFlagRequired("module")

	// Generator Optional Flags
	cmd.PersistentFlags().StringVar(&builderCfg.GeneratorConfig.ServerURL, "server-url", "", "Url of the server to use in the generated code, if not provided, it will be set to the server URL from the OpenAPI spec")
	cmd.PersistentFlags().BoolVar(&builderCfg.GeneratorConfig.WithModel, "with-model", false, "generate a model for the OpenAPI spec, this will generate a model in the output directory with the same name as the OpenAPI spec file, but with a .go extension. This is useful if you want to use the generated code in your own project.")
	cmd.PersistentFlags().StringVarP(&builderCfg.OutputDirectory, "output", "o", "out", "output directory for generated code - defaults to 'out' in the current directory.")
	cmd.PersistentFlags().StringVarP(&builderCfg.GeneratorConfig.CommandName, "name", "n", "", "The root command name (and usage), if not provided, it will be set to the info name from the OpenAPI spec. If it is not find, it will be a random name")

	// Compiler flags
	cmd.PersistentFlags().BoolVar(&builderCfg.CompilerConfig.Compile, "compile", false, "create binary using go compiler. If set to true, it would use by default the go compiler. You can override this by setting either --compile-with-go or --compile-with-docker to true.")
	cmd.PersistentFlags().BoolVar(&builderCfg.CompilerConfig.CompileWithGo, "compile-with-go", false, "create binary using go compiler. This will only work if you have go installed and in your PATH.")
	cmd.PersistentFlags().BoolVar(&builderCfg.CompilerConfig.CompileWithDocker, "compile-with-docker", false, "create binary using docker. This will only work if you have docker installed and in your PATH.")
	cmd.MarkFlagsMutuallyExclusive("compile-with-go", "compile-with-docker")
	cmd.PersistentFlags().StringVar(&builderCfg.CompilerConfig.TargetOs, "target-os", "", "OS for the generated binary. Would be setup as env var in the GOOS env while compiling. Defaults to the current OS if not specified.")
	cmd.PersistentFlags().StringVar(&builderCfg.CompilerConfig.TargetArch, "target-arch", "", "Architecture for the generated binary. Would be setup as env var in the GOARCH env while compiling. Defaults to the current architecture if not specified.")
	cmd.PersistentFlags().StringVarP(&builderCfg.CompilerConfig.BinaryName, "binary", "b", "", "Name of the binary file. If not specified, it will be the same as the command name.")

	return cmd
}
