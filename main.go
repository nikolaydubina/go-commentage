package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"log"
	"os/exec"
	"strings"
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

	inspect.Preorder([]ast.Node{&ast.FuncDecl{}}, func(n ast.Node) {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn == nil {
			return
		}

		cg := fn.Doc
		if cg == nil {
			return
		}

		if len(cg.List) == 0 {
			return
		}

		filename := pass.Fset.Position(cg.Pos()).Filename
		start := pass.Fset.Position(cg.Pos()).Line
		end := pass.Fset.Position(cg.End()).Line

		log.Println(filename, start, end, cg.List[0].Text)
	})

	return nil, nil
}

type DetailsForFile struct {
	CommitForLine     map[int]string
	LastUpdateForLine map[int]time.Time
}

// ProcessGitProvider makes OS process calls to git to fech data.
// Utilizes human-readable git blame output.
// TODO: Caches data per file.
type ProcessGitProvider struct {
	Files map[string]DetailsForFile
}

func (s ProcessGitProvider) processFile(filename string) error {
	cmd := exec.Command("git", "blame", "-t", "-e", filename)

	var stderr, stdout bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("can not execute command: %w", err)
	}

	if stderr.Len() > 0 {
		return fmt.Errorf("stderr > 0: %s", stderr.String())
	}

	// parse git blame to get line details for whole file
	commitForLine := make(map[int]string)
	lastUpdateForLine := make(map[int]time.Time)

	scanner := bufio.NewScanner(&stdout)
	for scanner.Scan() {
		// last line can be empty line in blame and in go file
		// not counted towards lines in source code
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		lineDetails, err := parseBlameLine(line)
		if err != nil {
			log.Println("error at line: %s", err)
			continue
		}

		commitForLine[lineDetails.lineNumber] = lineDetails.commit
		lastUpdateForLine[lineDetails.lineNumber] = lineDetails.createdAt
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("can not scan: %w", err)
	}

	// check that all lines collected
	if len(commitForLine) == 0 || len(lastUpdateForLine) == 0 {
		return errors.New("missing lines details")
	}
	if len(commitForLine) != len(lastUpdateForLine) {
		return fmt.Errorf("num lines with details mismatch")
	}
	// check line numbers are continious
	maxLine := 1
	for k := range commitForLine {
		if k > maxLine {
			maxLine = k
		}
	}
	for i := 1; i <= maxLine; i++ {
		if _, ok := commitForLine[i]; !ok {
			return fmt.Errorf("missing line(%d)", i)
		}
		if _, ok := lastUpdateForLine[i]; !ok {
			return fmt.Errorf("missing line(%d)", i)
		}
	}

	// store
	if s.Files == nil {
		s.Files = make(map[string]DetailsForFile)
	}
	s.Files[filename] = DetailsForFile{
		CommitForLine:     commitForLine,
		LastUpdateForLine: lastUpdateForLine,
	}

	return nil
}

type blameLine struct {
	lineNumber int
	commit     string
	createdAt  time.Time
}

// Example:
// ef0c9f0c5b8e pkg/util/node/node.go (<djmm@google.com>                 1464913558 -0700   2) Copyright 2015 The Kubernetes Authors.
func parseBlameLine(line string) (blameLine, error) {
	fields := strings.Fields(line)
	return blameLine{}, nil
}

func (s ProcessGitProvider) CommitForLine(filename string, line int) string {
	if _, ok := s.Files[filename]; !ok {
		s.processFile(filename)
	}
	return s.Files[filename].CommitForLine[line]
}

func (s ProcessGitProvider) LastUpdateForLine(filename string, line int) time.Time {
	if _, ok := s.Files[filename]; !ok {
		s.processFile(filename)
	}
	return s.Files[filename].LastUpdateForLine[line]
}

func main() { singlechecker.Main(Analyzer) }
