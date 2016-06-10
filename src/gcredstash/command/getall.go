package command

import (
	"encoding/json"
	"fmt"
	"gcredstash"
	"os"
	"strings"
)

type GetallCommand struct {
	Meta
}

func (c *GetallCommand) Run(args []string) int {
	newArgs, version, err := gcredstash.PerseVersion(args)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return 1
	}

	context, err := gcredstash.PerseContext(newArgs)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return 1
	}

	names := map[string]bool{}

	items, err := gcredstash.ListSecrets(c.Meta.Table)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return 1
	}

	for name, _ := range items {
		names[*name] = true
	}

	creds := map[string]string{}
	hasErr := false

	for name, _ := range names {
		plainText, err := gcredstash.GetSecret(name, version, c.Meta.Table, context)

		if err != nil {
			hasErr = true
			fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
			continue
		}

		creds[name] = plainText
	}

	jsonString, err := json.MarshalIndent(creds, "", "  ")

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return 1
	}

	fmt.Println(string(jsonString))

	if hasErr {
		return 1
	}

	return 0
}

func (c *GetallCommand) Synopsis() string {
	return "Get all credentials from the store"
}

func (c *GetallCommand) Help() string {
	helpText := `
usage: gcredstash getall [-v VERSION] [context [context ...]]
`
	return strings.TrimSpace(helpText)
}
