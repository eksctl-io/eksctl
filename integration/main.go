// +build integration

package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"regexp"
	"strings"
	"sync"
	"syscall"
)

func main() {
	// Handle interrupt signals and shutdown gracefully:
	ctx, interrupts := handleInterruptSignals()
	defer signal.Stop(interrupts)
	signal.Notify(interrupts, os.Interrupt, syscall.SIGTERM)

	// Run integration tests' suites in parallel, using the ginkgo binary, as
	// using "go test" together with Ginkgo's GinkgoWriter doesn't print to the
	// console as tests progress, and as using just one ginkgo call together
	// with its parallel feature randomises both suites AND tests within a
	// suite, which is problematic for our integration tests as they follow a
	// specific sequence.
	// See also:
	// - https://github.com/onsi/ginkgo/issues/633
	// - https://github.com/onsi/ginkgo/issues/526
	testSuites := listModules()
	wg := &sync.WaitGroup{}
	summaries := make(chan []string, 2*len(testSuites))
	statuses := make(chan int, len(testSuites))
	for i, dir := range testSuites {
		wg.Add(1)
		go runGinkgo(ctx, wg, summaries, statuses, i, dir)
	}
	wg.Wait() // Block until all ginkgo processes completed.
	log.Println("================== Summary ================================")
	close(summaries)
	for summary := range summaries {
		for _, line := range summary {
			fmt.Println(line)
		}
	}
	// Obligatory "Cowboy Bebop" reference to end integration tests:
	log.Println("================== See ya Space Cowboy! ===================")
	exit(statuses)
}

func exit(statuses chan int) {
	close(statuses)
	globalStatus := OK
	for status := range statuses {
		globalStatus += status
	}
	os.Exit(globalStatus)
}

func handleInterruptSignals() (context.Context, chan os.Signal) {
	interrupts := make(chan os.Signal, 1)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-interrupts // Blocks here until interrupted
		log.Println("=========================== Shutdown signal received ===========================")
		// Signal cancellation to context.Context, which aborts ginkgo child
		// processes, letting them run AfterSuite to destroy all clusters, and
		// then letting them naturally terminate.
		cancel()
	}()
	return ctx, interrupts
}

func listModules() []string {
	testsDir := path.Join("integration", "tests")
	files, err := ioutil.ReadDir(testsDir)
	if err != nil {
		log.Fatalf("failed to gather test suites: %v", err)
	}
	dirs := []string{}
	for _, f := range files {
		if f.IsDir() {
			dirs = append(dirs, path.Join(testsDir, f.Name()))
		}
	}
	return dirs
}

const (
	OK    = 0
	Error = 1
)

func runGinkgo(ctx context.Context, wg *sync.WaitGroup, summaries chan []string, statuses chan int, id int, dir string) {
	defer wg.Done()
	args := []string{"-tags", "integration", "-v", "--progress", fmt.Sprintf("%s/...", dir), "--"}
	args = append(args, os.Args[1:]...)
	name := "ginkgo"
	prefix := fmt.Sprintf("[%d]", id)
	fmt.Printf("%s Running: %s %s\n", prefix, name, strings.Join(args, " "))

	cmd := exec.CommandContext(ctx, name, args...)
	// We could run the command like this:
	//   cmd.Stdout = os.Stdout
	//   cmd.Stderr = os.Stderr
	//   if err := cmd.Run(); err != nil { ... }
	// and still have the output from the various ginkgo instances to be
	// gathered by the current process, but in order to more easily group logs
	// we use the following scanner which prefixes logs with the goroutine ID.
	outPipe, err := cmd.StdoutPipe()
	scan(outPipe, err, os.Stdout, prefix, summaries)
	errPipe, err := cmd.StderrPipe()
	scan(errPipe, err, os.Stderr, prefix, summaries)
	if err := cmd.Start(); err != nil {
		log.Fatalf("%s %v", prefix, err)
	}
	if err := cmd.Wait(); err != nil {
		log.Printf("%s %v", prefix, err)
		statuses <- Error
	} else {
		statuses <- OK
	}
	fmt.Println(prefix + "=================================== Done ===================================")
}

func scan(in io.ReadCloser, err error, out io.Writer, prefix string, summaries chan []string) {
	if err != nil {
		log.Fatalf("%s Failed to create stdout/stderr pipe: %v", prefix, err)
	}
	scanner := bufio.NewScanner(in)
	summary := []string{}
	captureSummary := false
	go func() {
		for scanner.Scan() {
			rawText := scanner.Text()
			text := fmt.Sprintf("%s %s", prefix, rawText)
			if captureSummary || enableCaptureSummary(rawText) {
				captureSummary = true
				summary = append(summary, text)
			}
			fmt.Fprintln(out, text)
		}
		summaries <- summary
	}()
}

var summaryRE = regexp.MustCompile(`Summarizing \d+ Failure:`)
var ranRE = regexp.MustCompile(`Ran \d+ of \d+ Specs in [\d\.]+ seconds`)

func enableCaptureSummary(out string) bool {
	return summaryRE.MatchString(out) || ranRE.MatchString(out)
}
