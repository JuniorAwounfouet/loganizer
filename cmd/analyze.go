package cmd

import (
	"fmt"
	"sync"

	"github.com/JuniorAwounfouet/go_loganizer/internal/analyzer"
	"github.com/JuniorAwounfouet/go_loganizer/internal/config"
	"github.com/JuniorAwounfouet/go_loganizer/internal/reporter"
	"github.com/spf13/cobra"
)

var (
	configPath string
	outputPath string
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyser les fichiers de logs spécifiés dans la configuration",
	Long: `Analyser plusieurs fichiers de logs en parallèle selon la configuration JSON.
Exemple: loganalyzer analyze -c config.json -o report.json`,
	Run: func(cmd *cobra.Command, args []string) {
		configs, err := config.LoadConfig(configPath)
		if err != nil {
			fmt.Printf("Erreur lors du chargement de la configuration: %v\n", err)
			return
		}

		fmt.Printf("Début de l'analyse de %d fichiers de logs...\n", len(configs))

		var wg sync.WaitGroup
		resultChan := make(chan analyzer.AnalysisResult, len(configs))

		// Lancer les goroutines d'analyse
		for _, cfg := range configs {
			wg.Add(1)
			go func(logID, filePath, logType string) {
				defer wg.Done()
				result := analyzer.AnalyzeLog(logID, filePath, logType)
				resultChan <- result
			}(cfg.ID, cfg.Path, cfg.Type)
		}

		// Attendre que toutes les analyses soient terminées
		wg.Wait()
		close(resultChan)

		// Collecter les résultats
		var results []analyzer.AnalysisResult
		for result := range resultChan {
			results = append(results, result)
		}

		fmt.Println("\nRésultats de l'analyse:")
		fmt.Println("======================")
		for _, result := range results {
			fmt.Printf("ID: %s\n", result.LogID)
			fmt.Printf("Chemin: %s\n", result.FilePath)
			fmt.Printf("Statut: %s\n", result.Status)
			fmt.Printf("Message: %s\n", result.Message)
			if result.ErrorDetails != "" {
				fmt.Printf("Erreur: %s\n", result.ErrorDetails)
			}
			fmt.Println("---")
		}

		// Exporter si un chemin de sortie est spécifié
		if outputPath != "" {
			err = reporter.ExportResults(results, outputPath)
			if err != nil {
				fmt.Printf("Erreur lors de l'exportation: %v\n", err)
			} else {
				fmt.Printf("Résultats exportés vers: %s\n", outputPath)
			}
		}

		successCount := 0
		for _, result := range results {
			if result.Status == "OK" {
				successCount++
			}
		}
		fmt.Printf("Résumé: %d/%d analyses réussies\n", successCount, len(results))
	},
}

func init() {
	analyzeCmd.Flags().StringVarP(&configPath, "config", "c", "", "Chemin vers le fichier de configuration JSON (requis)")
	analyzeCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Chemin pour exporter le rapport JSON")
	analyzeCmd.MarkFlagRequired("config")
}
