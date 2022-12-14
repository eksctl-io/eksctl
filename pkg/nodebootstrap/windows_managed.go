package nodebootstrap

import (
	"encoding/base64"
	"strings"

	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// ManagedWindows implements a bootstrapper for managed Windows nodegroups.
type ManagedWindows struct {
	NodeGroup *api.ManagedNodeGroup
}

// UserData returns the userdata.
func (w *ManagedWindows) UserData() (string, error) {
	if len(w.NodeGroup.PreBootstrapCommands) == 0 {
		return "", nil
	}
	commands := append([]string{`<powershell>`}, w.NodeGroup.PreBootstrapCommands...)
	commands = append(commands, "</powershell>")
	userData := base64.StdEncoding.EncodeToString([]byte(strings.Join(commands, "\n")))
	logger.Debug("user-data = %s", userData)
	return userData, nil
}
