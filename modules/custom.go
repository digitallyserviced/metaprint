package modules

import (
	"os/exec"
	"strings"

	"github.com/gookit/goutil/dump"
	"github.com/oxodao/metaprint/utils"
)

type Custom struct {
	Prefix  string
	Suffix  string
	Command string
	Format  string
	Scripts map[string]Script
}

type Script struct {
	Inline  bool `yaml:"inline"`
	Content string
	Shell   []string
	Script  string
	Args    []string
}

func (c Custom) Print(args []string) string {
	if len(c.Scripts) > 0 {
		dump.P(c)
	}
	out, err := exec.Command("bash", "-c", c.Command).Output()
	if err != nil {
		return ""
	}

	return utils.ReplaceVariables(c.Format, map[string]interface{}{
		"output": strings.Trim(string(out), " \r\n\t"),
	})
	// return strings.Trim(string(out), " \n\t")
}

func (c Custom) GetPrefix() string {
	return c.Prefix
}

func (c Custom) GetSuffix() string {
	return c.Suffix
}
