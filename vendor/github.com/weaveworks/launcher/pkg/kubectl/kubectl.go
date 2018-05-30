package kubectl

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	podSuccess = "Succeeded"
	podFailure = "Failed"
)

type pod struct {
	Status status `json:"status"`
}

type status struct {
	Phase string `json:"phase"`
}

// Client implements a kubectl client to execute commands
type Client interface {
	Execute(args ...string) (string, error)
	ExecuteOutputMatrix(args ...string) (stdout, stderr, combined string, err error)
}

// Execute executes kubectl <args> and returns the combined stdout/err output.
func Execute(c Client, args ...string) (string, error) {
	return c.Execute(args...)
}

// ExecuteJSON execute kubectl <args> and returns the combined json stdout and err output.
func ExecuteJSON(c Client, o interface{}, args ...string) error {
	a := append(args, "-ojson")
	oJSON, err := c.Execute(a...)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(oJSON), o)
	if err != nil {
		return err
	}

	return nil
}

// ClusterInfo describes a Kubernetes cluster
type ClusterInfo struct {
	Name          string
	ServerAddress string
}

// Example stdout:
// Client Version: version.Info{Major:"1", Minor:"9", GitVersion:"v1.9.2", ..., Platform:"linux/amd64"}
// Server Version: version.Info{Major:"1", Minor:"9", GitVersion:"v1.9.3", ..., Platform:"linux/amd64"}

// We don't care about the exact reason why the parsing failed, we'll display
// more context in the error message anyway.
var errParsing = errors.New("parse error")

func parseVersionLine(line string) (string, error) {
	// Only interested in what's between '{', '}'
	idx := strings.Index(line, "{")
	list := line[idx+1 : len(line)-2]

	parts := strings.Split(list, ",")
	for _, part := range parts {
		// parts are of the form key:"value"
		part := strings.TrimSpace(part)
		colon := strings.Index(part, ":")
		if colon == -1 {
			return "", errParsing
		}
		key := part[0:colon]

		if key == "GitVersion" {
			value := part[colon+2 : len(part)-1]
			return value, nil
		}
	}

	return "", nil
}

func parseVersionOutput(stdout string) (clientVersion, serverVersion string, err error) {
	// Protect against invalid input triggering panics (eg. out of bounds).
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("error parsing: %s", stdout)
		}
	}()

	lines := strings.Split(strings.TrimSuffix(stdout, "\n"), "\n")

	// Associate line prefixes to their respective version.
	a := []struct {
		prefix, version string
	}{
		{"Client Version:", ""},
		{"Server Version:", ""},
	}
	for _, line := range lines {
		for i := range a {
			if strings.HasPrefix(line, a[i].prefix) {
				version, err := parseVersionLine(line)
				if err != nil {
					return "", "", fmt.Errorf("error parsing: %s", line)
				}
				a[i].version = version
			}
		}
	}

	return a[0].version, a[1].version, nil
}

// GetVersionInfo returns the version metadata from kubectl
// May return a value for the kubectl client version, despite also returning an error
func GetVersionInfo(c Client) (string, string, error) {
	// Capture stdout only (to ignore server reachability errors)
	stdout, stderr, _, err := c.ExecuteOutputMatrix("version")
	clientVersion, serverVersion, parseErr := parseVersionOutput(stdout)
	// If the server is unreachable, we might have an error but parsable output
	if parseErr != nil {
		if err != nil {
			if stderr == "" {
				return "", "", err
			}
			return "", "", fmt.Errorf("kubectl error (%v): %s", err, stderr)
		}
		return "", "", fmt.Errorf("error parsing kubectl version output: %s", parseErr)
	}

	// kubectl couldn't contact the server.
	if serverVersion == "" {
		err = errors.New(stderr)
	}
	return clientVersion, serverVersion, err
}

// GetClusterInfo gets the current Kubernetes cluster information
func GetClusterInfo(c Client) (ClusterInfo, error) {
	currentContext, err := Execute(c, "config", "current-context")
	if err != nil {
		return ClusterInfo{}, err
	}

	name, err := Execute(c, "config", "view",
		fmt.Sprintf("-o=jsonpath='{.contexts[?(@.name == \"%s\")].context.cluster}'", currentContext),
	)
	if err != nil {
		return ClusterInfo{}, err
	}

	serverAddress, err := Execute(c,
		"config",
		"view",
		fmt.Sprintf("-o=jsonpath='{.clusters[?(@.name == \"%s\")].cluster.server}'", name),
	)
	if err != nil {
		return ClusterInfo{}, err
	}

	return ClusterInfo{
		Name:          name,
		ServerAddress: serverAddress,
	}, nil
}

// Apply applies via kubectl
func Apply(c Client, f string) error {
	// Escape ',' to prevent a url being interpreted as an array, and split
	_, err := Execute(c, "apply", "-f", strings.Replace(f, ",", "%2C", -1))
	return err
}

// ResourceExists return true if the resource exists
func ResourceExists(c Client, resourceType, namespace, resourceName string) (bool, error) {
	_, err := Execute(c, "get", resourceType, resourceName, fmt.Sprintf("--namespace=%s", namespace))
	if err != nil {
		// k8s 1.4 answers with "Error from server: secrets "weave-cloud" not found"
		// More recent versions with "Error from server (NotFound): secrets "weave-cloud" not found
		errorText := err.Error()
		if strings.Contains(errorText, "NotFound") ||
			strings.Contains(errorText, "not found") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// DeleteResource deletes a resource
func DeleteResource(c Client, resourceType, namespace, resourceName string) error {
	_, err := Execute(c, "delete", resourceType, resourceName, fmt.Sprintf("--namespace=%s", namespace))
	return err
}

func isPodReady(c Client, podName, ns string) error {
	// Timeout is set for 30 seconds, as Kubernetes requires some time to create a pod.
	timeout := time.After(1 * time.Minute)
	tick := time.Tick(1 * time.Second)

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timed out during DNS check")
		case <-tick:
			ok, err := checkPod(c, podName, ns)
			if err != nil {
				return err
			} else if ok {
				return nil
			}
		}
	}
}

func checkPod(c Client, podName, ns string) (bool, error) {
	// Retrieve current pod data.
	p := pod{}
	err := ExecuteJSON(c, &p, "get", "pod", podName, "-n", ns)
	if err != nil {
		return false, err
	}

	if p.Status.Phase != podSuccess && p.Status.Phase != podFailure {
		return false, nil
	}

	return true, nil
}

// randomizeName appends random string of numbers to the end of the passed name.
func randomizeName(name string) string {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		// If we fail, we should rather just return the original name.
		return name
	}

	return fmt.Sprintf("%s-%x", name, b)
}

//TestDNS creates a pod where a nslookup is called on a provided domain. It returns true only if the pod was successful.
func TestDNS(c Client, domain string) (bool, error) {
	// Generate a "random" pod name, to prevent "already exists" errors.
	podName := randomizeName("launcher-pre-flight")
	ns := "weave"

	// Create weave namespace, as this happens before any resources are created.
	_, err := CreateNamespace(c, ns)
	if err != nil {
		return false, err
	}

	// Create pod to perform nslookup on a passed domain to check DNS is working.
	_, err = Execute(c, "run", "-n", "weave", "--image", "busybox", "--restart=Never", "--command", podName, "nslookup", domain)
	if err != nil {
		return false, err
	}

	// Initially fetch the pod, which was created above.
	p := pod{}
	err = ExecuteJSON(c, &p, "get", "pod", podName, "-n", ns)
	if err != nil {
		return false, err
	}

	if p.Status.Phase != podSuccess && p.Status.Phase != podFailure {
		// If the state has not been reached yet, we enter a retry phase.
		// In isPodReady function we retry to get pod status phase for a minute and then timeout.
		err := isPodReady(c, podName, ns)
		// Either an error occurred or timeout was reached.
		if err != nil {
			// Attempt to cleanup pod.
			_ = DeleteResource(c, "pod", ns, podName)
			return false, err
		}
	}

	// Get fresh pod data.
	err = ExecuteJSON(c, &p, "get", "pod", podName, "-n", ns)
	if err != nil {
		return false, err
	}

	// If the final status of the pod was failed, we should return an error as DNS is not working.
	if p.Status.Phase == podFailure {
		err = DeleteResource(c, "pod", ns, podName)
		if err != nil {
			return false, err
		}
		return false, fmt.Errorf("checking DNS failed, the DNS in the Kubernetes cluster is not working correctly")
	}

	// This should not happen but still lets error out in case it does.
	if p.Status.Phase != podSuccess {
		return false, fmt.Errorf("checking DNS failed")
	}

	// Cleanup the pod.
	err = DeleteResource(c, "pod", ns, podName)
	if err != nil {
		// We should still return that DNS works, there was only a problem with deleting the resource.
		return true, err
	}

	// We are certain that pod is up and running so we return DNS okay.
	return true, nil
}

// CreateNamespace creates a new namespace and returns whether it was created or not
func CreateNamespace(c Client, namespace string) (bool, error) {
	_, err := Execute(c, "create", "namespace", namespace)
	if err != nil {
		if strings.Contains(err.Error(), "AlreadyExists") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// CreateClusterRoleBinding creates a new cluster role binding
func CreateClusterRoleBinding(c Client, name, role, user string) error {
	_, err := Execute(
		c,
		"create",
		"clusterrolebinding",
		name,
		"--clusterrole",
		role,
		"--user",
		user,
	)
	if err != nil && !strings.Contains(err.Error(), "AlreadyExists") {
		return err
	}
	return nil
}

// CreateSecretFromLiteral creates a new secret with a single (key,value) pair.
func CreateSecretFromLiteral(c Client, namespace, name, key, value string, override bool) (bool, error) {
	secretExists, err := ResourceExists(c, "secret", namespace, name)
	if err != nil {
		return false, err
	}

	if secretExists {
		if !override {
			return false, nil
		}
		err := DeleteResource(c, "secret", namespace, name)
		if err != nil {
			return false, err
		}
	}

	// Create the weave namespace and the weave-cloud secret
	_, err = CreateNamespace(c, namespace)
	if err != nil {
		return false, err
	}

	// Create the secret
	_, err = Execute(c,
		fmt.Sprintf("--namespace=%s", namespace),
		"create",
		"secret",
		"generic",
		name,
		fmt.Sprintf("--from-literal=%s=%s", key, value),
	)
	if err != nil {
		return false, err
	}

	return true, nil
}

type secretManifest struct {
	Data map[string]string
}

// GetSecretValue returns the value of a secret
func GetSecretValue(c Client, namespace, name, key string) (string, error) {
	var secretDefn secretManifest
	err := ExecuteJSON(c, &secretDefn, "get", "secret", name, fmt.Sprintf("--namespace=%s", namespace))
	if err != nil {
		return "", err
	}
	encoded, ok := secretDefn.Data[key]
	if !ok {
		return "", fmt.Errorf("Secret missing key %q", key)
	}
	valueBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(valueBytes), nil
}
