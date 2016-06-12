package command

import (
	"fmt"
	"gcredstash"
	"os"
	"sort"
	"strings"
)

type ListCommand struct {
	Meta
}

func (c *ListCommand) getLines(items map[*string]*string) []string {
	maxNameLen := gcredstash.MaxKeyLen(items)
	lines := []string{}

	for name, version := range items {
		versionNum := gcredstash.Atoi(*version)
		lines = append(lines, fmt.Sprintf("%-*s -- version: %d", maxNameLen, *name, versionNum))
	}

	return lines
}

func (c *ListCommand) RunImpl(args []string) (string, error) {
	if len(args) > 0 {
		return "", fmt.Errorf("too many arguments")
	}

	items, err := c.Driver.ListSecrets(c.Table)

	if err != nil {
		return "", err
	}

	lines := c.getLines(items)
	sort.Strings(lines)

	return strings.Join(lines, "\n"), nil
}

func (c *ListCommand) Run(args []string) int {
	out, err := c.RunImpl(args)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return 1
	}

	fmt.Println(out)

	return 0
}

func (c *ListCommand) Synopsis() string {
	return "list credentials and their version"
}

func (c *ListCommand) Help() string {
	helpText := `
usage: gcredstash list
`

	return strings.TrimSpace(helpText)
}
