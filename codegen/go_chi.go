package codegen

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"

	lazyapi "github.com/openyan-org/lazyapi/core"
)

const mainTemplate = `// Generated by OpenYan LazyAPI

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"database/sql"
	_ "{{if eq .DatabaseEngine "postgres"}}github.com/lib/pq{{else if eq .DatabaseEngine "mysql"}}github.com/go-sql-driver/mysql{{end}}" // Database driver
)

func main() {
	err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	db, err := sql.Open("{{.DatabaseEngine}}", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Define routes
	{{range .Endpoints}}
	r.{{.Method}}("{{.Path}}", {{.BodySchema.Name}}Handler)
	{{end}}

	log.Println("Server is running at http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
`

const repositoriesTemplate = `// Generated by OpenYan LazyAPI

package main

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

{{range .Models}}
// {{.Name}} represents the {{.Name}} model
type {{.Name}} struct {
	{{range .Fields}}{{.Name}} {{if eq .Type "uuid"}}uuid.UUID{{else if eq .Type "text"}}string{{else}}{{.Type}}{{end}} ` + "`json:\"{{.Name}}\"`" + `
	{{end}}
}

// Insert{{.Name}} inserts a new {{.Name}} into the database
func Insert{{.Name}}(db *sql.DB, {{range .Fields}}{{.Name}} {{if eq .Type "uuid"}}uuid.UUID{{else if eq .Type "text"}}string{{else}}{{.Type}}{{end}}, {{end}}) error {
	query := "INSERT INTO {{.Name}} ({{range $index, $field := .Fields}}{{if $index}}, {{end}}{{$field.Name}}{{end}}) VALUES ({{range $index, $field := .Fields}}{{if $index}}, {{end}}?{{end}})"
	_, err := db.Exec(query, {{range .Fields}}{{.Name}}, {{end}})
	if err != nil {
		return fmt.Errorf("failed to insert {{.Name}}: %w", err)
	}
	return nil
}
{{end}}
`

const handlersTemplate = `// Generated by OpenYan LazyAPI

package main

import (
	"encoding/json"
	"net/http"
)

{{range .Endpoints}}
// {{.BodySchema.Name}}Handler handles {{.Method}} requests to {{.Path}}
func {{.BodySchema.Name}}Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var body {{.BodySchema.Name}}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Call repository function
		err := Insert{{.BodySchema.Name}}(nil /* DB */, {{range .BodySchema.Fields}}body.{{.Name}}, {{end}})
		if err != nil {
			http.Error(w, "Failed to insert record", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Created successfully"})
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
{{end}}
`

const envTemplate = `# Generated by OpenYan LazyAPI
# Enter your environment secrets here!

# Database URI configuration
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/postgres
`

func GenerateNetHTTP(api lazyapi.API) error {
	baseOutDir := "out"
	err := os.MkdirAll(baseOutDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create base output directory: %w", err)
	}

	nextDir, err := getNextOutputDirectory(baseOutDir)
	if err != nil {
		return fmt.Errorf("failed to determine next output directory: %w", err)
	}

	err = os.MkdirAll(nextDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create output directory %s: %w", nextDir, err)
	}

	files := map[string]string{
		filepath.Join(nextDir, "main.go"):         mainTemplate,
		filepath.Join(nextDir, "repositories.go"): repositoriesTemplate,
		filepath.Join(nextDir, "handlers.go"):     handlersTemplate,
		filepath.Join(nextDir, ".env"):            envTemplate,
	}

	lazyapiFile := filepath.Join(nextDir, "lazyapi.json")
	fmt.Println("Generating lazyapi.json...")
	err = writeJSON(lazyapiFile, api)
	if err != nil {
		return fmt.Errorf("failed to write lazyapi.json: %w", err)
	}

	fmt.Println("Initializing Go module...")
	err = runCommand("go mod init lazyapi_api", nextDir)
	if err != nil {
		return fmt.Errorf("failed to initialize Go module: %w", err)
	}

	fmt.Println("Generating files...")

	for fileName, tmplContent := range files {
		err := writeTemplate(fileName, tmplContent, api)
		if err != nil {
			return fmt.Errorf("failed to write %s: %w", fileName, err)
		}
	}

	fmt.Println("Installing dependencies...")
	dependencies := []string{
		"github.com/google/uuid",
		"github.com/go-chi/chi/v5",
		"github.com/joho/godotenv",
	}
	if api.DatabaseEngine == "postgres" {
		dependencies = append(dependencies, "github.com/lib/pq")
	} else if api.DatabaseEngine == "mysql" {
		dependencies = append(dependencies, "github.com/go-sql-driver/mysql")
	}

	for _, dep := range dependencies {
		err = runCommand(fmt.Sprintf("go get %s", dep), nextDir)
		if err != nil {
			return fmt.Errorf("failed to add dependency %s: %w", dep, err)
		}
	}

	fmt.Println("Tidying up dependencies...")
	err = runCommand("go mod tidy", nextDir)
	if err != nil {
		return fmt.Errorf("failed to tidy modules: %w", err)
	}

	fmt.Printf("Source code generated in directory: %s\n", nextDir)
	return nil
}

func getNextOutputDirectory(baseOutDir string) (string, error) {
	entries, err := os.ReadDir(baseOutDir)
	if err != nil {
		return "", err
	}

	var serials []int
	for _, entry := range entries {
		if entry.IsDir() {
			parts := strings.Split(entry.Name(), "_")
			if len(parts) > 0 {
				if serial, err := strconv.Atoi(parts[0]); err == nil {
					serials = append(serials, serial)
				}
			}
		}
	}

	sort.Ints(serials)
	nextSerial := 1
	if len(serials) > 0 {
		nextSerial = serials[len(serials)-1] + 1
	}

	for {
		nextDir := fmt.Sprintf("%s/%04d_lazyapi_src", baseOutDir, nextSerial)
		if _, err := os.Stat(nextDir); os.IsNotExist(err) {
			return nextDir, nil
		}
		nextSerial++
	}
}

func writeTemplate(fileName, tmplContent string, data interface{}) error {
	tmpl, err := template.New(fileName).Parse(tmplContent)
	if err != nil {
		return err
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

func runCommand(cmd, dir string) error {
	parts := strings.Split(cmd, " ")
	c := exec.Command(parts[0], parts[1:]...)
	c.Dir = dir
	c.Stdout = nil
	c.Stderr = nil
	return c.Run()
}

func writeJSON(fileName string, data interface{}) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}