// Package namedtypes provides a go/analysis analyzer enforcing the gomatic Go
// standard that function parameters use named domain types rather than bare
// primitives. v1 covers non-method function declarations whose parameter type is
// a bare predeclared primitive identifier; methods and composite types are
// deferred.
package namedtypes

import (
	"go/ast"
	"go/types"

	goyze "github.com/gomatic/go-yze"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer reports function parameters declared with a bare primitive type.
var Analyzer = &analysis.Analyzer{
	Name:     "namedtypes",
	Doc:      "reports function parameters that use a bare primitive type instead of a named domain type",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

// Registration declares this analyzer to the yze framework.
var Registration = goyze.Registration{
	Name:       "namedtypes",
	Categories: []goyze.Category{"types"},
	URL:        "https://docs.gomatic.dev/yze/namedtypes",
	Analyzer:   Analyzer,
}

// run reports bare-primitive parameters of non-method function declarations.
func run(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	insp.Preorder([]ast.Node{(*ast.FuncDecl)(nil)}, func(n ast.Node) {
		if fn := n.(*ast.FuncDecl); fn.Recv == nil {
			checkParams(pass, fn.Type.Params)
		}
	})
	return nil, nil
}

// checkParams reports each parameter whose type is a bare primitive.
func checkParams(pass *analysis.Pass, params *ast.FieldList) {
	for _, field := range params.List {
		typ := paramType(field.Type)
		if name, ok := barePrimitiveName(pass, typ); ok {
			pass.Reportf(typ.Pos(), "parameter type %s is a bare primitive; define a named domain type", name)
		}
	}
}

// paramType yields the element type of a variadic parameter (the type after the
// `...`), or the field's type unchanged otherwise, so a variadic bare primitive
// (func f(nums ...int)) is checked and reported at its element type rather than
// being skipped because its field type is an *ast.Ellipsis.
func paramType(expr ast.Expr) ast.Expr {
	if ellipsis, ok := expr.(*ast.Ellipsis); ok {
		return ellipsis.Elt
	}
	return expr
}

// barePrimitiveName returns the identifier name when expr names a predeclared
// primitive type by its predeclared name, and false otherwise. The exemption is
// resolved through the identifier's declaring object rather than the type it
// presents, so it holds regardless of the gotypesalias GODEBUG setting (under
// which an alias to a primitive may otherwise surface as a bare *types.Basic):
//   - a named domain type — a defined type (type Count int) or an alias (type
//     Celsius = float64) — is declared in a package (Pkg() != nil), so exempt;
//   - a predeclared object lives in the universe scope (Pkg() == nil); only one
//     whose underlying type is a primitive (excludes the predeclared error/any
//     interfaces, while still catching the byte/rune primitive aliases) is flagged.
func barePrimitiveName(pass *analysis.Pass, expr ast.Expr) (string, bool) {
	ident, ok := expr.(*ast.Ident)
	if !ok {
		return "", false
	}
	obj := pass.TypesInfo.ObjectOf(ident)
	if obj.Pkg() != nil {
		return "", false
	}
	if _, ok := obj.Type().Underlying().(*types.Basic); !ok {
		return "", false
	}
	return ident.Name, true
}
