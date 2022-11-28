package main

import (
	"github.com/nikolaydubina/go-comment-age/commentage"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(commentage.Analyzer) }
