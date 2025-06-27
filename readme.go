package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const title = `
# gozo

[![Go Report Card](https://goreportcard.com/badge/github.com/sonirico/gozo)](https://goreportcard.com/report/github.com/sonirico/gozo)
[![Go Reference](https://pkg.go.dev/badge/github.com/sonirico/gozo.svg)](https://pkg.go.dev/github.com/sonirico/gozo)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

[![gozo Art](gozo.png)](https://github.com/sonirico/gozo/gozo.png)

The ultimate toolkit for Go developers. A comprehensive collection of functions, data structures, and utilities designed to enhance productivity and code quality.

## Modules

`

type (
	fun struct {
		name    string
		comment string
		body    *bytes.Buffer
	}

	mod struct {
		title       string
		description string
		funs        []fun
	}
)

// sortFunctionsByName sorts the functions within a module alphabetically by name.
func (m *mod) sortFunctionsByName() {
	sort.Slice(m.funs, func(i, j int) bool {
		return m.funs[i].name < m.funs[j].name
	})
}

func (m mod) String() string {
	buf := bytes.NewBuffer(nil)

	// Module header with Emoji
	// Use lowercase anchor link for the module itself
	moduleAnchor := strings.ToLower(m.title)
	emoji := moduleEmojis[m.title]
	if emoji == "" {
		emoji = "⚙️" // Default emoji
	}
	buf.WriteString(fmt.Sprintf("## <a name=\"%s\"></a>%s %s\n\n", moduleAnchor, emoji, strings.ToUpper(m.title[:1])+m.title[1:]))
	buf.WriteString(m.description)
	buf.WriteString("\n\n")

	// Module-specific table of contents
	buf.WriteString("### Functions\n\n")
	for _, fn := range m.funs {
		// Use lowercase anchor links consistent with global ToC
		buf.WriteString(fmt.Sprintf("- [%s](#%s)\n", fn.name, strings.ToLower(fn.name)))
	}
	buf.WriteString("\n")

	// Function details
	for _, fn := range m.funs {
		buf.WriteString(fn.String())
		// Use a standard markdown horizontal rule
		buf.WriteString("\n\n---\n\n")
	}

	// Add Back to Top link
	buf.WriteString("\n[⬆️ Back to Top](#table-of-contents)\n")

	return buf.String()
}

func (f fun) String() string {
	buf := bytes.NewBuffer(nil)

	// Function header (using H4 for better hierarchy)
	// Anchor link for the function name itself
	buf.WriteString(fmt.Sprintf("#### <a name=\"%s\"></a>%s\n\n", strings.ToLower(f.name), f.name))
	buf.WriteString(f.comment)

	// Code block with collapsible details
	buf.WriteString("\n\n<details><summary>Code</summary>\n\n")
	// Ensure consistent code block formatting
	buf.WriteString("```go\n" + strings.TrimSpace(f.body.String()) + "\n```\n\n</details>\n")

	return buf.String()
}

func cleanComment(c string) string {
	return strings.ReplaceAll(c, "// ", "")
}

var moduleEmojis = map[string]string{
	"slices":  "⛓️", // Chain/Sequence
	"maps":    "🗝️", // Keys
	"fp":      "🪄",  // Magic/Transformation
	"streams": "🌊",  // Water/Flow
}

func readme() {
	modules := []string{
		"slices", "maps", "fp", "streams",
	}

	// Open README.md for writing (truncate if exists)
	file, err := os.Create("README.md")
	if err != nil {
		fmt.Printf("Error creating README.md: %v\n", err)
		return
	}
	defer file.Close()

	buf := bufio.NewWriter(file)
	_, _ = buf.WriteString(title)

	// Use map to aggregate functions by package name
	modsMap := make(map[string]*mod)

	for _, modName := range modules {
		_ = filepath.WalkDir(modName, func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() || strings.Contains(path, "_test") || !strings.HasSuffix(path, ".go") {
				return nil
			}

			fset := token.NewFileSet()

			fmt.Printf("parsing '%s'...\n", path) // Debugging log
			// Parse src but stop after processing the imports.
			f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				fmt.Printf("Error parsing %s: %v\n", path, err) // Log parsing errors
				return nil                                      // Continue with other files
			}
			fmt.Println("...done") // Debugging log

			pkgName := f.Name.Name
			m, exists := modsMap[pkgName]
			if !exists {
				m = &mod{
					title:       pkgName,
					description: cleanComment(f.Doc.Text()),
					funs:        make([]fun, 0),
				}
				modsMap[pkgName] = m
			}

			// Add functions from this file
			for _, decl := range f.Decls {
				if fn, ok := decl.(*ast.FuncDecl); ok {
					if fn.Doc == nil || fn.Recv != nil || !fn.Name.IsExported() {
						continue
					}

					modFun := fun{
						name:    fn.Name.String(),
						comment: cleanComment(fn.Doc.Text()),
						body:    new(bytes.Buffer),
					}

					// remove comments from node in order to make go/format print the body without comments
					fn.Doc = nil

					if err = format.Node(modFun.body, fset, fn); err != nil {
						fmt.Printf("Error formatting node for %s in %s: %v\n", modFun.name, path, err) // Log formatting errors
						continue
					}

					m.funs = append(m.funs, modFun)
				}
			}

			return nil
		})
	}

	// Convert map to slice for sorting
	mods := make([]*mod, 0, len(modsMap))
	for _, m := range modsMap {
		mods = append(mods, m)
	}

	// Sort modules alphabetically by title
	sort.Slice(mods, func(i, j int) bool {
		return mods[i].title < mods[j].title
	})

	// Global table of contents
	// Add an anchor for the Table of Contents itself
	_, _ = buf.WriteString("## <a name=\"table-of-contents\"></a>Table of Contents\n\n")
	for _, m := range mods {
		// Sort functions within each module alphabetically before generating ToC
		m.sortFunctionsByName() // Renamed function call
		emoji := moduleEmojis[m.title]
		if emoji == "" {
			emoji = "⚙️" // Default emoji
		}
		// Use lowercase anchor links for modules
		moduleAnchor := strings.ToLower(m.title)
		_, _ = buf.WriteString(fmt.Sprintf("- [%s %s](#%s)\n", emoji, strings.ToUpper(m.title[:1])+m.title[1:], moduleAnchor))
		for _, fn := range m.funs {
			// Use lowercase anchor links for functions
			functionAnchor := strings.ToLower(fn.name)
			_, _ = buf.WriteString(fmt.Sprintf("  - [%s](#%s)\n", fn.name, functionAnchor))
		}
	}
	_, _ = buf.WriteString("\n")

	// Module content
	for _, m := range mods {
		// Functions are already sorted from the ToC generation step
		_, _ = buf.WriteString(m.String())
		// Add a visual break between modules
		_, _ = buf.WriteString("\n\n<br/>\n\n")
	}

	_ = buf.Flush()
}

func main() {
	readme()
}
