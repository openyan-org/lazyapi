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
	api := lazyapi.NewAPI("My API", lazyapi.Go, lazyapi.NetHTTP)
	api.SetDatabase(lazyapi.PostgreSQL)
	api.SetPathPrefix("/api")

	// userModel := lazyapi.NewModel("user", userFields, []lazyapi.Relationship{})
}
