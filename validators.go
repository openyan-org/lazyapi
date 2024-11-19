package lazyapi

import (
	"errors"
	"fmt"
)

func (api *API) Validate() error {
	if api.PackageName == "" {
		return errors.New("PackageName is required")
	}

	switch api.Language {
	case Go, Python, TypeScript:
		// language is valid
	default:
		return fmt.Errorf("invalid Language: %s", api.Language)
	}

	switch api.WebFramework {
	case NetHTTP, Flask, Hono:
		// framework is valid
	default:
		return fmt.Errorf("invalid WebFramework: %s", api.WebFramework)
	}

	switch api.DatabaseEngine {
	case PostgreSQL, MySQL, SQLite:
		// database engine is valid
	default:
		return fmt.Errorf("invalid DatabaseEngine: %s", api.DatabaseEngine)
	}

	if err := ValidatePathPrefix(api.PathPrefix); err != nil {
		return fmt.Errorf("invalid PathPrefix: %v", err)
	}

	return nil
}

func ValidatePathPrefix(path string) error {
	if path == "" {
		return nil
	}

	// must start with '/'
	if path[0] != '/' {
		return errors.New("PathPrefix must start with '/'")
	}

	// must not contain "//"
	if containsDoubleSlashes(path) {
		return errors.New("PathPrefix must not contain '//'")
	}

	// must not end with a trailing slash unless it's the root "/"
	if len(path) > 1 && path[len(path)-1] == '/' {
		return errors.New("PathPrefix must not end with a trailing slash")
	}

	return nil
}

func containsDoubleSlashes(path string) bool {
	for i := 1; i < len(path); i++ {
		if path[i] == '/' && path[i-1] == '/' {
			return true
		}
	}
	return false
}
