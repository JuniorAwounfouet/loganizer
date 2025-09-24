# Loganalyzer

Outil d'analyse de fichiers de logs développé en Go pour aider les administrateurs système à centraliser l'analyse de multiples logs en parallèle.

## Fonctionnalités

- Analyse concurrentielle de multiples fichiers de logs
- Gestion robuste des erreurs (fichiers introuvables, parsing, etc.)
- Configuration via fichier JSON
- Export des résultats au format JSON
- Interface CLI intuitive avec Cobra

## Installation

```bash
git clone du repo
cd go_loganizer
go build -o loganalyzer