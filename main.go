package main

import (
	"fmt"
	"go/ast"
	"log"
	"time"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/singlechecker"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "commentage",
	Doc:      "age (time/commits) of comments as compared with age of associated functions for git repositories",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	gitStatsProvider := &ProcessGitProvider{}
	stats := SimpleStats{
		GitProider: gitStatsProvider,
	}

	inspect.Preorder([]ast.Node{&ast.FuncDecl{}}, func(n ast.Node) {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn == nil {
			return
		}

		filename := pass.Fset.Position(fn.Pos()).Filename

		functionLineNumStart := pass.Fset.Position(fn.Pos()).Line
		functionLineNumEnd := pass.Fset.Position(fn.End()).Line

		functionLastUpdatedAt, err := stats.CommmentLastUpdatedAt(filename, functionLineNumStart, functionLineNumEnd)
		if err != nil {
			log.Println(err)
		}

		cg := fn.Doc
		if cg == nil {
			return
		}

		if len(cg.List) == 0 {
			return
		}

		commentLineNumStart := pass.Fset.Position(cg.Pos()).Line
		commentLineNumEnd := pass.Fset.Position(cg.End()).Line

		if functionLineNumStart >= commentLineNumEnd {
			log.Println("function overlaps comment line numbers")
		}

		commentLastUpdatedAt, err := stats.CommmentLastUpdatedAt(filename, commentLineNumStart, commentLineNumEnd)
		if err != nil {
			log.Println(err)
		}

		log.Println(filename, "function", functionLineNumStart, functionLineNumEnd, functionLastUpdatedAt, "comment", commentLineNumStart, commentLineNumEnd, commentLastUpdatedAt, cg.List[0].Text)
	})

	return nil, nil
}

type GitProider interface {
	CommitForLine(filename string, line int) (string, error)
	LastUpdateForLine(filename string, line int) (time.Time, error)
}

type SimpleStats struct {
	GitProider GitProider
}

func (s SimpleStats) lineRangeLastUpdatedAt(filename string, lineStart int, lineEnd int) (time.Time, error) {
	var maxTime time.Time
	for i := lineStart; i <= lineEnd; i++ {
		lineUpdatedTime, err := s.GitProider.LastUpdateForLine(filename, i)
		if err != nil {
			return time.Time{}, fmt.Errorf("can not get last updated at: %w", err)
		}
		if lineUpdatedTime.After(maxTime) {
			maxTime = lineUpdatedTime
		}
	}
	return maxTime, nil
}

func (s SimpleStats) CommmentLastUpdatedAt(filename string, lineStart int, lineEnd int) (time.Time, error) {
	return s.lineRangeLastUpdatedAt(filename, lineStart, lineEnd)
}

func (s SimpleStats) FunctionLastUpdatedAt(filename string, lineStart int, lineEnd int) (time.Time, error) {
	return s.lineRangeLastUpdatedAt(filename, lineStart, lineEnd)
}

type DetailsForFile struct {
	CommitForLine     map[int]string
	LastUpdateForLine map[int]time.Time
}

func main() { singlechecker.Main(Analyzer) }
