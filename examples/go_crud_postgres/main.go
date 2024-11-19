package main

import (
	"log"

	"github.com/openyan-org/lazyapi/codegen"
	lazyapi "github.com/openyan-org/lazyapi/core"
)

func main() {
	api := lazyapi.NewAPI("Cats API", "go", "chi")
	api.SetDatabase("postgres")

	// Always validate the inputs.
	err := api.Validate()
	if err != nil {
		log.Fatal(err)
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
		message string
	}{
		message: "Successfully created a cat!",
	})

	createCatEndpoint.SetAction("insert_record")

	api.AddEndpoint(createCatEndpoint)

	err = codegen.GenerateSourceCode(api)
	if err != nil {
		log.Fatal(err)
	}
}
