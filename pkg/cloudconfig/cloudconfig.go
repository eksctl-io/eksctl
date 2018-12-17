package cloudconfig

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/kris-nova/logger"
)

const (
	header = "#cloud-config"

	// Shell defines the used shell
	Shell = "/bin/bash"

	scriptDir = "/var/lib/cloud/scripts/per-instance/"

	defaultOwner             = "root:root"
	defaultScriptPermissions = "0755"
	defaultFilePermissions   = "0644"
)

// CloudConfig stores informaiton of the cloud config
type CloudConfig struct {
	Commands   []interface{} `json:"runcmd"`
	Packages   []string      `json:"packages"`
	WriteFiles []File        `json:"write_files"`
}

// File stores information about the file
type File struct {
	Content     string `json:"content"`
	Owner       string `json:"owner"`
	Path        string `json:"path"`
	Permissions string `json:"permissions"`
}

// New creates a new cloud config
func New() *CloudConfig {
	return &CloudConfig{}
}

// AddPackages adds packages, which should be installed on the node
func (c *CloudConfig) AddPackages(pkgs ...string) {
	c.Packages = append(c.Packages, pkgs...)
}

// AddCommands adds commands, which will be run on node start up
func (c *CloudConfig) AddCommands(cmds ...[]string) {
	for _, cmd := range cmds {
		c.Commands = append(c.Commands, cmd)
	}
}

// AddCommand adds a command, which will be run on node start up
func (c *CloudConfig) AddCommand(cmd ...string) {
	c.Commands = append(c.Commands, cmd)
}

// AddShellCommand adds a shell comannd, which will be run on node start up
func (c *CloudConfig) AddShellCommand(cmd string) {
	c.Commands = append(c.Commands, []string{Shell, "-c", cmd})
}

// AddFile adds a file, which will be placed on the node
func (c *CloudConfig) AddFile(f File) {
	if f.Owner == "" {
		f.Owner = defaultOwner
	}
	if f.Permissions == "" {
		f.Permissions = defaultFilePermissions
	}
	c.WriteFiles = append(c.WriteFiles, f)
}

// AddScript adds a scipt, which will be placed on the node
func (c *CloudConfig) AddScript(p, s string) {
	c.AddFile(File{
		Content:     s,
		Path:        p,
		Permissions: defaultScriptPermissions,
		Owner:       defaultOwner,
	})
}

// RunScript adds and runs a script on the node
func (c *CloudConfig) RunScript(name, s string) {
	p := scriptDir + name
	c.AddScript(p, s)
	c.AddCommand(p)
}

// Encode encodes the cloud config
func (c *CloudConfig) Encode() (string, error) {
	buf := &bytes.Buffer{}

	data, err := yaml.Marshal(c)
	if err != nil {
		return "", err
	}

	data = append([]byte(fmt.Sprintln(header)), data...)

	gw := gzip.NewWriter(buf)

	if _, err = gw.Write(data); err != nil {
		return "", err
	}
	if err = gw.Close(); err != nil {
		return "", err
	}
	data, err = ioutil.ReadAll(buf)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

// DecodeCloudConfig decodes the cloud config
func DecodeCloudConfig(s string) (*CloudConfig, error) {
	if s == "" {
		return nil, fmt.Errorf("cannot decode empty string")
	}
	c := New()

	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}

	gr, err := gzip.NewReader(ioutil.NopCloser(bytes.NewBuffer(data)))
	if err != nil {
		return nil, err
	}
	defer close(gr)
	data, err = ioutil.ReadAll(gr)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func close(c io.Closer) {
	if err := c.Close(); err != nil {
		logger.Debug("could not close file: %v", err)
	}
}
