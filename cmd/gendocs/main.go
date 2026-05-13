// gendocs extracts the package-level doc comment from each module under
// internal/modules/ and writes it to docs/modules/<name>.md.
//
// The package doc comment is the canonical source of truth for module docs.
// Run via: make gen-docs
package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	repoRoot, err := findRepoRoot()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	modulesDir := filepath.Join(repoRoot, "internal", "modules")
	docsDir := filepath.Join(repoRoot, "docs", "modules")

	entries, err := os.ReadDir(modulesDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		doc, err := packageDoc(filepath.Join(modulesDir, name))
		if err != nil {
			fmt.Fprintf(os.Stderr, "skip %s: %v\n", name, err)
			continue
		}
		if doc == "" {
			fmt.Fprintf(os.Stderr, "skip %s: no package doc comment\n", name)
			continue
		}
		out := filepath.Join(docsDir, name+".md")
		if err := os.WriteFile(out, []byte(doc+"\n"), 0o644); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(out)
	}
}

// packageDoc returns the package-level doc comment for the first .go file in dir
// that has one, with comment markers stripped.
func packageDoc(dir string) (string, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)
	if err != nil {
		return "", err
	}
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			if file.Doc == nil {
				continue
			}
			return strings.TrimSpace(file.Doc.Text()), nil
		}
	}
	return "", nil
}

// findRepoRoot walks up from the binary's working directory looking for go.mod.
func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found")
		}
		dir = parent
	}
}
