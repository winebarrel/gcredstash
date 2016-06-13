package command

import (
	"fmt"
	"os"
	"strings"
)

type SetupCommand struct {
	Meta
}

func (c *SetupCommand) RunImpl(args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("too many arguments")
	}

	err := c.Driver.CreateDdbTable(c.Meta.Table)

	if err != nil {
		return err
	}

	return nil
}

func (c *SetupCommand) Run(args []string) int {
	err := c.RunImpl(args)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return 1
	}

	return 0
}

func (c *SetupCommand) Synopsis() string {
	return "setup the credential store"
}

func (c *SetupCommand) Help() string {
	helpText := `
usage: gcredstash setup
`
	return strings.TrimSpace(helpText)
}
