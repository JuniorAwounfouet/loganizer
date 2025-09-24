package reporter

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/JuniorAwounfouet/go_loganizer/internal/analyzer"
)

func ExportResults(results []analyzer.AnalysisResult, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(results)
}

func PrintResults(results []analyzer.AnalysisResult) {
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
}
