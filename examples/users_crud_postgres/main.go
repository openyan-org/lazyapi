package main

import (
	"github.com/openyan-org/lazyapi"
)

func main() {
	api := lazyapi.NewAPI("Users API", "go", "net/http")
	api.SetDatabase("postgres")
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
}
