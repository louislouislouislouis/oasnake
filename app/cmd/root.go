/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "oasnake",
		Short: "Generate CLI REST Client",
		Long: `
	
	Use 'oasnake [command] --help' to get more information about a specific command.`,
	}

	rootCmd.AddCommand(NewGenerateCommand())

	return rootCmd
}
