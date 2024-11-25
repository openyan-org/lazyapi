package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/openyan-org/lazyapi/codegen/generators"
)

func writeFilesGo(baseDir string, goSrc generators.GoAPI) error {
	err := os.MkdirAll(baseDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	files := map[string]string{
		"main.go":         goSrc.MainFile,
		"repositories.go": goSrc.RepositoriesFile,
		"handlers.go":     goSrc.HandlersFile,
		".env":            goSrc.EnvFile,
		"setup.sh":        goSrc.SetupScriptFile,
		"lazyapi.json":    goSrc.LazyAPIFile,
	}

	for name, content := range files {
		filePath := filepath.Join(baseDir, name)
		err := os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			return fmt.Errorf("failed to write file %s: %w", filePath, err)
		}
	}

	return nil
}
