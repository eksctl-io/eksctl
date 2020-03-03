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
	"strings"
	"sync"
	"syscall"
)

func main() {
	// Handle interrupt signals and shutdown gracefully:
	ctx, interrupts := handleInterruptSignals()
	defer signal.Stop(interrupts)
	signal.Notify(interrupts, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Run integration tests' suites in parallel, using the ginkgo binary, as
	// using "go test" together with Ginkgo's GinkgoWriter doesn't print to the
	// console as tests progress, and as using just one ginkgo call together
	// with its parallel feature randomises both suites AND tests within a
	// suite, which is problematic for our integration tests as they follow a
	// specific sequence.
	// See also:
	// - https://github.com/onsi/ginkgo/issues/633
	// - https://github.com/onsi/ginkgo/issues/526
	wg := &sync.WaitGroup{}
	for i, dir := range listModules() {
		wg.Add(1)
		go runGinkgo(ctx, wg, i, dir)
	}
	wg.Wait() // Block until all ginkgo processes completed.

	// Obligatory "Cowboy Bebop" reference to end integration tests:
	log.Println("================== See ya Space Cowboy! ===================")
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

func runGinkgo(ctx context.Context, wg *sync.WaitGroup, id int, dir string) {
	defer wg.Done()
	args := []string{"-tags", "integration", "-timeout", "150m", "-v", "--progress", fmt.Sprintf("%s/...", dir), "--"}
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
	scan(outPipe, err, os.Stdout, prefix)
	errPipe, err := cmd.StderrPipe()
	scan(errPipe, err, os.Stderr, prefix)
	if err := cmd.Start(); err != nil {
		log.Fatalf("%s %v", prefix, err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatalf("%s %v", prefix, err)
	}
	fmt.Println(prefix + "=================================== Done ===================================")
}

func scan(in io.ReadCloser, err error, out io.Writer, prefix string) {
	if err != nil {
		log.Fatalf("%s Failed to create stdout/stderr pipe: %v", prefix, err)
	}
	scanner := bufio.NewScanner(in)
	go func() {
		for scanner.Scan() {
			fmt.Fprintln(out, prefix+scanner.Text())
		}
	}()
}
