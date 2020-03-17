package utils

import (
	"bytes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"strconv"
	"strings"
	"testing"
)

func TestValidateLoggingFlags(t *testing.T) {
	var (
		errBothEmpty        = "at least one flag has to be provided"
		errAllBoth          = "cannot use `all` for both"
		errOverlappingTypes = "log types cannot be part of --enable-types and --disable-types"
		errUnknownType      = "unknown log type"
	)

	flagTests := []struct {
		toEnable   []string
		toDisable  []string
		errPattern string
	}{
		{
			toEnable:   []string{"all"},
			toDisable:  []string{"audit", "scheduler", "api"},
			errPattern: "",
		},
		{
			toEnable:   []string{"audit", "scheduler", "api"},
			toDisable:  []string{"controllerManager", "authenticator"},
			errPattern: "",
		},
		{
			toEnable:   []string{"all"},
			toDisable:  []string{"controllerManager", "authenticator"},
			errPattern: "",
		},
		{
			toEnable:   []string{"all"},
			toDisable:  []string{"all"},
			errPattern: errAllBoth,
		},
		{
			toEnable:  []string{"all", "api"},
			toDisable: []string{"all"},
			// TODO improve error reporting for {"all", "api", ...}
			errPattern: errUnknownType,
		},
		{
			toEnable:   []string{""},
			toDisable:  []string{"all"},
			errPattern: errUnknownType,
		},
		{
			toEnable:   []string{"all", "invalid"},
			errPattern: errUnknownType,
		},
		{
			toDisable:  []string{"all", "api", "invalid", "scheduler"},
			errPattern: errUnknownType,
		},
		{
			toEnable:   []string{""},
			errPattern: errUnknownType,
		},
		{
			toEnable:   []string{""},
			toDisable:  []string{""},
			errPattern: errUnknownType,
		},
		{
			errPattern: errBothEmpty,
		},
		{
			toEnable:   []string{"api", "audit"},
			toDisable:  []string{"audit", "authenticator", "scheduler"},
			errPattern: errOverlappingTypes,
		},
		{
			toEnable:   []string{"audit", "authenticator", "scheduler"},
			toDisable:  []string{"api", "scheduler", "controllerManager"},
			errPattern: errOverlappingTypes,
		},
	}

	for i, tt := range flagTests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			err := validateLoggingFlags(tt.toEnable, tt.toDisable)
			if err != nil {
				if tt.errPattern == "" {
					t.Errorf("unexpected error: %v", err)
				} else if !strings.Contains(err.Error(), tt.errPattern) {
					t.Errorf("expected error %q to match %q", err, tt.errPattern)
				}
			} else if tt.errPattern != "" {
				t.Errorf("expected error %q; got nil", tt.errPattern)
			}
		})
	}

}

var _ = Describe("utils", func() {
	Describe("invalid-resource", func() {
		It("with no flag", func() {
			cmd := newMockCmd("invalid-resource")
			out, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("unknown command \"invalid-resource\" for \"utils\""))
			Expect(out).To(ContainSubstring("usage"))
		})
		It("with invalid-resource and some flag", func() {
			cmd := newMockCmd("invalid-resource", "--invalid-flag", "foo")
			out, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("unknown command \"invalid-resource\" for \"utils\""))
			Expect(out).To(ContainSubstring("usage"))
		})
		It("with invalid-resource and additional argument", func() {
			cmd := newMockCmd("invalid-resource", "foo")
			out, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("unknown command \"invalid-resource\" for \"utils\""))
			Expect(out).To(ContainSubstring("usage"))
		})
	})
})

func newMockCmd(args ...string) *mockVerbCmd {
	flagGrouping := cmdutils.NewGrouping()
	cmd := Command(flagGrouping)
	cmd.SetArgs(args)
	return &mockVerbCmd{
		parentCmd: cmd,
	}
}

type mockVerbCmd struct {
	parentCmd *cobra.Command
	cmd       *cmdutils.Cmd
}

func (c mockVerbCmd) execute() (string, error) {
	buf := new(bytes.Buffer)
	c.parentCmd.SetOut(buf)
	err := c.parentCmd.Execute()
	return buf.String(), err
}
