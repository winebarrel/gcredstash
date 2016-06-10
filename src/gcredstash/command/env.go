package command

import (
	"fmt"
	"gcredstash"
	"github.com/ryanuber/go-glob"
	"os"
	"regexp"
	"strings"
)

type EnvCommand struct {
	Meta
}

func escapeShellword(word string) string {
	rep := regexp.MustCompile(`([^A-Za-z0-9_\-.,:\/@\n])`)
	return rep.ReplaceAllString(word, `\$1`)
}

func convShellKeyword(word string) string {
	rep := regexp.MustCompile(`([^A-Za-z0-9_]+)`)
	name := rep.ReplaceAllString(word, "_")
	return strings.ToUpper(name)
}

func parseArgs(args []string) ([]string, string, string, error) {
	argsWithoutPrefix, prefix, err := gcredstash.GetOptionValue(args, "-p")

	if err != nil {
		return nil, "", "", err
	}

	newArgs, version, err := gcredstash.PerseVersion(argsWithoutPrefix)

	if err != nil {
		return nil, "", "", err
	}

	return newArgs, version, prefix, nil
}

func getCredentials(credential string, version string, table string, context map[string]string) (map[string]string, error) {
	names := map[string]bool{}

	if strings.Contains(credential, "*") {
		items, err := gcredstash.ListSecrets(table)

		if err != nil {
			return nil, err
		}

		for name, _ := range items {
			names[*name] = true
		}
	} else {
		names[credential] = true
	}

	creds := map[string]string{}

	for name, _ := range names {
		if !glob.Glob(credential, name) {
			continue
		}

		plainText, err := gcredstash.GetSecret(name, version, table, context)

		if err != nil {
			fmt.Fprintf(os.Stderr, "# error: %s\n", err.Error())
			continue
		}

		creds[name] = plainText
	}

	return creds, nil
}

func printEnvs(creds map[string]string, prefix string) {
	for name, value := range creds {
		if strings.HasPrefix(name, prefix) {
			name = name[len(prefix):]
		}

		name = convShellKeyword(name)
		value = escapeShellword(value)
		fmt.Printf("export %s=%s\n", name, value)
	}
}

func (c *EnvCommand) Run(args []string) int {
	newArgs, version, prefix, err := parseArgs(args)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return 1
	}

	if len(newArgs) < 1 {
		fmt.Fprintf(os.Stderr, "error: too few arguments\n")
		return 1
	}

	credential := args[0]
	context, err := gcredstash.PerseContext(newArgs[1:])

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return 1
	}

	creds, err := getCredentials(credential, version, c.Meta.Table, context)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return 1
	}

	printEnvs(creds, prefix)

	return 0
}

func (c *EnvCommand) Synopsis() string {
	return "Display the commands to set environment variables"
}

func (c *EnvCommand) Help() string {
	helpText := `
usage: gcredstash env [-v VERSION] [-p PREFIX] credential [context [context ...]]
`
	return strings.TrimSpace(helpText)
}
