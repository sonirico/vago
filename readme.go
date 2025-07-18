// Package main generates the README.md file for the vago project
// by extracting examples from test files and creating documentation.
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Constants and template definitions

const title = `
<div align="center">
  <img src="vago.png" alt="a golang gopher napping in the beach" title="a golang gopher napping in the beach" width="200"/>

  # vago
  
  The ultimate toolkit for vaGo developers. A comprehensive collection of functions, data structures, and utilities designed to enhance productivity and code quality with no learning curve and less effort.

  [![Go Report Card](https://goreportcard.com/badge/github.com/sonirico/vago)](https://goreportcard.com/report/github.com/sonirico/vago)
  [![Go Reference](https://pkg.go.dev/badge/github.com/sonirico/vago.svg)](https://pkg.go.dev/github.com/sonirico/vago)
  [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
  [![Release](https://img.shields.io/github/v/release/sonirico/vago.svg)](https://github.com/sonirico/vago/releases)
</div>

📖 **[View full documentation and examples on pkg.go.dev →](https://pkg.go.dev/github.com/sonirico/vago)**

## ✨ Workspace Architecture

This project leverages Go workspaces to provide **isolated dependencies** for each module. This means:

- 🎯 **Lightweight imports**: When you import ` + "`fp`" + ` or ` + "`streams`" + `, you won't download database drivers or logging dependencies
- 🔧 **Modular design**: Each module (` + "`db`" + `, ` + "`lol`" + `, ` + "`num`" + `) maintains its own ` + "`go.mod`" + ` with specific dependencies
- 📦 **Zero bloat**: Use only what you need without carrying unnecessary dependencies
- 🚀 **Fast builds**: Smaller dependency graphs lead to faster compilation and smaller binaries

**Example**: Importing ` + "`github.com/sonirico/vago/fp`" + ` will only pull functional programming utilities, not database connections or logging frameworks.

## Modules

`

var moduleEmojis = map[string]string{
	"ent":     "🪾",
	"fp":      "🪄",
	"maps":    "🗝️",
	"slices":  "⛓️",
	"streams": "🌊",
	"lol":     "📝",
	"num":     "🔢",
	"db":      "🗃️",
	"zero":    "🔞",
}

var moduleDescriptions = map[string]string{
	"ent":     "Environment variable management utilities with type-safe retrieval and validation.",
	"streams": "Powerful data streaming and processing utilities with fluent API for functional programming patterns.",
	"slices":  "Comprehensive slice manipulation utilities with functional programming patterns.",
	"maps":    "Map manipulation and transformation utilities.",
	"fp":      "Functional programming utilities including Option and Result types.",
	"lol":     "Structured logging utilities with multiple backends and APM integration.",
	"num":     "Numeric utilities including high-precision decimal operations.",
	"db":      "Database utilities and adapters for PostgreSQL, MongoDB, Redis, and ClickHouse.\n\n**Note:** This module always appears in the documentation, even if only interface or example tests are present, to ensure discoverability.",
	"zero":    "Zero-value utilities and string manipulation functions.",
}

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

// String generates the markdown representation of a module.
func (m mod) String() string {
	buf := bytes.NewBuffer(nil)

	// Module header with emoji
	moduleAnchor := strings.ToLower(m.title)
	emoji := moduleEmojis[m.title]
	if emoji == "" {
		emoji = "⚙️"
	}
	buf.WriteString(
		fmt.Sprintf(
			"## <a name=\"%s\"></a>%s %s\n\n",
			moduleAnchor,
			emoji,
			strings.ToUpper(m.title[:1])+m.title[1:],
		),
	)
	buf.WriteString(m.description)
	buf.WriteString("\n\n")

	// Table of contents
	buf.WriteString("### Functions\n\n")
	for _, fn := range m.funs {
		anchor := strings.ToLower(strings.ReplaceAll(m.title+" "+fn.name, " ", "-"))
		buf.WriteString(fmt.Sprintf("- [%s](#%s)\n", fn.name, anchor))
	}
	buf.WriteString("\n")

	// Function details
	for _, fn := range m.funs {
		buf.WriteString(fmt.Sprintf("#### %s %s\n\n", m.title, fn.name))
		buf.WriteString(fn.comment)

		buf.WriteString("\n\n<details><summary>Code</summary>\n\n")
		buf.WriteString("```go\n" + strings.TrimSpace(fn.body.String()) + "\n```\n\n</details>\n")
		buf.WriteString("\n\n[⬆️ Back to Top](#table-of-contents)\n\n---\n\n")
	}

	buf.WriteString("\n[⬆️ Back to Top](#table-of-contents)\n")

	return buf.String()
}

// cleanComment removes Go comment markers (//) and cleans up comment text.
func cleanComment(c string) string {
	return strings.ReplaceAll(c, "// ", "")
}

// createModuleFromExamples creates a module by parsing example functions from test files.
func createModuleFromExamples(moduleName string) *mod {
	m := &mod{
		title: moduleName,
		funs:  make([]fun, 0),
	}

	// Get module description from main module file if it exists
	mainFile := fmt.Sprintf("%s/%s.go", moduleName, moduleName)
	if _, err := os.Stat(mainFile); err == nil {
		if fset := token.NewFileSet(); fset != nil {
			if f, err := parser.ParseFile(fset, mainFile, nil, parser.ParseComments); err == nil &&
				f.Doc != nil {
				m.description = cleanComment(f.Doc.Text())
			}
		}
	}

	// Set default descriptions for known modules
	if m.description == "" {
		if desc, ok := moduleDescriptions[moduleName]; ok {
			m.description = desc
		}
	}

	// Find all test files in the module directory
	testPattern := fmt.Sprintf("%s/*_test.go", moduleName)
	testFiles, err := filepath.Glob(testPattern)
	if err != nil {
		return m
	}

	for _, filePath := range testFiles {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
		if err != nil {
			continue
		}

		// Extract Example functions
		for _, decl := range f.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok {
				if fn.Doc == nil || !strings.HasPrefix(fn.Name.String(), "Example") {
					continue
				}

				// Extract the function name without "Example" prefix
				funcName := strings.TrimPrefix(fn.Name.String(), "Example")
				if funcName == "" {
					continue
				}

				modFun := fun{
					name:    funcName,
					comment: cleanComment(fn.Doc.Text()),
					body:    new(bytes.Buffer),
				}

				// For example functions, extract source with comments preserved
				if err := extractSourceWithComments(modFun.body, fset, fn, filePath); err != nil {
					continue
				}

				m.funs = append(m.funs, modFun)
			}
		}
	}

	return m
}

// extractSourceWithComments extracts the source code of a function preserving comments
func extractSourceWithComments(
	buf *bytes.Buffer,
	fset *token.FileSet,
	fn *ast.FuncDecl,
	filePath string,
) error {
	// Read the source file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Get the position of the function in the source
	start := fset.Position(fn.Pos())
	end := fset.Position(fn.End())

	// Convert to lines
	lines := strings.Split(string(content), "\n")

	// Extract function lines (adjust for 1-based line numbers)
	if start.Line > len(lines) || end.Line > len(lines) {
		return fmt.Errorf("function position out of bounds")
	}

	// Extract the function body (skip the func declaration line, keep the body)
	bodyLines := []string{}
	inBody := false
	braceCount := 0

	for i := start.Line - 1; i < len(lines) && i < end.Line; i++ {
		line := lines[i]

		// Find opening brace to start body extraction
		if !inBody && strings.Contains(line, "{") {
			inBody = true
			// Extract everything after the opening brace
			bracePos := strings.Index(line, "{")
			if bracePos < len(line)-1 {
				bodyContent := line[bracePos+1:]
				if strings.TrimSpace(bodyContent) != "" {
					bodyLines = append(bodyLines, bodyContent)
				}
			}
			braceCount = 1
			continue
		}

		if inBody {
			// Count braces to know when we've reached the end
			openBraces := strings.Count(line, "{")
			closeBraces := strings.Count(line, "}")
			braceCount += openBraces - closeBraces

			// If we're at the closing brace, extract content before it
			if braceCount == 0 {
				if closeBracePos := strings.LastIndex(line, "}"); closeBracePos > 0 {
					bodyContent := line[:closeBracePos]
					bodyLines = append(bodyLines, bodyContent)
				}
				break
			} else {
				// Always add the line as-is to preserve comments and formatting
				bodyLines = append(bodyLines, line)
			}
		}
	}

	// Write the body content with proper indentation removed
	if len(bodyLines) > 0 {
		// Find common indentation to remove (ignore empty lines and pure comment lines)
		minIndent := 1000
		for _, line := range bodyLines {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" && !strings.HasPrefix(trimmed, "//") {
				indent := 0
				for _, char := range line {
					if char == ' ' || char == '\t' {
						indent++
					} else {
						break
					}
				}
				if indent < minIndent {
					minIndent = indent
				}
			}
		}

		// If no code lines found, use 0 as minimum indent
		if minIndent == 1000 {
			minIndent = 0
		}

		// Write function signature first
		buf.WriteString(fmt.Sprintf("func %s() {\n", fn.Name.Name))

		// Write body with adjusted indentation
		for _, line := range bodyLines {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" {
				// Empty line - preserve it
				buf.WriteString("\n")
			} else if len(line) > minIndent {
				// Remove common indentation and add one tab
				buf.WriteString("\t" + line[minIndent:] + "\n")
			} else {
				// Line is shorter than minIndent (probably a comment or short line)
				buf.WriteString("\t" + line + "\n")
			}
		}
		buf.WriteString("}")
	}

	return nil
}

// readme generates the complete README.md file for the vago project.
func readme() {
	modules := []string{
		"ent", "slices", "maps", "fp", "streams", "lol", "num", "db", "zero",
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

	// Create modules from examples
	mods := make([]*mod, 0, len(modules))
	for _, modName := range modules {
		fmt.Printf("Processing module '%s'...\n", modName)
		if exampleMod := createModuleFromExamples(modName); exampleMod != nil {
			if len(exampleMod.funs) > 0 {
				mods = append(mods, exampleMod)
				fmt.Printf("Found %d examples in '%s'\n", len(exampleMod.funs), modName)
			} else {
				fmt.Printf("No examples found in '%s'\n", modName)
			}
		}
	}

	// Sort modules alphabetically by title
	sort.Slice(mods, func(i, j int) bool {
		return mods[i].title < mods[j].title
	})

	// Global table of contents
	_, _ = buf.WriteString("## <a name=\"table-of-contents\"></a>Table of Contents\n\n")
	for _, m := range mods {
		m.sortFunctionsByName()
		emoji := moduleEmojis[m.title]
		if emoji == "" {
			emoji = "⚙️"
		}
		moduleAnchor := strings.ToLower(m.title)
		_, _ = buf.WriteString(
			fmt.Sprintf(
				"- [%s %s](#%s) - %d functions\n",
				emoji,
				strings.ToUpper(m.title[:1])+m.title[1:],
				moduleAnchor,
				len(m.funs),
			),
		)
	}
	_, _ = buf.WriteString("\n")

	// Module content
	for _, m := range mods {
		_, _ = buf.WriteString(m.String())
		_, _ = buf.WriteString("\n\n<br/>\n\n")
	}

	_ = buf.Flush()
}

func main() {
	readme()
}
