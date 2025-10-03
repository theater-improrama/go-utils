// enumvalidator generates Validate() methods for enum types
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
	"text/template"
)

var (
	typeNames = flag.String("type", "", "comma-separated list of type names; must be set")
	output    = flag.String("output", "", "output file name; default srcdir/<type>_validation.go")
)

// TypeInfo holds information about a single enum type
type TypeInfo struct {
	TypeName  string
	Constants []string
	LowerType string
}

// TemplateData holds the data for the validation template
type TemplateData struct {
	PackageName string
	Types       []TypeInfo
}

const validationTemplate = `// Code generated github.com/theater-improrama/go-utils/tools/enumvalidator DO NOT EDIT.
package {{ .PackageName }}

import (
	"errors"
	"github.com/theater-improrama/go-utils/validator"
)

{{- range .Types }}

var _ validator.Validator = (*{{ .TypeName }})(nil)

// ErrInvalid{{ .TypeName }} is returned when an invalid value is passed to Validate()
var ErrInvalid{{ .TypeName }} = errors.New("invalid {{ .TypeName }}")

// Validate returns an error if the enum value is invalid
func (e {{ .TypeName }}) Validate() error {
	switch e {
{{- range .Constants }}
	case {{ . }}:
{{- end }}
	default:
		return ErrInvalid{{ .TypeName }}
	}

	return nil
}
{{- end }}
`

func main() {
	flag.Parse()
	if len(*typeNames) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	// Parse the package in the current directory
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, ".", nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	// Find the package (should be only one)
	var pkg *ast.Package
	for _, p := range pkgs {
		if !strings.HasSuffix(p.Name, "_test") {
			pkg = p
			break
		}
	}
	if pkg == nil {
		log.Fatal("no package found")
	}

	types := strings.Split(*typeNames, ",")
	if err := generateValidations(pkg, types, fset); err != nil {
		log.Fatalf("generating validations: %v", err)
	}
}

func generateValidations(pkg *ast.Package, typeNames []string, fset *token.FileSet) error {
	var packageName string
	var sourceFileName string
	var typeInfos []TypeInfo

	// Collect all type information
	for _, typeName := range typeNames {
		typeName = strings.TrimSpace(typeName)
		typeInfo, srcFile, pkgName, err := findTypeInfo(pkg, typeName, fset)
		if err != nil {
			return fmt.Errorf("finding type %s: %v", typeName, err)
		}
		typeInfos = append(typeInfos, typeInfo)
		if packageName == "" {
			packageName = pkgName
		}
		if sourceFileName == "" {
			sourceFileName = srcFile
		}
	}

	// Generate all validations in a single file
	tmpl, err := template.New("validation").Parse(validationTemplate)
	if err != nil {
		return fmt.Errorf("parsing template: %v", err)
	}

	data := TemplateData{
		PackageName: packageName,
		Types:       typeInfos,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("executing template: %v", err)
	}

	// Format the generated code
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("formatting generated code: %v", err)
	}

	// Determine output filename based on source file
	outputFile := *output
	if outputFile == "" {
		base := strings.TrimSuffix(sourceFileName, ".go")
		outputFile = fmt.Sprintf("%s_enumvalidator.go", base)
	}

	// Write to file
	if err := os.WriteFile(outputFile, formatted, 0644); err != nil {
		return fmt.Errorf("writing output file: %v", err)
	}

	fmt.Printf("Generated validations for %v in %s\n", typeNames, outputFile)
	return nil
}

func findTypeInfo(pkg *ast.Package, typeName string, fset *token.FileSet) (TypeInfo, string, string, error) {
	// Find the type declaration and its constants
	var typeDecl *ast.TypeSpec
	var constants []string
	var packageName string
	var sourceFileName string

	for _, file := range pkg.Files {
		if packageName == "" {
			packageName = file.Name.Name
		}

		ast.Inspect(file, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.TypeSpec:
				if node.Name.Name == typeName {
					typeDecl = node
					// Extract just the filename from the full path
					pos := fset.Position(node.Pos())
					sourceFileName = strings.TrimPrefix(pos.Filename, "./")
					if idx := strings.LastIndex(sourceFileName, "/"); idx >= 0 {
						sourceFileName = sourceFileName[idx+1:]
					}
				}
			case *ast.GenDecl:
				if node.Tok == token.CONST {
					for _, spec := range node.Specs {
						if valueSpec, ok := spec.(*ast.ValueSpec); ok {
							for _, name := range valueSpec.Names {
								// Check if this constant matches our type pattern
								if strings.HasPrefix(name.Name, typeName) {
									constants = append(constants, name.Name)
								}
							}
						}
					}
				}
			}
			return true
		})
	}

	if typeDecl == nil {
		return TypeInfo{}, "", "", fmt.Errorf("type %s not found", typeName)
	}

	if len(constants) == 0 {
		return TypeInfo{}, "", "", fmt.Errorf("no constants found for type %s", typeName)
	}

	typeInfo := TypeInfo{
		TypeName:  typeName,
		Constants: constants,
		LowerType: strings.ToLower(typeName),
	}

	return typeInfo, sourceFileName, packageName, nil
}
