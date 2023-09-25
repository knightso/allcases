// The allcases command runs the allcases analyzer.
package main

import (
	"github.com/esh2n/allcases"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(allcases.Analyzer) }
