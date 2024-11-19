package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/openyan-org/lazyapi"
)

func main() {
	r := chi.NewRouter()

	r.Get("/", handler)

	log.Println("LazyAPI testing client is running at http://localhost:4123")

	http.ListenAndServe(":4123", r)
}

func handler(w http.ResponseWriter, r *http.Request) {
	api := lazyapi.NewAPI("My API", "go", "net/http")
	api.SetDatabase("postgresql")
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

	userModel := lazyapi.NewModel("User", userFields, []lazyapi.Relationship{})
	api.AddModel(userModel)

	createUserEndpoint := lazyapi.NewEndpoint("Post", "/users")
	createUserEndpoint.SetBodySchema(userModel)
	createUserEndpoint.SetResponseSchema(struct {
		user_id int
	}{
		user_id: 1,
	})
	createUserEndpoint.SetAction("insert_record")

	api.AddEndpoint(createUserEndpoint)

	err := api.GenerateSourceCode()
	if err != nil {
		panic(err)
	}

	writeJSON(w, 200, api)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}
