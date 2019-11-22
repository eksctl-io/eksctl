package cloudconfig

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/kris-nova/logger"

	"sigs.k8s.io/yaml"
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

	// The coreos field is documented on:
	// https://github.com/coreos/coreos-cloudinit/blob/f1f0405491dfd073bbf074f7e374c9ef85600691/Documentation/cloud-config.md
	Coreos Coreos `json:"coreos,omitempty"`
}

// File stores information about the file
type File struct {
	Content     string `json:"content"`
	Owner       string `json:"owner"`
	Path        string `json:"path"`
	Permissions string `json:"permissions"`
}

type Coreos struct {
	Units []SystemdUnit `json:"units,omitempty"`
}

type SystemdUnit struct {
	Name    string `json:"name"`
	Enable  bool   `json:"enable,omitempty"`
	Command string `json:"command,omitempty"`
	Content string `json:"content,omitempty"`
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

// AddSystemdUnit adds a systemd unit
func (c *CloudConfig) AddSystemdUnit(name string, enable bool, command, content string) {
	c.Coreos.Units = append(c.Coreos.Units, SystemdUnit{name, enable, command, content})
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
	data, err := yaml.Marshal(c)
	if err != nil {
		return "", err
	}

	data = append([]byte(fmt.Sprintln(header)), data...)

	return encodeUserData(data)
}

func encodeUserData(data []byte) (string, error) {
	var (
		buf bytes.Buffer
		gw  = gzip.NewWriter(&buf)
	)

	if _, err := gw.Write(data); err != nil {
		return "", err
	}
	if err := gw.Close(); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
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
	defer safeClose(gr)
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

func safeClose(c io.Closer) {
	if err := c.Close(); err != nil {
		logger.Debug("could not close file: %v", err)
	}
}
