package analyzer

import (
	"errors"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	ErrFileNotFound    = errors.New("file not found")
	ErrFileNotReadable = errors.New("file not readable")
	ErrParsingFailed   = errors.New("parsing failed")
)

type AnalysisResult struct {
	LogID        string `json:"log_id"`
	FilePath     string `json:"file_path"`
	Status       string `json:"status"`
	Message      string `json:"message"`
	ErrorDetails string `json:"error_details,omitempty"`
}

func AnalyzeLog(logID, filePath, logType string) AnalysisResult {
	result := AnalysisResult{
		LogID:    logID,
		FilePath: filePath,
	}

	// Vérifier si le fichier existe
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		result.Status = "FAILED"
		result.Message = "Fichier introuvable."
		result.ErrorDetails = err.Error()
		return result
	}

	// Vérifier si le fichier est lisible
	file, err := os.Open(filePath)
	if err != nil {
		result.Status = "FAILED"
		result.Message = "Fichier inaccessible."
		result.ErrorDetails = err.Error()
		return result
	}
	defer file.Close()

	// Vérifier si le fichier est vide
	fileInfo, err := file.Stat()
	if err != nil {
		result.Status = "FAILED"
		result.Message = "Impossible de lire les informations du fichier."
		result.ErrorDetails = err.Error()
		return result
	}

	if fileInfo.Size() == 0 {
		result.Status = "OK"
		result.Message = "Fichier vide - aucune analyse nécessaire."
		return result
	}

	// Simuler l'analyse avec un délai aléatoire
	rand.Seed(time.Now().UnixNano())
	delay := time.Duration(rand.Intn(151) + 50)
	time.Sleep(delay * time.Millisecond)

	// Vérifier le contenu pour détecter les erreurs de parsing
	content, err := os.ReadFile(filePath)
	if err != nil {
		result.Status = "FAILED"
		result.Message = "Erreur lors de la lecture du fichier."
		result.ErrorDetails = err.Error()
		return result
	}

	// Détection simple d'erreurs de parsing
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.Contains(line, "INVALID_LINE") || strings.Contains(line, "cannot be parsed") {
			result.Status = "FAILED"
			result.Message = "Erreur de parsing détectée."
			result.ErrorDetails = ErrParsingFailed.Error() + ": " + line
			return result
		}
	}

	result.Status = "OK"
	result.Message = "Analyse terminée avec succès."
	return result
}

func AnalyzeLogsConcurrently(configs []struct {
	ID   string
	Path string
	Type string
}) []AnalysisResult {
	results := make([]AnalysisResult, len(configs))
	resultChan := make(chan AnalysisResult, len(configs))

	for _, config := range configs {
		go func(logID, filePath, logType string) {
			result := AnalyzeLog(logID, filePath, logType)
			resultChan <- result
		}(config.ID, config.Path, config.Type)
	}

	for i := 0; i < len(configs); i++ {
		results[i] = <-resultChan
	}

	return results
}

func HandleError(err error) string {
	if errors.Is(err, ErrFileNotFound) {
		return "Erreur: Fichier non trouvé"
	}
	if errors.Is(err, ErrFileNotReadable) {
		return "Erreur: Fichier non lisible"
	}
	if errors.Is(err, ErrParsingFailed) {
		return "Erreur: Échec du parsing"
	}
	return "Erreur inconnue: " + err.Error()
}
