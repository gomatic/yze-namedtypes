package namedtypes

// White-box tests for the fix-eligibility contract on shapes the analysistest
// corpus cannot express: a bodyless function declaration (an assembly stub
// would not type-check in testdata) and a _test.go source path (the production
// driver loads packages without test files, so testdata cannot carry one
// without duplicating diagnostics across package variants).

import (
	"go/ast"
	"testing"

	"github.com/stretchr/testify/assert"
)

// eligibleFixture yields a minimal unexported single-string-parameter function
// declaration that satisfies every fixEligible condition.
func eligibleFixture() (*ast.FuncDecl, *ast.Field) {
	field := &ast.Field{
		Names: []*ast.Ident{ast.NewIdent("s")},
		Type:  ast.NewIdent("string"),
	}
	fn := &ast.FuncDecl{
		Name: ast.NewIdent("lower"),
		Type: &ast.FuncType{Params: &ast.FieldList{List: []*ast.Field{field}}},
		Body: &ast.BlockStmt{},
	}
	return fn, field
}

func TestFixEligible(t *testing.T) {
	tests := []struct {
		mutate   func(fn *ast.FuncDecl, field *ast.Field)
		name     string
		path     sourcePath
		eligible bool
	}{
		{
			name:     "eligible unexported function in a non-test file",
			path:     "a.go",
			mutate:   func(*ast.FuncDecl, *ast.Field) {},
			eligible: true,
		},
		{
			name:     "exported function",
			path:     "a.go",
			mutate:   func(fn *ast.FuncDecl, _ *ast.Field) { fn.Name = ast.NewIdent("Upper") },
			eligible: false,
		},
		{
			name:     "test file",
			path:     "a_test.go",
			mutate:   func(*ast.FuncDecl, *ast.Field) {},
			eligible: false,
		},
		{
			name:     "bodyless declaration",
			path:     "a.go",
			mutate:   func(fn *ast.FuncDecl, _ *ast.Field) { fn.Body = nil },
			eligible: false,
		},
		{
			name: "shared multi-name field",
			path: "a.go",
			mutate: func(_ *ast.FuncDecl, field *ast.Field) {
				field.Names = append(field.Names, ast.NewIdent("t"))
			},
			eligible: false,
		},
		{
			name:     "unnamed parameter",
			path:     "a.go",
			mutate:   func(_ *ast.FuncDecl, field *ast.Field) { field.Names = nil },
			eligible: false,
		},
		{
			name: "variadic parameter",
			path: "a.go",
			mutate: func(_ *ast.FuncDecl, field *ast.Field) {
				field.Type = &ast.Ellipsis{Elt: field.Type}
			},
			eligible: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn, field := eligibleFixture()
			tt.mutate(fn, field)
			assert.Equal(t, tt.eligible, fixEligible(tt.path, fn, field))
		})
	}
}
