package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

func getValidFolderName() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	color.RGB(255, 255, 255).Print("Project Name: ")
	folderName, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}

	folderName = strings.TrimSpace(folderName)
	if folderName == "" {
		return "", fmt.Errorf("folder name cannot be empty")
	}

	if strings.ContainsAny(folderName, `\/:*?"<>|`) {
		return "", fmt.Errorf("folder name contains invalid characters")
	}

	return folderName, nil
}
