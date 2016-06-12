package command

import (
	"fmt"
	"gcredstash"
	"os"
	"strings"
)

type PutCommand struct {
	Meta
}

func (c *PutCommand) parseArgs(args []string) (string, string, string, map[string]string, bool, error) {
	argsWithoutA, autoVersion := gcredstash.HasOption(args, "-a")
	newArgs, version, err := gcredstash.ParseVersion(argsWithoutA)

	if err != nil {
		return "", "", "", nil, false, err
	}

	if len(newArgs) < 2 {
		return "", "", "", nil, false, fmt.Errorf("too few arguments")
	}

	credential := newArgs[0]
	value := newArgs[1]
	context, err := gcredstash.ParseContext(newArgs[2:])

	return credential, value, version, context, autoVersion, err
}

func (c *PutCommand) RunImpl(args []string) error {
	credential, value, version, context, autoVersion, err := c.parseArgs(args)

	if err != nil {
		return err
	}

	if value == "-" {
		value = gcredstash.ReadStdin()
	}

	if autoVersion {
		latestVersion, err := c.Driver.GetHighestVersion(credential, c.Table)

		if err != nil {
			return err
		}

		latestVersion += 1
		version = gcredstash.VersionNumToStr(latestVersion)
	} else if version == "" {
		version = gcredstash.VersionNumToStr(1)
	}

	err = c.Driver.PutSecret(credential, value, version, c.KmsKey, c.Table, context)

	if err != nil {
		return err
	}

	fmt.Printf("%s has been stored\n", credential)

	return nil
}

func (c *PutCommand) Run(args []string) int {
	err := c.RunImpl(args)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return 1
	}

	return 0
}

func (c *PutCommand) Synopsis() string {
	return "Put a credential into the store"
}

func (c *PutCommand) Help() string {
	helpText := `
usage: gcredstash put [-k KEY] [-v VERSION] [-a] credential value [context [context ...]]
`
	return strings.TrimSpace(helpText)
}
