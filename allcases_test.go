package allcases_test

import (
	"testing"

	"github.com/knightso/allcases"

	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, allcases.Analyzer, "c")
}
