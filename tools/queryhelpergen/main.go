package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"golang.org/x/tools/go/packages"
)

func main() {
	var (
		filterIF string
		orderIF  string
		outFile  string
	)

	flag.StringVar(&filterIF, "filter", "", "Name of the FilterBuilder interface (e.g., FilterBuilder)")
	flag.StringVar(&orderIF, "order", "", "Name of the OrderByBuilder interface (e.g., OrderByBuilder)")
	flag.StringVar(&outFile, "out", "", "Output file path for generated helpers. Defaults to <GOFILE>_query_helpers.go")
	flag.Parse()

	if filterIF == "" && orderIF == "" {
		fatalf("provide at least one of -filter or -order")
	}

	pkg, err := loadPackage()
	if err != nil {
		fatalf("load package: %v", err)
	}

	// Discover interfaces by name
	var (
		filterIface *types.Interface
		orderIface  *types.Interface
	)

	lookupInterface := func(name string) *types.Interface {
		if name == "" {
			return nil
		}
		obj := pkg.Types.Scope().Lookup(name)
		if obj == nil {
			fatalf("interface %q not found in package %s", name, pkg.PkgPath)
		}
		named, ok := obj.Type().(*types.Named)
		if !ok {
			fatalf("%q is not a named type", name)
		}
		iface, ok := named.Underlying().(*types.Interface)
		if !ok {
			fatalf("%q is not an interface type", name)
		}
		// Ensure methods are computed
		iface.Complete()
		return iface
	}

	filterIface = lookupInterface(filterIF)
	orderIface = lookupInterface(orderIF)

	g := newGenerator(pkg)

	// Collect methods and imports for Filter and OrderBy wrappers
	if filterIface != nil {
		names := declaredInterfaceMethodNames(pkg, filterIF)
		methods := methodsByName(filterIface, names)
		methods = filterOutLogical(methods)
		g.addFilterWrappers(filterIF, methods)
	}
	if orderIface != nil {
		names := declaredInterfaceMethodNames(pkg, orderIF)
		methods := methodsByName(orderIface, names)
		g.addOrderByWrappers(orderIF, methods)
	}
	if filterIF != "" && orderIF != "" {
		g.addQueryHelpers(filterIF, orderIF)
	}

	// Assemble file
	src := g.render()
	formatted, err := format.Source([]byte(src))
	if err != nil {
		fatalf("format generated code: %v\n---\n%s", err, src)
	}

	if outFile == "" {
		// If invoked via go:generate, GOFILE is the source file containing the directive.
		// Derive output as <GOFILE base>_query_helpers.go in current directory
		gofile := os.Getenv("GOFILE")
		base := "generated_query_helpers"
		if gofile != "" {
			name := strings.TrimSuffix(gofile, filepath.Ext(gofile))
			if name != "" {
				base = name + "_query_helpers"
			}
		}
		outFile = base + ".go"
	}
	if err := os.MkdirAll(filepath.Dir(outFile), 0o755); err != nil {
		fatalf("ensure out dir: %v", err)
	}
	if err := os.WriteFile(outFile, formatted, 0o644); err != nil {
		fatalf("write output: %v", err)
	}
}

func loadPackage() (*packages.Package, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedModule,
	}
	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		return nil, err
	}
	if packages.PrintErrors(pkgs) > 0 {
		return nil, errors.New("package load error")
	}
	if len(pkgs) != 1 {
		return nil, fmt.Errorf("expected 1 package in current directory, got %d", len(pkgs))
	}
	return pkgs[0], nil
}

// generator holds state for code generation, including import resolution.

type generator struct {
	pkg       *packages.Package
	imports   map[string]string // importPath -> alias
	usedAlias map[string]bool
	needQuery bool

	// template data
	hasFilter     bool
	filterIFName  string
	filterMethods []methodSpec

	hasOrder     bool
	orderIFName  string
	orderMethods []methodSpec

	hasQuery bool
}

func newGenerator(pkg *packages.Package) *generator {
	g := &generator{
		pkg:       pkg,
		imports:   make(map[string]string),
		usedAlias: make(map[string]bool),
	}
	return g
}

const queryImport = "github.com/theater-improrama/go-utils/query"

func (g *generator) ensureImport(path string, suggested string) string {
	if path == "" {
		return ""
	}
	if alias, ok := g.imports[path]; ok {
		return alias
	}
	alias := suggested
	if alias == "" {
		parts := strings.Split(path, "/")
		alias = parts[len(parts)-1]
	}
	if alias == g.pkg.Name || alias == "query" {
		alias = alias + "pkg"
	}
	// prevent collisions
	base := alias
	for i := 1; g.usedAlias[alias]; i++ {
		alias = fmt.Sprintf("%s%d", base, i)
	}
	g.imports[path] = alias
	g.usedAlias[alias] = true
	return alias
}

func (g *generator) typeString(t types.Type) string {
	qual := func(p *types.Package) string {
		if p == nil {
			return ""
		}
		if p.Path() == g.pkg.PkgPath {
			return ""
		}
		if p.Path() == queryImport {
			g.needQuery = true
			return "query"
		}
		alias := g.ensureImport(p.Path(), p.Name())
		return alias
	}
	return types.TypeString(t, qual)
}

func (g *generator) objNameString(v *types.Var, idx int) string {
	name := v.Name()
	if name == "" || name == "_" {
		return fmt.Sprintf("p%d", idx)
	}
	return name
}

func (g *generator) addFilterWrappers(filterIFName string, methods []*types.Func) {
	g.needQuery = true
	g.hasFilter = true
	g.filterIFName = filterIFName
	for _, m := range methods {
		name := m.Name()
		sig := m.Type().(*types.Signature)
		params := sig.Params()
		isVariadic := sig.Variadic()

		var plist []string
		var args []string
		for i := 0; i < params.Len(); i++ {
			p := params.At(i)
			pname := g.objNameString(p, i)
			pt := p.Type()
			var pdecl string
			if isVariadic && i == params.Len()-1 {
				if slice, ok := pt.(*types.Slice); ok {
					pdecl = fmt.Sprintf("%s ...%s", pname, g.typeString(slice.Elem()))
					args = append(args, pname+"...")
				} else {
					pdecl = fmt.Sprintf("%s %s", pname, g.typeString(pt))
					args = append(args, pname)
				}
			} else {
				pdecl = fmt.Sprintf("%s %s", pname, g.typeString(pt))
				args = append(args, pname)
			}
			plist = append(plist, pdecl)
		}
		g.filterMethods = append(g.filterMethods, methodSpec{
			Name:      name,
			ParamList: strings.Join(plist, ", "),
			ArgList:   strings.Join(args, ", "),
		})
	}
}

func (g *generator) addOrderByWrappers(orderIFName string, methods []*types.Func) {
	g.needQuery = true
	g.hasOrder = true
	g.orderIFName = orderIFName
	for _, m := range methods {
		name := m.Name()
		sig := m.Type().(*types.Signature)
		params := sig.Params()
		isVariadic := sig.Variadic()

		var plist []string
		var args []string
		for i := 0; i < params.Len(); i++ {
			p := params.At(i)
			pname := g.objNameString(p, i)
			pt := p.Type()
			var pdecl string
			if isVariadic && i == params.Len()-1 {
				if slice, ok := pt.(*types.Slice); ok {
					pdecl = fmt.Sprintf("%s ...%s", pname, g.typeString(slice.Elem()))
					args = append(args, pname+"...")
				} else {
					pdecl = fmt.Sprintf("%s %s", pname, g.typeString(pt))
					args = append(args, pname)
				}
			} else {
				pdecl = fmt.Sprintf("%s %s", pname, g.typeString(pt))
				args = append(args, pname)
			}
			plist = append(plist, pdecl)
		}
		g.orderMethods = append(g.orderMethods, methodSpec{
			Name:      name,
			ParamList: strings.Join(plist, ", "),
			ArgList:   strings.Join(args, ", "),
		})
	}
}

type methodSpec struct {
	Name      string
	ParamList string
	ArgList   string
}

// addQueryHelpers marks that we should render ListOption helpers using the generic
// query.QueryBuilder[FilterIF, OrderIF] type.
func (g *generator) addQueryHelpers(filterIFName, orderIFName string) {
	g.hasQuery = true
	g.needQuery = true
	g.hasFilter = true
	g.filterIFName = filterIFName
	g.hasOrder = true
	g.orderIFName = orderIFName
}

type importSpec struct{ Alias, Path string }

type templateData struct {
	Package      string
	SourceFile   string
	Imports      []importSpec
	HasFilter    bool
	FilterIFName string
	Filter       []methodSpec
	HasOrder     bool
	OrderIFName  string
	Order        []methodSpec
	HasQuery     bool
}

const fileTemplate = `// Code generated by query-builder-generator; DO NOT EDIT.
// Source: {{.SourceFile}}

package {{.Package}}

{{- if .Imports }}
import (
{{- range .Imports }}
    {{ .Alias }} "{{ .Path }}"
{{- end }}
)
{{- end }}

{{- if .HasFilter }}
var Filter _Filter

type _Filter struct {
    query.FilterBase[{{ .FilterIFName }}]
}
{{- range .Filter }}
func (_Filter) {{ .Name }}({{ .ParamList }}) query.FilterPredicate[{{ $.FilterIFName }}] {
    return func(b {{ $.FilterIFName }}) {{ $.FilterIFName }} {
        return b.{{ .Name }}({{ .ArgList }})
    }
}
{{- end }}
{{- end }}

{{- if .HasOrder }}
var OrderBy _OrderBy

type _OrderBy struct{}
{{- range .Order }}
func (_OrderBy) {{ .Name }}({{ .ParamList }}) query.OrderByFunc[{{ $.OrderIFName }}] {
    return func(b {{ $.OrderIFName }}) {{ $.OrderIFName }} {
        return b.{{ .Name }}({{ .ArgList }})
    }
}
{{- end }}
{{- end }}

{{- if .HasQuery }}
func WithPagination(offset, limit int) ListOption {
    return func(b query.QueryBuilder[{{ .FilterIFName }}, {{ .OrderIFName }}]) {
        b.Paginate(offset, limit)
    }
}
{{- if .HasFilter }}
func WithFilter(fn query.FilterPredicate[{{ .FilterIFName }}]) ListOption {
    return func(b query.QueryBuilder[{{ .FilterIFName }}, {{ .OrderIFName }}]) {
        b.Filter(fn)
    }
}
{{- end }}
{{- if .HasOrder }}
func WithOrderBy(fns ...query.OrderByFunc[{{ .OrderIFName }}]) ListOption {
    return func(b query.QueryBuilder[{{ .FilterIFName }}, {{ .OrderIFName }}]) {
        b.OrderBy(fns...)
    }
}
{{- end }}
{{- end }}
`

func (g *generator) render() string {
	// Build imports list
	var imports []importSpec
	if g.needQuery {
		imports = append(imports, importSpec{Alias: "query", Path: queryImport})
		g.usedAlias["query"] = true
	}
	for path, alias := range g.imports {
		if path == queryImport {
			continue
		}
		imports = append(imports, importSpec{Alias: alias, Path: path})
	}
	sort.Slice(imports, func(i, j int) bool { return imports[i].Alias < imports[j].Alias })

	data := templateData{
		Package:      g.pkg.Name,
		SourceFile:   os.Getenv("GOFILE"),
		Imports:      imports,
		HasFilter:    g.hasFilter,
		FilterIFName: g.filterIFName,
		Filter:       g.filterMethods,
		HasOrder:     g.hasOrder,
		OrderIFName:  g.orderIFName,
		Order:        g.orderMethods,
		HasQuery:     g.hasQuery,
	}

	tpl := template.Must(template.New("file").Parse(fileTemplate))
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		fatalf("execute template: %v", err)
	}
	return buf.String()
}

// declaredInterfaceMethodNames returns only the method names declared directly
// on the named interface (excluding embedded interface methods) by inspecting AST.
func declaredInterfaceMethodNames(pkg *packages.Package, ifaceName string) map[string]bool {
	names := map[string]bool{}
	for _, f := range pkg.Syntax {
		for _, d := range f.Decls {
			gd, ok := d.(*ast.GenDecl)
			if !ok || gd.Tok != token.TYPE {
				continue
			}
			for _, s := range gd.Specs {
				ts, ok := s.(*ast.TypeSpec)
				if !ok || ts.Name == nil || ts.Name.Name != ifaceName {
					continue
				}
				it, ok := ts.Type.(*ast.InterfaceType)
				if !ok || it.Methods == nil {
					continue
				}
				for _, field := range it.Methods.List {
					// Only named fields are explicit method declarations.
					if len(field.Names) == 0 {
						continue // embedded interface, skip
					}
					names[field.Names[0].Name] = true
				}
			}
		}
	}
	return names
}

// methodsByName returns the funcs from iface whose names are in the provided set.
func methodsByName(iface *types.Interface, names map[string]bool) []*types.Func {
	if len(names) == 0 {
		return nil
	}
	var out []*types.Func
	for i := 0; i < iface.NumMethods(); i++ {
		m := iface.Method(i)
		if names[m.Name()] {
			out = append(out, m)
		}
	}
	return out
}

// filterOutLogical removes Not/And/Or methods if present.
func filterOutLogical(methods []*types.Func) []*types.Func {
	var out []*types.Func
	for _, m := range methods {
		switch m.Name() {
		case "Not", "And", "Or":
			continue
		}
		out = append(out, m)
	}
	return out
}

func fatalf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
	os.Exit(2)
}
