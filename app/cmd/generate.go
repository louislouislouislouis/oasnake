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
		Long:  `Scans `,
		Run: func(cmd *cobra.Command, args []string) {
			myBuilder := builder.NewBuilder(builderCfg)

			if err := myBuilder.Build(); err != nil {
				handleError(err)
			}
		},
	}

	// Required Flags
	cmd.PersistentFlags().StringVarP(&builderCfg.GeneratorConfig.InputFilePath, "input", "i", "", "the input OpenAPI file path")
	cmd.MarkPersistentFlagRequired("input")
	cmd.PersistentFlags().StringVarP(&builderCfg.GeneratorConfig.Module, "module", "m", "", "The module name for the generated code")
	cmd.MarkPersistentFlagRequired("module")

	// Optional Flags
	cmd.PersistentFlags().StringVar(&builderCfg.GeneratorConfig.ServerURL, "server-url", "", "Url of the server to use in the generated code, if not provided, it will be set to the server URL from the OpenAPI spec")
	cmd.PersistentFlags().BoolVar(&builderCfg.GeneratorConfig.WithModel, "with-model", false, "generate a model for the OpenAPI spec, this will generate a model in the output directory with the same name as the OpenAPI spec file, but with a .go extension. This is useful if you want to use the generated code in your own project.")
	cmd.PersistentFlags().StringVarP(&builderCfg.GeneratorConfig.OutputDirectory, "output", "o", "out", "output directory for generated code - defaults to 'out' in the current directory.")
	cmd.PersistentFlags().StringVarP(&builderCfg.GeneratorConfig.CommandName, "name", "n", "", "The root command name (and usage), if not provided, it will be set to the info name from the OpenAPI spec. If it is not find, it will be a random name")

	// Binary and installation flags
	cmd.PersistentFlags().StringVarP(&builderCfg.CompilerConfig.BinaryName, "binary", "b", "", "Name of the binary file. If not specified, it will be the same as the command name.")
	cmd.PersistentFlags().BoolVar(&builderCfg.CompilerConfig.InstallWithGo, "install-with-go", false, "generate and install binary using Go install. This will only work if you have Go installed and in your PATH.")
	cmd.PersistentFlags().BoolVar(&builderCfg.CompilerConfig.InstallWithDocker, "install-with-docker", false, "generate and install binary using Docker. This will only work if you have Docker installed and in your PATH.")
	cmd.MarkFlagsMutuallyExclusive("install-with-go", "install-with-docker")

	return cmd
}
