package main

import (
	"context"
	"strings"
)

type Git struct{}

const (
	GIT_IMAGE = "alpine/git:2.43.0"
	WORK_DIR  = "/src/"
)

func (m *Git) GetLatestTag(ctx context.Context, source *Directory) (string, error) {
	args := []string{
		"describe", "--tags", "--abbrev=0",
	}
	tag, err := dag.Container().
		From(GIT_IMAGE).
		WithMountedDirectory(WORK_DIR, source).
		WithWorkdir(WORK_DIR).
		WithExec(args).
		Stdout(ctx)
	if err != nil {
		return "", err
	}
	tag = strings.Trim(tag, "\n")
	return tag, nil
}

func (m *Git) GetLatestBranch(ctx context.Context, source *Directory) (string, error) {
	args := []string{
		"rev-parse", "--abbrev-ref", "HEAD",
	}

	branch, err := dag.Container().
		From(GIT_IMAGE).
		WithMountedDirectory(WORK_DIR, source).
		WithWorkdir(WORK_DIR).
		WithExec(args).
		Stdout(ctx)
	if err != nil {
		return "", err
	}
	branch = strings.Trim(branch, "\n")
	return branch, nil
}

func (m *Git) GetHash(ctx context.Context, source *Directory) (string, error) {
	args := []string{
		"rev-parse", "HEAD",
	}

	hash, err := dag.Container().
		From(GIT_IMAGE).
		WithMountedDirectory(WORK_DIR, source).
		WithWorkdir(WORK_DIR).
		WithExec(args).
		Stdout(ctx)
	if err != nil {
		return "", err
	}
	hash = strings.Trim(hash, "\n")
	return hash, nil
}
