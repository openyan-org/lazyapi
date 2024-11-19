package codegen

import (
	"errors"

	lazyapi "github.com/openyan-org/lazyapi/core"
)

func GenerateSourceCode(api *lazyapi.API) error {
	switch api.WebFramework {
	case "net/http":
		return GenerateNetHTTP(*api)
	}

	return errors.New("the given web frammework is not supported")
}