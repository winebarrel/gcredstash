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

func (c *GetCommand) parseArgs(args []string) (string, string, map[string]string, bool, bool, string, error) {
	argsWithoutN, noNL := gcredstash.HasOption(args, "-n")
	argsWithoutNS, noErr := gcredstash.HasOption(argsWithoutN, "-s")
	argsWithoutNSE, errOut, err := gcredstash.ParseOptionWithValue(argsWithoutNS, "-e")

	if err != nil {
		return "", "", nil, false, false, "", err
	}

	newArgs, version, err := gcredstash.ParseVersion(argsWithoutNSE)

	if err != nil {
		return "", "", nil, false, false, "", err
	}

	if len(newArgs) < 1 {
		return "", "", nil, false, false, "", fmt.Errorf("too few arguments")
	}

	credential := newArgs[0]
	context, err := gcredstash.ParseContext(newArgs[1:])

	return credential, version, context, noNL, noErr, errOut, err
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

func (c *GetCommand) write(filename string, message string) {
	fp, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)

	if err != nil {
		return
	}

	defer fp.Close()

	fp.WriteString(message)
}

func (c *GetCommand) RunImpl(args []string) (string, error) {
	credential, version, context, noNL, noErr, errOut, err := c.parseArgs(args)

	if err != nil {
		return "", err
	}

	if strings.Contains(credential, "*") {
		value, err := c.getCredentials(credential, version, context)

		if err != nil && errOut != "" {
			c.write(errOut, fmt.Sprintf("error: gcredstash get %v: %s\n", args, err.Error()))
		}

		return value, err
	} else {
		value, err := c.getCredential(credential, version, context)

		if err != nil {
			if errOut != "" {
				c.write(errOut, fmt.Sprintf("error: gcredstash get %v: %s\n", args, err.Error()))
			}

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
usage: gcredstash get [-v VERSION] [-n] [-s] [-e ERROUT] credential [context [context ...]]
`
	return strings.TrimSpace(helpText)
}
