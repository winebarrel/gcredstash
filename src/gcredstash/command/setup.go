package command

import (
	"fmt"
	"gcredstash"
	"os"
	"strings"
)

type SetupCommand struct {
	Meta
}

func (c *SetupCommand) Run(args []string) int {
	if len(args) > 0 {
		fmt.Fprintf(os.Stderr, "error: too many arguments\n")
		return 1
	}

	createErr := gcredstash.CreateDdbTable(c.Meta.Table)

	if createErr != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", createErr.Error())
		return 1
	}

	return 0
}

func (c *SetupCommand) Synopsis() string {
	return "setup the credential store"
}

func (c *SetupCommand) Help() string {
	helpText := `
usage: credstash setup
`
	return strings.TrimSpace(helpText)
}
