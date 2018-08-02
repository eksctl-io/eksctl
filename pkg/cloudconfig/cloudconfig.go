package cloudconfig

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
)

const (
	Header = "#cloud-config"
	Shell  = "/bin/bash"

	scriptDir = "/var/lib/cloud/scripts/per-instance/"

	defaultOwner             = "root:root"
	defaultScriptPermissions = "0755"
	defaultFilePermissions   = "0644"
)

type CloudConfig struct {
	Commands   []interface{} `json:"runcmd"`
	Packages   []string      `json:"packages"`
	WriteFiles []File        `json:"write_files"`
}

type File struct {
	Content     string `json:"content"`
	Owner       string `json:"owner"`
	Path        string `json:"path"`
	Permissions string `json:"permissions"`
}

func New() *CloudConfig {
	return &CloudConfig{}
}

func (c *CloudConfig) AddPackages(pkgs ...string) {
	c.Packages = append(c.Packages, pkgs...)
}

func (c *CloudConfig) AddCommands(cmds ...[]string) {
	for _, cmd := range cmds {
		c.Commands = append(c.Commands, cmd)
	}
}

func (c *CloudConfig) AddCommand(cmd ...string) {
	c.Commands = append(c.Commands, cmd)
}

func (c *CloudConfig) AddShellCommand(cmd string) {
	c.Commands = append(c.Commands, []string{Shell, "-c", cmd})
}

func (c *CloudConfig) AddFile(f File) {
	if f.Owner == "" {
		f.Owner = defaultOwner
	}
	if f.Permissions == "" {
		f.Permissions = defaultFilePermissions
	}
	c.WriteFiles = append(c.WriteFiles, f)
}

func (c *CloudConfig) AddScript(p, s string) {
	c.AddFile(File{
		Content:     s,
		Path:        p,
		Permissions: defaultScriptPermissions,
		Owner:       defaultOwner,
	})
}

func (c *CloudConfig) RunScript(name, s string) {
	p := scriptDir + name
	c.AddScript(p, s)
	c.AddCommand(p)
}

func (c *CloudConfig) Encode() (string, error) {
	buf := &bytes.Buffer{}

	data, err := yaml.Marshal(c)
	if err != nil {
		return "", err
	}

	data = append([]byte(fmt.Sprintln(Header)), data...)

	gw := gzip.NewWriter(buf)

	if _, err := gw.Write(data); err != nil {
		return "", err
	}
	gw.Close()
	data, err = ioutil.ReadAll(buf)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

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
	defer gr.Close()
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
