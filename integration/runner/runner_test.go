// +build !windows

package runner

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("command runner", func() {
	It("can declare commands", func() {
		c1 := NewCmd("foo")

		c2 := c1.WithArgs("bar1", "bar2")

		c3 := c2.WithArgs("bar3").WithEnv("bar=3")

		c4 := c3.WithCleanEnv()

		Expect(c1.execPath).To(Equal("foo"))
		Expect(c1.args).To(BeEmpty())
		Expect(c1.String()).To(Equal("foo"))

		Expect(c2.execPath).To(Equal("foo"))
		Expect(c2.args).To(ConsistOf("bar1", "bar2"))

		Expect(c3.execPath).To(Equal("foo"))
		Expect(c3.args).To(ConsistOf("bar1", "bar2", "bar3"))

		Expect(c3.String()).To(Equal(`foo "bar1" "bar2" "bar3"`))

		Expect(c3.cleanEnv).To(BeFalse())
		Expect(c3.env).To(ConsistOf("bar=3"))

		Expect(c4.cleanEnv).To(BeTrue())
		Expect(c4.env).To(BeEmpty())
	})

	It("can run commands with env vars and match output", func() {
		c4 := NewCmd("env").WithEnv("FOO=bar")

		Expect(c4).To(RunSuccessfullyWithOutputString(ContainSubstring("FOO=bar\n")))
		Expect(c4).To(RunSuccessfullyWithOutputString(ContainSubstring("PATH=")))

		Expect(c4).To(RunSuccessfullyWithOutputStringLines(
			ContainElement(MatchRegexp(".*")),
			ContainElement("FOO=bar"),
			Not(ContainElement("FOO=baz")),
			ContainElement(MatchRegexp("PATH=.*")),
		))

		c5 := c4.WithCleanEnv().WithEnv("BAR=foo")

		Expect(c5).To(RunSuccessfullyWithOutputString(Equal("BAR=foo\n")))
		Expect(c5).To(RunSuccessfullyWithOutputString(Not(ContainSubstring("PATH="))))

		Expect(c5).To(RunSuccessfullyWithOutputStringLines(
			ConsistOf("BAR=foo", ""),
		))

		c6 := NewCmd("echo").WithArgs("{}")

		Expect(c6).To(RunSuccessfullyWithOutputString(MatchJSON("{}")))
	})

	It("can run multiple commands", func() {
		willSucceed := []Cmd{
			NewCmd("echo"),
			NewCmd("true"),
		}

		Expect(willSucceed).To(RunSuccessfully())

		willFail := []Cmd{
			NewCmd("true"),
			NewCmd("false"),
		}
		Expect(willFail).ToNot(RunSuccessfully())
	})

	It("can start a command and interrupt it", func() {
		session := NewCmd("sleep").WithArgs("20").Start()
		Expect(session.Command.Process).ToNot(BeNil())
		session.Interrupt().Wait()
		Expect(session.ExitCode()).ToNot(BeZero())
	})
})
