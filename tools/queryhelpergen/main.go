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
		filterableIF string
		orderableIF  string
		outFile      string
	)

	flag.StringVar(&filterableIF, "filterable", "", "Name of the Filterable interface (abstract filter definitions)")
	flag.StringVar(&orderableIF, "orderable", "", "Name of the Orderable interface (abstract order definitions)")
	flag.StringVar(&outFile, "out", "", "Output file path for generated code. Defaults to <GOFILE>_queryhelper.go")
	flag.Parse()

	if filterableIF == "" && orderableIF == "" {
		fatalf("provide at least one of -filterable or -orderable")
	}

	pkg, err := loadPackage()
	if err != nil {
		fatalf("load package: %v", err)
	}

	g := newGenerator(pkg)

	// Process Filterable interface -> generate FilterBuilder interface + helpers
	if filterableIF != "" {
		iface := lookupInterface(pkg, filterableIF)
		names := declaredInterfaceMethodNames(pkg, filterableIF)
		methods := methodsByName(iface, names)
		g.addFilterable(filterableIF, methods)
	}

	// Process Orderable interface -> generate OrderByBuilder interface + helpers
	if orderableIF != "" {
		iface := lookupInterface(pkg, orderableIF)
		names := declaredInterfaceMethodNames(pkg, orderableIF)
		methods := methodsByName(iface, names)
		g.addOrderable(orderableIF, methods)
	}

	if filterableIF != "" && orderableIF != "" {
		g.addQueryHelpers()
	}

	// Assemble file
	src := g.render()
	formatted, err := format.Source([]byte(src))
	if err != nil {
		fatalf("format generated code: %v\n---\n%s", err, src)
	}

	if outFile == "" {
		gofile := os.Getenv("GOFILE")
		base := "generated_queryhelper"
		if gofile != "" {
			name := strings.TrimSuffix(gofile, filepath.Ext(gofile))
			if name != "" {
				base = name + "_queryhelper"
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

func lookupInterface(pkg *packages.Package, name string) *types.Interface {
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
	iface.Complete()
	return iface
}

type generator struct {
	pkg       *packages.Package
	imports   map[string]string // importPath -> alias
	usedAlias map[string]bool
	needQuery bool

	// Filterable -> FilterBuilder
	hasFilter          bool
	filterableIFName   string // e.g., "TransactionFilterable"
	filterBuilderName  string // e.g., "TransactionFilterBuilder"
	filterMethods      []filterMethodSpec
	filterHelperPrefix string // e.g., "Transaction"

	// Orderable -> OrderByBuilder
	hasOrder            bool
	orderableIFName     string // e.g., "TransactionOrderable"
	orderByBuilderName  string // e.g., "TransactionOrderByBuilder"
	orderMethods        []orderMethodSpec
	orderByHelperPrefix string // e.g., "Transaction"

	hasQuery bool
}

type filterMethodSpec struct {
	Name      string
	ParamList string // e.g., "amount apd.Decimal"
	ArgList   string // e.g., "amount"
}

type orderMethodSpec struct {
	Name string
}

func newGenerator(pkg *packages.Package) *generator {
	return &generator{
		pkg:       pkg,
		imports:   make(map[string]string),
		usedAlias: make(map[string]bool),
	}
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

// deriveBuilderName transforms "TransactionFilterable" -> "TransactionFilterBuilder"
// and "TransactionOrderable" -> "TransactionOrderByBuilder"
func deriveFilterBuilderName(filterableName string) string {
	if strings.HasSuffix(filterableName, "Filterable") {
		prefix := strings.TrimSuffix(filterableName, "Filterable")
		return prefix + "FilterBuilder"
	}
	return filterableName + "Builder"
}

func deriveOrderByBuilderName(orderableName string) string {
	if strings.HasSuffix(orderableName, "Orderable") {
		prefix := strings.TrimSuffix(orderableName, "Orderable")
		return prefix + "OrderByBuilder"
	}
	return orderableName + "Builder"
}

func deriveHelperPrefix(name string) string {
	for _, suffix := range []string{"Filterable", "Orderable", "FilterBuilder", "OrderByBuilder"} {
		if strings.HasSuffix(name, suffix) {
			return strings.TrimSuffix(name, suffix)
		}
	}
	return name
}

func (g *generator) addFilterable(filterableIFName string, methods []*types.Func) {
	g.needQuery = true
	g.hasFilter = true
	g.filterableIFName = filterableIFName
	g.filterBuilderName = deriveFilterBuilderName(filterableIFName)
	g.filterHelperPrefix = deriveHelperPrefix(filterableIFName)

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
		g.filterMethods = append(g.filterMethods, filterMethodSpec{
			Name:      name,
			ParamList: strings.Join(plist, ", "),
			ArgList:   strings.Join(args, ", "),
		})
	}
}

func (g *generator) addOrderable(orderableIFName string, methods []*types.Func) {
	g.needQuery = true
	g.hasOrder = true
	g.orderableIFName = orderableIFName
	g.orderByBuilderName = deriveOrderByBuilderName(orderableIFName)
	g.orderByHelperPrefix = deriveHelperPrefix(orderableIFName)

	for _, m := range methods {
		g.orderMethods = append(g.orderMethods, orderMethodSpec{
			Name: m.Name(),
		})
	}
}

func (g *generator) addQueryHelpers() {
	g.hasQuery = true
	g.needQuery = true
}

type importSpec struct{ Alias, Path string }

type templateData struct {
	Package              string
	SourceFile           string
	Imports              []importSpec
	HasFilter            bool
	FilterableIFName     string
	FilterBuilderName    string
	FilterMethods        []filterMethodSpec
	FilterHelperPrefix   string
	HasOrder             bool
	OrderableIFName      string
	OrderByBuilderName   string
	OrderMethods         []orderMethodSpec
	OrderByHelperPrefix  string
	HasQuery             bool
}

const fileTemplate = `// Code generated by queryhelpergen; DO NOT EDIT.
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

// {{ .FilterBuilderName }} is the fluent builder interface for constructing filters.
// Implementations are provided by database adapters.
type {{ .FilterBuilderName }} interface {
    query.FilterBuilderLogic[{{ .FilterBuilderName }}]
{{- range .FilterMethods }}
    {{ .Name }}({{ .ParamList }}) {{ $.FilterBuilderName }}
{{- end }}
}

// {{ .FilterHelperPrefix }}Filter provides helper methods for constructing filter predicates.
var {{ .FilterHelperPrefix }}Filter _{{ .FilterHelperPrefix }}Filter

type _{{ .FilterHelperPrefix }}Filter struct {
    query.FilterBase[{{ .FilterBuilderName }}]
}
{{- range .FilterMethods }}

func (_{{ $.FilterHelperPrefix }}Filter) {{ .Name }}({{ .ParamList }}) query.FilterPredicate[{{ $.FilterBuilderName }}] {
    return func(b {{ $.FilterBuilderName }}) {{ $.FilterBuilderName }} {
        return b.{{ .Name }}({{ .ArgList }})
    }
}
{{- end }}
{{- end }}

{{- if .HasOrder }}

// {{ .OrderByBuilderName }} is the fluent builder interface for constructing order clauses.
// Implementations are provided by database adapters.
type {{ .OrderByBuilderName }} interface {
{{- range .OrderMethods }}
    {{ .Name }}(order query.Order) {{ $.OrderByBuilderName }}
{{- end }}
}

// {{ .OrderByHelperPrefix }}OrderBy provides helper methods for constructing order clauses.
var {{ .OrderByHelperPrefix }}OrderBy _{{ .OrderByHelperPrefix }}OrderBy

type _{{ .OrderByHelperPrefix }}OrderBy struct{}
{{- range .OrderMethods }}

func (_{{ $.OrderByHelperPrefix }}OrderBy) {{ .Name }}(order query.Order) query.OrderByFunc[{{ $.OrderByBuilderName }}] {
    return func(b {{ $.OrderByBuilderName }}) {{ $.OrderByBuilderName }} {
        return b.{{ .Name }}(order)
    }
}
{{- end }}
{{- end }}

{{- if .HasQuery }}

func {{ .FilterHelperPrefix }}WithPagination(offset, limit int) query.Option[{{ .FilterBuilderName }}, {{ .OrderByBuilderName }}] {
    return func(b query.Builder[{{ .FilterBuilderName }}, {{ .OrderByBuilderName }}]) {
        b.Paginate(offset, limit)
    }
}
{{- if .HasFilter }}

func {{ .FilterHelperPrefix }}WithFilter(fn query.FilterPredicate[{{ .FilterBuilderName }}]) query.Option[{{ .FilterBuilderName }}, {{ .OrderByBuilderName }}] {
    return func(b query.Builder[{{ .FilterBuilderName }}, {{ .OrderByBuilderName }}]) {
        b.Filter(fn)
    }
}
{{- end }}
{{- if .HasOrder }}

func {{ .FilterHelperPrefix }}WithOrderBy(fns ...query.OrderByFunc[{{ .OrderByBuilderName }}]) query.Option[{{ .FilterBuilderName }}, {{ .OrderByBuilderName }}] {
    return func(b query.Builder[{{ .FilterBuilderName }}, {{ .OrderByBuilderName }}]) {
        b.OrderBy(fns...)
    }
}
{{- end }}
{{- end }}
`

func (g *generator) render() string {
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
		Package:              g.pkg.Name,
		SourceFile:           os.Getenv("GOFILE"),
		Imports:              imports,
		HasFilter:            g.hasFilter,
		FilterableIFName:     g.filterableIFName,
		FilterBuilderName:    g.filterBuilderName,
		FilterMethods:        g.filterMethods,
		FilterHelperPrefix:   g.filterHelperPrefix,
		HasOrder:             g.hasOrder,
		OrderableIFName:      g.orderableIFName,
		OrderByBuilderName:   g.orderByBuilderName,
		OrderMethods:         g.orderMethods,
		OrderByHelperPrefix:  g.orderByHelperPrefix,
		HasQuery:             g.hasQuery,
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

func fatalf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
	os.Exit(2)
}
