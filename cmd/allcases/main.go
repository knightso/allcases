// The allcases command runs the allcases analyzer.
package main

import (
	"allcases"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(allcases.Analyzer) }
