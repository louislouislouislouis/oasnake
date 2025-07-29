package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func NewDocCommand(documentedCmd *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "doc",
		Short: "Generate documentation for the CLI",
		Long:  `Generates documentation for the CLI in Markdown format.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			docDir := "./doc"

			// 1. If the 'docs' directory exists, remove it
			if _, err := os.Stat(docDir); err == nil {
				log.Debug().Msg("🧹 Removing existing 'docs' directory...")
				if err := os.RemoveAll(docDir); err != nil {
					return fmt.Errorf("failed to remove existing '%s': %w", docDir, err)
				}
			}

			// 2. Create a fresh 'docs' directory
			log.Debug().Msg("📁 Creating new 'docs' directory...")
			if err := os.MkdirAll(docDir, 0755); err != nil {
				return fmt.Errorf("failed to create directory '%s': %w", docDir, err)
			}

			// 3. Generate Markdown files for all Cobra commands
			log.Debug().Msg("📝 Generating CLI documentation in Markdown...")
			if err := doc.GenMarkdownTree(documentedCmd, docDir); err != nil {
				return fmt.Errorf("failed to generate documentation: %w", err)
			}

			log.Info().Msgf("✅ Documentation successfully generated in: %s", filepath.Clean(docDir))
			return nil
		},
	}

	return cmd
}
