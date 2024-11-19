package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/openyan-org/lazyapi"
)

func main() {
	r := chi.NewRouter()

	r.Get("/api/test", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	api := lazyapi.NewAPI("My API", "go", "net/http")
	api.SetDatabase("postgresql")
	api.SetPathPrefix("/api")
	api.Validate()

	userFields := []lazyapi.Field{
		{
			Name: "id",
			Type: lazyapi.UUID,
			Constraints: lazyapi.FieldConstraints{
				Unique:   true,
				Required: true,
			},
		},
		{
			Name: "first_name",
			Type: lazyapi.Text,
			Constraints: lazyapi.FieldConstraints{
				Unique:   false,
				Required: true,
			},
		},
	}

	userModel := lazyapi.NewModel("user", userFields, []lazyapi.Relationship{})
	api.AddModel(userModel)
}
