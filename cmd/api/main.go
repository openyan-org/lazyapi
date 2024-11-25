package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	codegen "github.com/openyan-org/lazyapi/codegen"
	lazyapi "github.com/openyan-org/lazyapi/core"
)

func main() {
	r := chi.NewRouter()

	r.Get("/", handler)

	log.Println("LazyAPI is listening on port 5000...")

	http.ListenAndServe(":5000", r)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("Failed to write JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, status int, err error) {
	log.Printf("Error: %v", err)
	http.Error(w, http.StatusText(status), status)
}

func handler(w http.ResponseWriter, r *http.Request) {
	api := lazyapi.NewAPI("Cats API", "go", "chi")
	api.SetDatabase("postgres")

	err := api.Validate()
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	fields := []lazyapi.Field{
		{
			Name: "id",
			Type: lazyapi.UUID,
			Constraints: lazyapi.FieldConstraints{
				Unique:   true,
				Required: true,
			},
		},
		{
			Name: "breed",
			Type: lazyapi.Text,
			Constraints: lazyapi.FieldConstraints{
				Unique:   false,
				Required: true,
			},
		},
		{
			Name: "name",
			Type: lazyapi.Text,
			Constraints: lazyapi.FieldConstraints{
				Unique:   false,
				Required: true,
			},
		},
	}

	catModel := lazyapi.NewModel("Cat", fields, []lazyapi.Relationship{})
	api.AddModel(catModel)

	createCatEndpoint := lazyapi.NewEndpoint("Post", "/cats")
	createCatEndpoint.SetBodySchema(catModel)
	createCatEndpoint.SetResponseSchema(struct {
		Message string `json:"message"`
	}{
		Message: "Successfully created a cat!",
	})

	createCatEndpoint.SetAction("insert_record")

	api.AddEndpoint(createCatEndpoint)

	src, err := codegen.GenerateSourceCode(api)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, src)
}
