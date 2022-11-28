package commentage

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/nikolaydubina/go-commentage/gitipc"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "commentage",
	Doc:      "collect details on age of comments and associated code in terms of time and commits",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	gitStatsProvider := &gitipc.ProcessGitProvider{}
	stats := SimpleStatsComputer{
		GitProider: gitStatsProvider,
	}

	var errf error

	inspect.Preorder([]ast.Node{&ast.FuncDecl{}}, func(n ast.Node) {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn == nil {
			return
		}

		filename := pass.Fset.Position(fn.Pos()).Filename

		if !strings.HasSuffix(filename, ".go") {
			return
		}

		functionLineNumStart := pass.Fset.Position(fn.Pos()).Line
		functionLineNumEnd := pass.Fset.Position(fn.End()).Line

		functionLastUpdatedAt, err := stats.CommmentLastUpdatedAt(filename, functionLineNumStart, functionLineNumEnd)
		if err != nil {
			errf = err
			return
		}

		fndoc := fn.Doc
		if fndoc == nil {
			return
		}

		if len(fndoc.List) == 0 {
			return
		}

		docLineNumStart := pass.Fset.Position(fndoc.Pos()).Line
		docLineNumEnd := pass.Fset.Position(fndoc.End()).Line

		if docLineNumEnd >= functionLineNumStart {
			err := fmt.Errorf("function '%s' (%d:%d) overlaps doc(%d:%d)", fn.Name, functionLineNumStart, functionLineNumEnd, docLineNumStart, docLineNumEnd)
			errf = err
			return
		}

		commentLastUpdatedAt, err := stats.CommmentLastUpdatedAt(filename, docLineNumStart, docLineNumEnd)
		if err != nil {
			errf = err
			return
		}

		fnstat := FunctionStats{
			Name:             fn.Name.Name,
			LastUpdatedAt:    functionLastUpdatedAt,
			DocLastUpdatedAt: commentLastUpdatedAt,
		}

		pass.Reportf(fn.Pos(), fnstat.String())
	})

	return nil, errf
}
