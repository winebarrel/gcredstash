package gcredstash

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func PerseVersion(args []string) ([]string, string, error) {
	newArgs := []string{}
	version := ""
	next_version := false

	for _, arg := range args {
		if next_version {
			version = arg
			next_version = false
		} else if arg == "-v" {
			next_version = true
		} else {
			newArgs = append(newArgs, arg)
		}
	}

	if next_version {
		return nil, "", errors.New("option requires an argument -- v")
	}

	if version != "" {
		ver, atoiErr := strconv.Atoi(version)

		if atoiErr != nil {
			return nil, "", atoiErr
		}

		version = fmt.Sprintf("%019d", ver)
	}

	return newArgs, version, nil
}

func PerseContext(contextStrs []string) (map[string]string, error) {
	context := map[string]string{}

	for _, ctxStr := range contextStrs {
		kv := strings.SplitN(ctxStr, "=", 2)

		if len(kv) < 2 || kv[0] == "" || kv[1] == "" {
			return nil, fmt.Errorf("invalid context -- %s", ctxStr)
		}

		context[kv[0]] = kv[1]
	}

	return context, nil
}
