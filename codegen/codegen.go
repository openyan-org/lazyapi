package codegen

import (
	"errors"

	"github.com/openyan-org/lazyapi/codegen/generators"
	lazyapi "github.com/openyan-org/lazyapi/core"
)

type APISourceCode struct {
	Language string      `json:"language"`
	Src      interface{} `json:"src"`
}

func GenerateSourceCode(api *lazyapi.API) (APISourceCode, error) {
	switch api.WebFramework {
	case "chi":
		src, err := generators.GenerateChi(*api)
		if err != nil {
			return APISourceCode{}, err
		}

		return APISourceCode{
			Language: "go",
			Src:      src,
		}, nil
	}

	return APISourceCode{}, errors.New("the given web framework is not supported")
}
