package utils

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/kubicorn/kubicorn/pkg/namer"
	"github.com/pkg/errors"

	"k8s.io/client-go/tools/clientcmd"
)

// ClusterName generates a neme string when a and b are empty strings.
// If either a or b are non-empty, it returns whichever is non-empty.
// If neither a nor b are empty, it returns empty name, to indicate
// ambigous usage.
func ClusterName(a, b string) string {
	if a != "" && b != "" {
		return ""
	}
	if a != "" {
		return a
	}
	if b != "" {
		return b
	}
	return fmt.Sprintf("%s-%d", namer.RandomName(), time.Now().Unix())
}

func ConfigPath(name string) string {
	return path.Join(clientcmd.RecommendedConfigDir, "eksctl", "clusters", name)
}

func isValidConfig(p, name string) error {
	clientConfig, err := clientcmd.LoadFromFile(p)
	if err != nil {
		return errors.Wrapf(err, "unable to load config %q", p)
	}

	if err := clientcmd.ConfirmUsable(*clientConfig, ""); err != nil {
		return errors.Wrapf(err, "unable to parse config %q", p)
	}

	// we want to make sure we only delete config files that haven't be modified by the user
	// checking context name is a good start, we might want ot do deeper checks later, e.g. checksum,
	// as we don't want to delete any files by accident that didn't belong to us
	ctxFmtErr := fmt.Errorf("unable to verify ownership of config %q, unexpected contex name %q", p, clientConfig.CurrentContext)

	ctx := strings.Split(clientConfig.CurrentContext, "@")
	if len(ctx) != 2 {
		return ctxFmtErr
	}
	if strings.HasPrefix(ctx[1], name+".") && strings.HasSuffix(ctx[1], ".eksctl.io") {
		return nil
	}
	return ctxFmtErr
}

func tryDeleteConfig(p, name string) {
	if err := isValidConfig(p, name); err != nil {
		logger.Debug("ignoring error while checking config file %q: %s", p, err.Error())
		return
	}
	if err := os.Remove(p); err != nil {
		logger.Debug("ignoring error while removing config file %q: %s", p, err.Error())
	}
}

func MaybeDeleteConfig(name string) {
	p := ConfigPath(name)

	autoConfExists, err := FileExists(p)
	if err != nil {
		logger.Debug("error checking if auto-generated kubeconfig file exists %q: %s", p, err.Error())
		return
	}
	if autoConfExists {
		if err := os.Remove(p); err != nil {
			logger.Debug("ignoring error while removing auto-generated config file %q: %s", p, err.Error())
		}
		return
	}

	// Print message to manually remove from config file
	logger.Warning("as you are not using the auto-generated kubeconfig file you will need to remove the details of cluster %s manually", name)
}

// FileExists checks to see if a file exists.
func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
