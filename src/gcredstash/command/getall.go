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
	newArgs, version, parseErr := gcredstash.PerseVersion(args)

	if parseErr != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", parseErr.Error())
		return 1
	}

	context, parseCtxErr := gcredstash.PerseContext(newArgs)

	if parseCtxErr != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", parseCtxErr.Error())
		return 1
	}

	names := map[string]bool{}

	items, listErr := gcredstash.ListSecrets(c.Meta.Table)

	if listErr != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", listErr.Error())
		return 1
	}

	for name, _ := range items {
		names[*name] = true
	}

	creds := map[string]string{}
	hasErr := false

	for name, _ := range names {
		plainText, getSecErr := gcredstash.GetSecret(name, version, c.Meta.Table, context)

		if getSecErr != nil {
			hasErr = true
			fmt.Fprintf(os.Stderr, "error: %s\n", getSecErr.Error())
			continue
		}

		creds[name] = plainText
	}

	jsonString, jsonErr := json.MarshalIndent(creds, "", "  ")

	if jsonErr != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", jsonErr.Error())
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
