package accessentry

import (
	"context"
	"crypto/sha1"
	"encoding/base32"
	"errors"
	"fmt"
	"strings"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

// A Creator creates access entries.
type Creator struct {
	ClusterName  string
	StackCreator StackCreator
}

// Create creates the specified access entries.
func (m *Creator) Create(ctx context.Context, accessEntries []api.AccessEntry) error {
	taskTree := m.CreateTasks(ctx, accessEntries)
	if errs := taskTree.DoAllSync(); len(errs) > 0 {
		var allErrs []string
		for _, err := range errs {
			allErrs = append(allErrs, err.Error())
		}
		return errors.New(strings.Join(allErrs, "\n"))
	}
	return nil
}

// CreateTasks creates a TaskTree for creating access entries.
func (m *Creator) CreateTasks(ctx context.Context, accessEntries []api.AccessEntry) *tasks.TaskTree {
	taskTree := &tasks.TaskTree{
		Parallel: true,
	}
	for _, ae := range accessEntries {
		taskTree.Append(&accessEntryTask{
			ctx:          ctx,
			info:         fmt.Sprintf("create access entry for principal ARN %s", ae.PrincipalARN),
			clusterName:  m.ClusterName,
			accessEntry:  ae,
			stackCreator: m.StackCreator,
		})
	}
	return taskTree
}

// MakeStackName creates a stack name for the specified access entry.
func MakeStackName(clusterName string, accessEntry api.AccessEntry) string {
	s := sha1.Sum([]byte(accessEntry.PrincipalARN.String()))
	return fmt.Sprintf("eksctl-%s-accessentry-%s", clusterName, base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(s[:]))
}
