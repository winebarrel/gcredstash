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

func (c *EnvCommand) Run(args []string) int {
	newArgs, version, parseErr := gcredstash.PerseVersion(args)

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

	names := map[string]bool{}

	if strings.Contains(credential, "*") {
		items, listErr := gcredstash.ListSecrets(c.Meta.Table)

		if listErr != nil {
			fmt.Fprintf(os.Stderr, "# error: %s\n", listErr.Error())
			return 1
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

		plainText, getSecErr := gcredstash.GetSecret(name, version, c.Meta.Table, context)

		if getSecErr != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", getSecErr.Error())
			continue
		}

		creds[name] = plainText
	}

	for name, value := range creds {
		name = convShellKeyword(name)
		value = escapeShellword(value)
		fmt.Printf("export %s=%s\n", name, value)
	}

	return 0
}

func (c *EnvCommand) Synopsis() string {
	return "Display the commands to set environment variables"
}

func (c *EnvCommand) Help() string {
	helpText := `
usage: gcredstash env [-v VERSION] credential [context [context ...]]
`
	return strings.TrimSpace(helpText)
}
