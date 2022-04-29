package anywhere_test

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/actions/anywhere"
	"github.com/weaveworks/eksctl/pkg/version"
)

var _ = Describe("Anywhere", func() {
	Context("IsAnywhereCommand", func() {
		It("returns true when anywhere is the argument", func() {
			isAnywhereCommand, err := anywhere.IsAnywhereCommand([]string{"anywhere"})
			Expect(err).NotTo(HaveOccurred())
			Expect(isAnywhereCommand).To(BeTrue())
		})

		When("its not the anywhere command", func() {
			It("returns false", func() {
				isAnywhereCommand, err := anywhere.IsAnywhereCommand([]string{"create", "cluster"})
				Expect(err).NotTo(HaveOccurred())
				Expect(isAnywhereCommand).To(BeFalse())
			})
		})

		When("errors when a flag is present before the anywhere command", func() {
			It("returns false", func() {
				_, err := anywhere.IsAnywhereCommand([]string{"--foo", "anywhere"})
				Expect(err).To(MatchError("flags cannot be placed before the anywhere command"))
			})
		})

		When("when no args are given", func() {
			It("returns false", func() {
				isAnywhereCommand, err := anywhere.IsAnywhereCommand([]string{})
				Expect(err).NotTo(HaveOccurred())
				Expect(isAnywhereCommand).To(BeFalse())
			})
		})
	})

	Context("RunAnywhereCommand", func() {
		var (
			tmpDir       string
			originalPath string
		)
		BeforeEach(func() {
			var err error
			tmpDir, err = os.MkdirTemp("", "anywhere-command")
			Expect(err).NotTo(HaveOccurred())
			err = os.WriteFile(filepath.Join(tmpDir, anywhere.BinaryFileName), []byte(`#!/usr/bin/env sh
echo $@
echo "EKSCTL_VERSION=$EKSCTL_VERSION"
>&2 echo "stderr outputted"
exit 0`), 0777)
			Expect(err).NotTo(HaveOccurred())

			originalPath = os.Getenv("PATH")
			Expect(os.Setenv("PATH", fmt.Sprintf("%s:%s", originalPath, tmpDir))).To(Succeed())
		})

		AfterEach(func() {
			_ = os.RemoveAll(tmpDir)
			Expect(os.Setenv("PATH", originalPath)).To(Succeed())
		})

		It("runs the command", func() {
			newStdoutReader, newStdoutWriter, _ := os.Pipe()
			newStderrReader, newStderrWriter, _ := os.Pipe()
			originalStdout := os.Stdout
			originalStderr := os.Stdout
			defer func() {
				os.Stdout = originalStdout
				os.Stderr = originalStderr
			}()

			os.Stdout = newStdoutWriter
			os.Stderr = newStderrWriter

			exitCode, err := anywhere.RunAnywhereCommand([]string{"anywhere", "--do", "something"})
			Expect(err).NotTo(HaveOccurred())
			Expect(exitCode).To(BeZero())
			newStdoutWriter.Close()
			newStderrWriter.Close()

			By("printing the binary stderr to os.Stderr")
			stderr, _ := io.ReadAll(newStderrReader)
			Expect(string(stderr)).To(Equal("stderr outputted\n"))

			By("printing the binary stdout to os.Stdout and passing the args to the binary")
			stdout, _ := io.ReadAll(newStdoutReader)
			stdoutLines := strings.Split(strings.TrimSuffix(string(stdout), "\n"), "\n")
			Expect(stdoutLines).To(HaveLen(2))
			Expect(stdoutLines[0]).To(Equal("--do something"))

			By("setting the EKSCTL_VERSION environment variable")
			Expect(stdoutLines[1]).To(Equal(fmt.Sprintf("EKSCTL_VERSION=%s", version.GetVersion())))
		})

		When("the binary exits non-zero", func() {
			BeforeEach(func() {
				err := os.WriteFile(filepath.Join(tmpDir, anywhere.BinaryFileName), []byte(`#!/usr/bin/env sh
exit 33`), 0777)
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns the exit code", func() {
				exitCode, err := anywhere.RunAnywhereCommand([]string{"anywhere", "--do", "something"})
				Expect(err).NotTo(HaveOccurred())
				Expect(exitCode).To(Equal(33))
			})
		})
		When("the binary doesn't exist", func() {
			It("returns an understandable error", func() {
				// delete the binary
				_ = os.RemoveAll(tmpDir)
				exitCode, err := anywhere.RunAnywhereCommand([]string{"anywhere", "--do", "something"})
				Expect(err).To(MatchError(fmt.Sprintf("%q plugin was not found on your path", anywhere.BinaryFileName)))
				Expect(exitCode).To(Equal(1))
			})
		})
	})
})
