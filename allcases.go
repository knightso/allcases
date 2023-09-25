// Package allcases defines an Analyzer that checks for switch statements to have all cases.
package allcases

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/types"
	"sort"
	"strings"

	"github.com/gostaticanalysis/comment"
	"github.com/gostaticanalysis/comment/passes/commentmap"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:      "allcases",
	Doc:       Doc,
	Run:       run,
	Requires:  []*analysis.Analyzer{inspect.Analyzer, commentmap.Analyzer},
	FactTypes: []analysis.Fact{new(fact)},
}

// map[<type name>]<constant names>
type fact map[string][]string

func (*fact) AFact() {}
func (f *fact) String() string {
	// sort by key for test
	keys := make([]string, 0, len(*f))
	for key := range *f {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	buf := new(bytes.Buffer)
	for i, key := range keys {
		if i > 0 {
			fmt.Fprint(buf, " ")
		}
		fmt.Fprintf(buf, "%s:%s", key, (*f)[key])
	}

	return buf.String()
}

const Doc = `check for switch statements to have all cases to have all cases

This checker checks for switch statements which have an annotation comment "// allcases"
to have all cases for consts of the expression type.

For example if there is a type and consts below:

	type TestKind int

	const (
		TestKindHoge TestKind = iota
		TestKindFuga
		TestKindPiyo
	)

and switch statements like:

	// allcases
	switch v {
	case TestKindHoge:
		// do something
	case TestKindFuga:
		// do something
	}

then the checker reports that it doesn't have the TestKindPiyo case.
`

const annotation = "allcases"

func run(pass *analysis.Pass) (interface{}, error) {
	importedPkgs := make(map[string]*types.Package)
	importedPkgs[pass.Pkg.Path()] = pass.Pkg

	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	cmaps := pass.ResultOf[commentmap.Analyzer].(comment.Maps)

	// collect imported packages
	inspect.Preorder([]ast.Node{(*ast.ImportSpec)(nil)}, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.ImportSpec:
			obj, ok := pass.TypesInfo.Implicits[n]
			if !ok {
				obj = pass.TypesInfo.Defs[n.Name] // renaming import
			}
			imported := obj.(*types.PkgName).Imported()
			importedPkgs[imported.Path()] = imported
		}
	})

	// collect constants
	inspect.Preorder([]ast.Node{(*ast.Ident)(nil)}, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.Ident:

			// check const
			if n.Obj == nil || n.Obj.Kind != ast.Con {
				return
			}

			o := pass.TypesInfo.ObjectOf(n)
			t := pass.TypesInfo.TypeOf(n)

			var f fact
			if !pass.ImportPackageFact(pass.Pkg, &f) {
				f = make(fact)
			}

			on := objName(o)

			// check for duplicate
			for _, c := range f[t.String()] {
				if c == on {
					return
				}
			}

			f[t.String()] = append(f[t.String()], objName(o))

			pass.ExportPackageFact(&f)
		}
	})

	nodeFilter := []ast.Node{
		(*ast.SwitchStmt)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.SwitchStmt:
			if !cmaps.Annotated(n, annotation) {
				break
			}

			tv, ok := pass.TypesInfo.Types[n.Tag]
			if !ok {
				// TODO: consider this case can happen
				break
			}
			pkg := pass.Pkg
			tn := tv.Type.String()
			ss := strings.SplitN(tn, ".", 2)
			if len(ss) == 2 {
				pkg = importedPkgs[ss[0]]
			}

			var f fact
			if !pass.ImportPackageFact(pkg, &f) {
				pass.Reportf(n.Switch, "unexpected type:%s", tn)
				break
			}

			used := make(map[string]bool)
			for _, stmt := range n.Body.List {
				if cc, ok := stmt.(*ast.CaseClause); ok {
					for _, expr := range cc.List {
						switch expr := expr.(type) {
						case *ast.SelectorExpr:
							o := pass.TypesInfo.ObjectOf(expr.Sel)
							used[objName(o)] = true
						case *ast.Ident:
							o := pass.TypesInfo.ObjectOf(expr)
							used[objName(o)] = true
						}
					}
				}
			}

			// check finally
			var unused []string
			for _, c := range f[tn] {
				if len(unused) >= 3 {
					// omit
					unused = append(unused, "more")
					break
				}
				if !used[c] {
					unused = append(unused, c)
					continue
				}
			}
			if len(unused) >= 2 {
				unused[len(unused)-1] = "and " + unused[len(unused)-1]
			}
			if len(unused) > 0 {
				pass.Reportf(n.Switch, "no case of %s", strings.Join(unused, ", "))
			}
		}
	})

	return nil, nil
}

func objName(o types.Object) string {
	return fmt.Sprintf("%s.%s", o.Pkg().Path(), o.Name())
}
