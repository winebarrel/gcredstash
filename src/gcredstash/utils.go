package gcredstash

import (
	"fmt"
	"strconv"
	"strings"
)

func GetOptionValue(args []string, opt string) ([]string, string, error) {
	newArgs := []string{}
	val := ""
	nextOpt := false

	for _, arg := range args {
		if nextOpt {
			val = arg
			nextOpt = false
		} else if arg == opt {
			nextOpt = true
		} else {
			newArgs = append(newArgs, arg)
		}
	}

	if nextOpt {
		return nil, "", fmt.Errorf("option requires an argument -- %s", opt)
	}

	return newArgs, val, nil
}

func PerseVersion(args []string) ([]string, string, error) {
	newArgs, version, err := GetOptionValue(args, "-v")

	if err != nil {
		return nil, "", err
	}

	if version != "" {
		ver, err := strconv.Atoi(version)

		if err != nil {
			return nil, "", err
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

func HasOption(args []string, opt string) ([]string, bool) {
	newArgs := []string{}
	hasOpt := false

	for _, arg := range args {
		if arg == opt {
			hasOpt = true
		} else {
			newArgs = append(newArgs, arg)
		}
	}

	return newArgs, hasOpt
}
