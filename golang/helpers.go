package main

import (
	"context"

	"golang.org/x/mod/modfile"
)

// Parses the go.mod file to get the version of Go for this project.
func inspectModVersion(ctx context.Context, src *Directory) (string, error) {
	goMod := "go.mod"
	mod, err := src.File(goMod).Contents(ctx)
	if err != nil {
		return "", err
	}

	f, err := modfile.Parse(goMod, []byte(mod), nil)
	if err != nil {
		return "", err
	}
	return f.Go.Version, nil
}
