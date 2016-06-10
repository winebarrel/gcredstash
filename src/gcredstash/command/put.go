package command

import (
	"bufio"
	"fmt"
	"gcredstash"
	"io/ioutil"
	"os"
	"strings"
)

type PutCommand struct {
	Meta
}

func readStdin() string {
	reader := bufio.NewReader(os.Stdin)
	input, err := ioutil.ReadAll(reader)

	if err != nil {
		panic(err)
	}

	return strings.TrimRight(string(input), "\n")
}

func (c *PutCommand) Run(args []string) int {
	argsWithoutA, autoVersion := gcredstash.HasOption(args, "-a")
	newArgs, version, err := gcredstash.PerseVersion(argsWithoutA)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return 1
	}

	if len(newArgs) < 2 {
		fmt.Fprintf(os.Stderr, "error: too few arguments\n")
		return 1
	}

	credential := newArgs[0]
	value := newArgs[1]
	context, err := gcredstash.PerseContext(newArgs[2:])

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return 1
	}

	if value == "-" {
		value = readStdin()
	}

	if autoVersion {
		latestVersion, err := gcredstash.GetHighestVersion(credential, c.Meta.Table)

		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
			return 1
		}

		latestVersion += 1
		version = fmt.Sprintf("%019d", latestVersion)
	} else if version == "" {
		version = fmt.Sprintf("%019d", 1)
	}

	err = gcredstash.PutSecret(credential, value, version, c.Meta.KmsKey, c.Meta.Table, context)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return 1

	}

	fmt.Printf("%s has been stored\n", credential)

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
