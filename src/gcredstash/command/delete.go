package command

import (
	"fmt"
	"gcredstash"
	"os"
	"strings"
)

type DeleteCommand struct {
	Meta
}

func (c *DeleteCommand) Run(args []string) int {
	newArgs, version, err := gcredstash.PerseVersion(args)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return 1
	}

	if len(newArgs) < 1 {
		fmt.Fprintf(os.Stderr, "error: too few arguments\n")
		return 1
	}

	credential := args[0]

	err = gcredstash.DeleteSecrets(credential, version, c.Meta.Table)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return 1
	}

	return 0
}

func (c *DeleteCommand) Synopsis() string {
	return "Delete a credential from the store"
}

func (c *DeleteCommand) Help() string {
	helpText := `
usage: gcredstash delete [-v VERSION] credential
`
	return strings.TrimSpace(helpText)
}
