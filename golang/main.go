// Golang build tools
package main

import (
	"context"
	"runtime"
)

type Golang struct{}

// TODO: Make this dynamic
const (
	BASE_IMAGE = "golang:1.21-bookworm"
	OUT_DIR    = "/out/"
	SOURCE_DIR = "/src/"
)

// Sets up the base container for re-use across all Functions in this module
func base(ctx context.Context) *Container {
	ctr := dag.Container().From(BASE_IMAGE)
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
	if arch == "" {
		arch = runtime.GOARCH
	}
	if os == "" {
		os = runtime.GOOS
	}

	// // Default to populating this stuff from the git directory
	// if ldFlags == []string(nil){
	//   ldFlags = []string{
	//     fmt.Sprintf("-X main.gitHash=%s", )
	//   }
	// }

	command := append([]string{"go", "build", "-o", OUT_DIR, packagePath})
	ctx := context.Background()
	ctr := base(ctx).
		WithMountedDirectory(SOURCE_DIR, source).
		WithWorkdir(SOURCE_DIR).
		WithExec(command)

	return ctr.Directory(OUT_DIR)
}
