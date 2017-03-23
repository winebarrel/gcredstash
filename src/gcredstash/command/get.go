package command

import (
	"fmt"
	"gcredstash"
	"github.com/ryanuber/go-glob"
	"os"
	"strings"
)

type GetCommand struct {
	Meta
}

func (c *GetCommand) parseArgs(args []string) (string, string, map[string]string, bool, bool, error) {
	argsWithoutN, noNL := gcredstash.HasOption(args, "-n")
	argsWithoutNS, noErr := gcredstash.HasOption(argsWithoutN, "-s")
	newArgs, version, err := gcredstash.ParseVersion(argsWithoutNS)

	if err != nil {
		return "", "", nil, false, false, err
	}

	if len(newArgs) < 1 {
		return "", "", nil, false, false, fmt.Errorf("too few arguments")
	}

	credential := newArgs[0]
	context, err := gcredstash.ParseContext(newArgs[1:])

	return credential, version, context, noNL, noErr, err
}

func (c *GetCommand) getCredential(credential string, version string, context map[string]string) (string, error) {
	value, err := c.Driver.GetSecret(credential, version, c.Table, context)

	if err != nil {
		return "", err
	}

	return value, nil
}

func (c *GetCommand) getCredentials(credential string, version string, context map[string]string) (string, error) {
	names := map[string]bool{}
	items, err := c.Driver.ListSecrets(c.Table)

	if err != nil {
		return "", err
	}

	for name, _ := range items {
		names[*name] = true
	}

	creds := map[string]string{}

	for name, _ := range names {
		if !glob.Glob(credential, name) {
			continue
		}

		value, err := c.Driver.GetSecret(name, version, c.Table, context)

		if err != nil {
			continue
		}

		creds[name] = value
	}

	return gcredstash.MapToJson(creds) + "\n", nil
}

func (c *GetCommand) RunImpl(args []string) (string, error) {
	credential, version, context, noNL, noErr, err := c.parseArgs(args)

	if err != nil {
		return "", err
	}

	if strings.Contains(credential, "*") {
		return c.getCredentials(credential, version, context)
	} else {
		value, err := c.getCredential(credential, version, context)

		if err != nil {
			if noErr {
				return "", nil
			} else {
				return "", err
			}
		}

		if noNL {
			return value, nil
		} else {
			return value + "\n", nil
		}
	}
}

func (c *GetCommand) Run(args []string) int {
	out, err := c.RunImpl(args)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return 1
	}

	fmt.Print(out)

	return 0
}

func (c *GetCommand) Synopsis() string {
	return "Get a credential from the store"
}

func (c *GetCommand) Help() string {
	helpText := `
usage: gcredstash get [-v VERSION] [-n] [-s] credential [context [context ...]]
`
	return strings.TrimSpace(helpText)
}
