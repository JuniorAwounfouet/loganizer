package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "loganalyzer",
	Short: "Outil d'analyse de fichiers de logs",
	Long: `Loganalyzer est un outil CLI pour analyser des fichiers de logs 
provenant de diverses sources (serveurs, applications).`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
}
