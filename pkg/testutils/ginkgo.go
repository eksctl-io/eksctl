package testutils

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
)

// RegisterAndRun setup and run Ginkgo tests
func RegisterAndRun(t *testing.T, testDecsription string) {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	filename := reg.ReplaceAllString(testDecsription, "")

	RegisterFailHandler(Fail)
	reportPath := os.Getenv("JUNIT_REPORT_FOLDER")
	if reportPath != "" {
		reportPath := fmt.Sprintf("%s/%s_%d.xml", reportPath, filename, config.GinkgoConfig.ParallelNode)
		fmt.Printf("test result output: %s\n", reportPath)
		junitReporter := reporters.NewJUnitReporter(reportPath)
		RunSpecsWithDefaultAndCustomReporters(t, testDecsription, []Reporter{junitReporter})
	} else {
		RunSpecs(t, testDecsription)
	}

}
