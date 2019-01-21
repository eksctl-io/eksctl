package testutils

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
)

// RegisterAndRun setup and run Ginkgo tests
func RegisterAndRun(t *testing.T) {
	_, suitePath, _, _ := runtime.Caller(1)
	RegisterFailHandler(Fail)
	reportPath := os.Getenv("JUNIT_REPORT_DIR")
	if reportPath != "" {
		name := regexp.MustCompile("[^a-zA-Z0-9]+").ReplaceAllString(suitePath, "__")
		reportPath := fmt.Sprintf("%s/%s_%d.xml", reportPath, name, config.GinkgoConfig.ParallelNode)
		fmt.Printf("test result output: %s\n", reportPath)
		junitReporter := reporters.NewJUnitReporter(reportPath)
		RunSpecsWithDefaultAndCustomReporters(t, suitePath, []Reporter{junitReporter})
	} else {
		RunSpecs(t, suitePath)
	}
}
