package cmd

import (
	"github.com/louislouislouislouis/oasnake/app/pkg/config"
	"github.com/louislouislouislouis/oasnake/app/pkg/generator"
	"github.com/spf13/cobra"
)

func NewGenerateCommand() *cobra.Command {
	cfg := config.NewGenerateConfig()
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a binary terminal CLI for REST",
		Long:  `Scans `,
		Run: func(cmd *cobra.Command, args []string) {
			generator := generator.NewGenerator(cfg)
			if err := generator.Generate(); err != nil {
				handleError(err)
			}
		},
	}

	// Necessary Flags
	cmd.PersistentFlags().StringVarP(&cfg.InputFilePath, "input", "i", "", "the input OpenAPI file path")
	cmd.MarkPersistentFlagRequired("input")
	cmd.PersistentFlags().StringVarP(&cfg.Module, "module", "m", "", "The module name for the generated code")
	cmd.MarkPersistentFlagRequired("module")

	// Optional Flags
	cmd.PersistentFlags().StringVar(&cfg.ServerURL, "server-url", "", "Url of the server to use in the generated code, if not provided, it will be set to the server URL from the OpenAPI spec")
	cmd.PersistentFlags().BoolVar(&cfg.WithModel, "with-model", false, "generate a model for the OpenAPI spec, this will generate a model in the output directory with the same name as the OpenAPI spec file, but with a .go extension. This is useful if you want to use the generated code in your own project.")
	cmd.PersistentFlags().StringVarP(&cfg.OutputDirectory, "output", "o", "out", "output directory for generated code - defaults to 'out' in the current directory.")
	cmd.PersistentFlags().StringVarP(&cfg.CommandName, "name", "n", "", "The root command name (and usage), if not provided, it will be set to the info name from the OpenAPI spec. If it is not find, it will be a random name")

	// Binary and installation flags
	cmd.PersistentFlags().BoolVar(&cfg.Installation.HasToBeInstalled, "install", false, "If set to true, the binary will be installed after generation. If not set, the binary will not be installed. This will only work if you have Go or Docker installed and in your PATH.")
	cmd.PersistentFlags().StringVarP(&cfg.Installation.BinaryName, "binary", "b", "", "If specified, install flag will be set to true and the binary will be installed with the specified name. If not specified, it will be the same as the command name.")
	cmd.PersistentFlags().BoolVar(&cfg.Installation.InstallWithGo, "install-with-go", true, "generate and install binary using Go install. This will only work if you have Go installed and in your PATH.")
	cmd.PersistentFlags().BoolVar(&cfg.Installation.InstallWithDocker, "install-with-docker", false, "generate and install binary using Docker. This will only work if you have Docker installed and in your PATH.")
	cmd.MarkFlagsMutuallyExclusive("install-with-go", "install-with-docker")

	return cmd
}
