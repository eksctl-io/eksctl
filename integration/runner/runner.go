package runner

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"

	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/types"
)

// Cmd holds definition of a command
type Cmd struct {
	execPath string
	args     []string
	env      []string
	cleanEnv bool
	timeout  time.Duration
}

// NewCmd constructs a new command
func NewCmd(execPath string) Cmd {
	return Cmd{
		execPath: execPath,
		timeout:  20 * time.Second,
	}
}

// NewCmd presents a the command as a string; each argument is quoted,
// so that the command can be copied and called from a shell
func (c Cmd) String() string {
	s := c.execPath
	for _, arg := range c.args {
		s += fmt.Sprintf(" %q", arg)
	}
	return s
}

// WithArgs returns a copy of the command with new arguments
func (c Cmd) WithArgs(args ...string) Cmd {
	c.args = append(c.args, args...)
	return c
}

// WithoutArg removes an existing argument
func (c Cmd) WithoutArg(arg, val string) Cmd {
	var argIdx int
	var found bool
	for i, earg := range c.args {
		if arg == earg {
			argIdx = i
			found = true
		}
	}
	if found {
		endIdx := argIdx + 1
		if val != "" {
			endIdx = argIdx + 2
		}
		c.args = append(c.args[:argIdx], c.args[endIdx:]...)
	}
	return c
}

// WithEnv returns a copy of the command with new environment variables
func (c Cmd) WithEnv(env ...string) Cmd {
	c.env = append(c.env, env...)
	return c
}

// WithCleanEnv returns a copy of the command with all environment variables
// reset (including OS environment variables)
func (c Cmd) WithCleanEnv() Cmd {
	c.env = []string{}
	c.cleanEnv = true
	return c
}

// WithTimeout returns a copy of the command with new timeout value
func (c Cmd) WithTimeout(timeout time.Duration) Cmd {
	c.timeout = timeout
	return c
}

// Start the command and returns underlying session
func (c Cmd) Start() *gexec.Session {
	command := exec.Command(c.execPath, c.args...)

	if c.cleanEnv {
		command.Env = []string{}
	} else {
		command.Env = os.Environ()
	}
	command.Env = append(command.Env, c.env...)

	fmt.Fprintf(GinkgoWriter, "starting '%s'\n", c.String())

	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	if err != nil {
		Fail(fmt.Sprintf("error starting process: %v\n", err), 1)
	}

	return session
}

// Run the command and wait for it to return
func (c Cmd) Run() *gexec.Session {
	session := c.Start()
	session.Wait(c.timeout)
	return session
}

type runCmdMatcher struct {
	session   *gexec.Session
	failedCmd string
}

func (m *runCmdMatcher) run(cmd Cmd) bool {
	m.session = cmd.Run()
	if m.session.ExitCode() == 0 {
		return true
	}
	m.failedCmd = cmd.String()
	return false
}

// RunSuccessfully matches successful execution of a command
func RunSuccessfully() types.GomegaMatcher {
	return &runCmdMatcher{}
}

func (m *runCmdMatcher) Match(actual interface{}) (bool, error) {
	switch act := actual.(type) {
	case Cmd:
		return m.run(act), nil
	case []Cmd:
		for _, cmd := range act {
			if !m.run(cmd) {
				return false, nil
			}
		}
		return true, nil
	default:
		return false, fmt.Errorf("not a Cmd or []Cmd")
	}
}

func (m *runCmdMatcher) FailureMessage(_ interface{}) string {
	return fmt.Sprintf("Expected '%s' to succeed, got return code %d", m.failedCmd, m.session.ExitCode())
}

func (m *runCmdMatcher) NegatedFailureMessage(_ interface{}) string {
	return "Expected command NOT to succeed"
}

type runCmdOutputMatcher struct {
	outputMatchers []types.GomegaMatcher
	failureMessage string
	splitLines     bool
	*runCmdMatcher
}

// RunSuccessfullyWithOutputString matches successful excution of a command and passes the output string
// to another matcher (the string will include stdout and stderr)
func RunSuccessfullyWithOutputString(outputMatchers ...types.GomegaMatcher) types.GomegaMatcher {
	return &runCmdOutputMatcher{
		outputMatchers: outputMatchers,
		splitLines:     false,
		runCmdMatcher:  &runCmdMatcher{},
	}
}

// RunSuccessfullyWithOutputStringLines matches successful excution of a command and passes the output string
// to another matcher split into lines (the string will include stdout and stderr)
func RunSuccessfullyWithOutputStringLines(outputMatchers ...types.GomegaMatcher) types.GomegaMatcher {
	return &runCmdOutputMatcher{
		outputMatchers: outputMatchers,
		splitLines:     true,
		runCmdMatcher:  &runCmdMatcher{},
	}
}

func (m *runCmdOutputMatcher) Match(actual interface{}) (bool, error) {
	cmd, ok := actual.(Cmd)
	if !ok {
		return false, fmt.Errorf("not a Cmd")
	}

	if !m.runCmdMatcher.run(cmd) {
		m.failureMessage = m.runCmdMatcher.FailureMessage(nil)
		return false, nil
	}

	outputString := string(m.runCmdMatcher.session.Buffer().Contents())

	var output interface{}
	if m.splitLines {
		output = strings.Split(outputString, "\n")
	} else {
		output = outputString
	}

	for _, outputMatcher := range m.outputMatchers {
		if ok, err := outputMatcher.Match(output); !ok {
			m.failureMessage = outputMatcher.FailureMessage(output)
			return false, err
		}
	}

	return true, nil
}

func (m *runCmdOutputMatcher) FailureMessage(_ interface{}) string {
	return m.failureMessage
}

func (m *runCmdOutputMatcher) NegatedFailureMessage(_ interface{}) string {
	return "Expected command NOT to succeed"
}
