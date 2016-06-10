package command

import (
	"encoding/json"
	"fmt"
	"gcredstash"
	"github.com/ryanuber/go-glob"
	"os"
	"strings"
)

type GetCommand struct {
	Meta
}

func (c *GetCommand) Run(args []string) int {
	argsWithoutN, noNL := gcredstash.HasOption(args, "-n")
	newArgs, version, parseErr := gcredstash.PerseVersion(argsWithoutN)

	if parseErr != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", parseErr.Error())
		return 1
	}

	if len(newArgs) < 1 {
		fmt.Fprintf(os.Stderr, "error: too few arguments\n")
		return 1
	}

	credential := newArgs[0]
	context, parseCtxErr := gcredstash.PerseContext(newArgs[1:])

	if parseCtxErr != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", parseCtxErr.Error())
		return 1
	}

	if strings.Contains(credential, "*") {
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
			if !glob.Glob(credential, name) {
				continue
			}

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
	} else {
		plainText, getSecErr := gcredstash.GetSecret(credential, version, c.Meta.Table, context)

		if getSecErr != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", getSecErr.Error())
			return 1
		}

		fmt.Print(plainText)

		if !noNL {
			fmt.Println()
		}
	}

	return 0
}

func (c *GetCommand) Synopsis() string {
	return "Get a credential from the store"
}

func (c *GetCommand) Help() string {
	helpText := `
usage: gcredstash get [-v VERSION] [-n] credential [context [context ...]]
`
	return strings.TrimSpace(helpText)
}
