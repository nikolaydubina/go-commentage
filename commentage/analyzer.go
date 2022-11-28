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
	Doc:      "collect details on age(eg, time, commits) of comments and associated code",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

var (
	enableTimeInfo   bool
	enableCommitInfo bool
	verbose          bool
)

func init() {
	Analyzer.Flags.BoolVar(&enableTimeInfo, "time", true, `enable time collection`)
	Analyzer.Flags.BoolVar(&enableCommitInfo, "commit", false, `enable commit collection`)
	Analyzer.Flags.BoolVar(&verbose, "verbose", false, `return diagnostics with more details`)
}

func run(pass *analysis.Pass) (_ interface{}, errf error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	gitStatsProvider := &gitipc.ProcessGitProvider{}
	stats := SimpleStatsComputer{
		GitProider: gitStatsProvider,
	}

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

		// compute stats
		fnstat := FunctionStats{
			Name: fn.Name.Name,
		}

		if enableTimeInfo {
			functionLastUpdatedAt, err := stats.CommmentLastUpdatedAt(filename, functionLineNumStart, functionLineNumEnd)
			if err != nil {
				errf = err
				return
			}

			commentLastUpdatedAt, err := stats.CommmentLastUpdatedAt(filename, docLineNumStart, docLineNumEnd)
			if err != nil {
				errf = err
				return
			}

			fnstat.TimeStats = &TimeStats{
				LastUpdatedAt:    functionLastUpdatedAt,
				DocLastUpdatedAt: commentLastUpdatedAt,
			}
		}

		if enableCommitInfo {
			functionLastCommit, err := stats.LastCommitForRange(filename, functionLineNumStart, functionLineNumEnd)
			if err != nil {
				errf = err
				return
			}

			commentLastCommit, err := stats.LastCommitForRange(filename, docLineNumStart, docLineNumEnd)
			if err != nil {
				errf = err
				return
			}

			commentBehindCommits, err := stats.CommitDifference(filename, docLineNumStart, docLineNumEnd, functionLineNumStart, functionLineNumEnd)
			if err != nil {
				errf = err
				return
			}

			fnstat.CommitStats = &CommitStats{
				LastCommit:       functionLastCommit,
				DocLastCommit:    commentLastCommit,
				DocBehindCommits: commentBehindCommits,
			}
		}

		var diagnosticsMessage string
		if verbose {
			diagnosticsMessage = fnstat.StringVerbose()
		} else {
			diagnosticsMessage = fnstat.String()
		}

		pass.Reportf(fn.Pos(), diagnosticsMessage)
	})

	return nil, errf
}
