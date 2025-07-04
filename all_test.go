package cli_test

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestGolangCILint(t *testing.T) {
	rungo(t, "run", "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest", "run")
}

func TestGoModTidy(t *testing.T) {
	rungo(t, "mod", "tidy", "-diff")
}

func TestGovulncheck(t *testing.T) {
	rungo(t, "run", "golang.org/x/vuln/cmd/govulncheck@latest", "./...")
}

func TestGodocLint(t *testing.T) {
	rungo(t, "run", "github.com/godoc-lint/godoc-lint/cmd/godoclint@latest",
		"-exclude", "cmd/librarian/main.go",
		"-exclude", "internal/statepb",
		"./...")
}

func TestCoverage(t *testing.T) {
	rungo(t, "test", "-coverprofile=coverage.out", "./...")
}

func rungo(t *testing.T, args ...string) {
	t.Helper()

	cmd := exec.Command("go", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		if ee := (*exec.ExitError)(nil); errors.As(err, &ee) && len(ee.Stderr) > 0 {
			t.Fatalf("%v: %v\n%s", cmd, err, ee.Stderr)
		}
		t.Fatalf("%v: %v\n%s", cmd, err, output)
	}
}

func TestExportedSymbolsHaveDocs(t *testing.T) {
	err := filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") ||
			strings.HasSuffix(path, "_test.go") || strings.HasSuffix(path, ".pb.go") {
			return nil
		}

		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			t.Errorf("failed to parse file %q: %v", path, err)
			return nil
		}

		// Visit every top-level declaration in the file.
		for _, decl := range node.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if ok && (gen.Tok == token.TYPE || gen.Tok == token.VAR) {
				for _, spec := range gen.Specs {
					switch s := spec.(type) {
					case *ast.TypeSpec:
						checkDoc(t, s.Name, gen.Doc, path)
					case *ast.ValueSpec:
						for _, name := range s.Names {
							checkDoc(t, name, gen.Doc, path)
						}
					}
				}
			}
			if fn, ok := decl.(*ast.FuncDecl); ok {
				checkDoc(t, fn.Name, fn.Doc, path)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func checkDoc(t *testing.T, name *ast.Ident, doc *ast.CommentGroup, path string) {
	t.Helper()
	if !name.IsExported() {
		return
	}
	if doc == nil {
		t.Errorf("%s: %q is missing doc comment",
			path, name.Name)
	}
}
