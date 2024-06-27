package architecture_test

import (
	"testing"

	"golang.org/x/tools/go/packages"
)

// checkImports tests if certain packages only import allowed packages.
func checkImports(t *testing.T, pkgName string, disallowedImports []string) {
	cfg := &packages.Config{
		Mode: packages.NeedImports,
	}
	pkgs, err := packages.Load(cfg, pkgName)
	if err != nil {
		t.Fatalf("failed to load packages: %v", err)
	}

	for _, pkg := range pkgs {
		for impPath := range pkg.Imports {
			for _, disallowed := range disallowedImports {
				if impPath == disallowed {
					t.Errorf("%s should not import %s", pkgName, disallowed)
				}
			}
		}
	}
}

func TestHandlerDependencies(t *testing.T) {
	disallowed := []string{
		"github.com/vladyslavpavlenko/genesis-api-project/internal/dbrepo",
		"github.com/vladyslavpavlenko/genesis-api-project/internal/dbrepo/gormrepo",
		"github.com/vladyslavpavlenko/genesis-api-project/internal/rateapi",
		"github.com/vladyslavpavlenko/genesis-api-project/internal/rateapi/chain",
	}
	checkImports(t, "github.com/vladyslavpavlenko/genesis-api-project/internal/handlers", disallowed)
}
