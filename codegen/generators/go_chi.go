package generators

import (
	"fmt"

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

func GenerateChi(api lazyapi.API) (*GoAPI, error) {
	mainFile, err := renderTemplate(mainTemplate, api)
	if err != nil {
		return nil, fmt.Errorf("failed to render main.go template: %w", err)
	}

	repositoriesFile, err := renderTemplate(repositoriesTemplate, api)
	if err != nil {
		return nil, fmt.Errorf("failed to render repositories.go template: %w", err)
	}

	handlersFile, err := renderTemplate(handlersTemplate, api)
	if err != nil {
		return nil, fmt.Errorf("failed to render handlers.go template: %w", err)
	}

	envFile, err := renderTemplate(envTemplate, api)
	if err != nil {
		return nil, fmt.Errorf("failed to render .env template: %w", err)
	}

	lazyAPIFile, err := generateJSON(api)
	if err != nil {
		return nil, fmt.Errorf("failed to generate lazyapi.json: %w", err)
	}

	return &GoAPI{
		MainFile:         mainFile,
		RepositoriesFile: repositoriesFile,
		HandlersFile:     handlersFile,
		EnvFile:          envFile,
		LazyAPIFile:      lazyAPIFile,
	}, nil
}
