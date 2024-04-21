// Golang build tools
package main

import (
	"context"
	"fmt"
	"runtime"
	"strings"
)

type Golang struct{}

const (
	// TODO: Make these dynamic
	GO_BASE_IMAGE   = "golang:1.21-bookworm"
	LINT_BASE_IMAGE = "golangci/golangci-lint:v1.57.2"

	OUT_DIR    = "/out/"
	SOURCE_DIR = "/src/"
)

// Sets up the base container for re-use across all Functions in this module
func base(ctx context.Context, baseImage string) *Container {
	ctr := dag.Container().From(baseImage)
	goCacheEnv, _ := ctr.WithExec([]string{"go", "env", "GOCACHE"}).Stdout(ctx)
	goModCacheEnv, _ := ctr.WithExec([]string{"go", "env", "GOMODCACHE"}).Stdout(ctx)

	gomod := dag.CacheVolume("gomod")
	gobuild := dag.CacheVolume("gobuild")

	return ctr.
		WithMountedCache(goModCacheEnv, gomod).
		WithMountedCache(goCacheEnv, gobuild)
}

func (m *Golang) Build(
	// The packages to build within this go project
	// +optional
	// +default ./...
	packagePath string,

	// GOOS envvar
	// +optional
	os string,

	// GOARCH envvar
	// +optional
	arch string,

	// Source code of the go project
	// +optional
	source *Directory,

	// LDFlags to pass to build
	// +optional
	ldFlags []string,
) *Directory {
	ctx := context.Background()

	if arch == "" {
		arch = runtime.GOARCH
	}
	if os == "" {
		os = runtime.GOOS
	}

	// Default to populating this stuff from the git directory
	if ldFlags == nil {
		gitHash, err := dag.Git().GetHash(ctx, source)
		if err != nil {
			gitHash = "unknown"
		}
		gitBranch, err := dag.Git().GetLatestBranch(ctx, source)
		if err != nil {
			gitBranch = "unknown"
		}
		gitTag, err := dag.Git().GetLatestTag(ctx, source)
		if err != nil {
			gitTag = "unknown"
		}
		ldFlags = []string{
			fmt.Sprintf("-X main.gitHash=%s", gitHash),
			fmt.Sprintf("-X main.gitBranch=%s", gitBranch),
			fmt.Sprintf("-X main.version=%s", gitTag),
		}
	}

	ldFlagsString := strings.Join(ldFlags, " ")
	command := []string{"go", "build", "-ldflags", ldFlagsString, "-o", OUT_DIR, packagePath}
	return base(ctx, GO_BASE_IMAGE).
		WithMountedDirectory(SOURCE_DIR, source).
		WithWorkdir(SOURCE_DIR).
		WithExec(command).
		Directory(OUT_DIR)
}

func (m *Golang) Test(
	// The source code to run test on
	source *Directory,

	// Which packages to run test against.
	// +optional
	// +default="./..."
	packages string,

	// Arguments to go test
	// +optional
	// +default=["-short", "-shuffle=on"]
	args []string,

	// Filter tests to run via regex
	// +optional
	run string,

	// Skip tests via regex
	// +optional
	skip string,
) (string, error) {
	ctx := context.Background()
	command := []string{"go", "test"}
	if run != "" {
		command = append(command, "-run", run)
	}
	if skip != "" {
		command = append(command, "-skip", skip)
	}
	command = append(command, args...)
	command = append(command, packages)
	return base(ctx, GO_BASE_IMAGE).
		WithMountedDirectory(SOURCE_DIR, source).
		WithWorkdir(SOURCE_DIR).
		WithExec(command).
		Stdout(ctx)
}

func (m *Golang) Lint(
	// The source code to lint
	source *Directory,

	// Format to use
	// +optional
	// +default="colored-line-number"
	format string,
) (string, error) {
	ctx := context.Background()
	command := []string{
		"golangci-lint",
		"run",
		"--out-format",
		format,
	}
	return base(ctx, LINT_BASE_IMAGE).
		WithMountedDirectory(SOURCE_DIR, source).
		WithWorkdir(SOURCE_DIR).
		WithExec(command).
		Stdout(ctx)
}
