package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/openyan-org/lazyapi/codegen"
	"github.com/openyan-org/lazyapi/codegen/generators"
)

func main() {
	if len(os.Args) != 2 {
		color.Red("Error: You must pass in at least 1 argument.")
		return
	}

	apiSpecURL := os.Args[1]

	resp, err := http.Get(apiSpecURL)
	if err != nil {
		color.Red("Error: Failed to fetch API spec: %v\n", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		color.Red("Error: received status code %d\n", resp.StatusCode)
		return
	}

	var apiPayload codegen.APISourceCode
	err = json.NewDecoder(resp.Body).Decode(&apiPayload)
	if err != nil {
		color.Red("Error: Failed to decode JSON: %v\n", err)
		return
	}

	switch apiPayload.Language {
	case "go":
		var goSrc generators.GoAPI
		srcData, err := json.Marshal(apiPayload.Src)
		if err != nil {
			color.Red("Error: Failed to marshal Src: %v\n", err)
			os.Exit(1)
		}

		err = json.Unmarshal(srcData, &goSrc)
		if err != nil {
			color.Red("Error: Failed to unmarshal Src to GoAPI: %v\n", err)
			os.Exit(1)
			return
		}

		folderName, err := getValidFolderName()
		if err != nil {
			color.Red("Error: %v\n", err)
			os.Exit(1)
			return
		}

		err = writeFilesGo(folderName, goSrc)
		if err != nil {
			color.Red("Error generating files: %v\n", err)
			os.Exit(1)
			return
		}

		successfulFeedback(folderName)
	default:
		color.Red("Error: Invalid language for LazyAPI source code.")
		os.Exit(1)
	}

}

func successfulFeedback(folderName string) {
	color.RGB(107, 114, 200).Println("\n------------ [OpenYan LazyAPI] ------------")
	color.Green("Successfully generated source code.")

	color.RGB(107, 114, 128).Println("\nGet started by running:")
	color.Magenta("$ cd %s", folderName)
	color.Magenta("$ bash setup.sh")
	color.RGB(107, 114, 200).Println("------------ [OpenYan LazyAPI] ------------")
}
