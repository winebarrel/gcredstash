package command

import (
	"fmt"
	"gcredstash"
	"os"
	"strings"
)

type GetallCommand struct {
	Meta
}

func (c *GetallCommand) getNames() ([]string, error) {
	namesMap := map[string]bool{}
	names := []string{}

	items, err := c.Driver.ListSecrets(c.Table)

	if err != nil {
		return nil, err
	}

	for name, _ := range items {
		namesMap[*name] = true
	}

	for name, _ := range namesMap {
		names = append(names, name)
	}

	return names, nil
}

func (c *GetallCommand) getCredentials(names []string, context map[string]string) map[string]string {
	creds := map[string]string{}

	for _, name := range names {
		value, err := c.Driver.GetSecret(name, "", c.Table, context)

		if err != nil {
			continue
		}

		creds[name] = value
	}

	return creds
}

func (c *GetallCommand) RunImpl(args []string) (string, error) {
	context, err := gcredstash.ParseContext(args)

	if err != nil {
		return "", err
	}

	names, err := c.getNames()

	if err != nil {
		return "", err
	}

	creds := c.getCredentials(names, context)

	return gcredstash.MapToJson(creds) + "\n", nil
}

func (c *GetallCommand) Run(args []string) int {
	out, err := c.RunImpl(args)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return 1
	}

	fmt.Print(out)

	return 0
}

func (c *GetallCommand) Synopsis() string {
	return "Get all credentials from the store"
}

func (c *GetallCommand) Help() string {
	helpText := `
usage: gcredstash getall [context [context ...]]
`
	return strings.TrimSpace(helpText)
}
