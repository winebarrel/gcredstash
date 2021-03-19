package command

import (
	"github.com/mitchellh/cli"
	"github.com/kgaughan/gcredstash/src/gcredstash"
)

// Meta contain the meta-option that nearly all subcommand inherited.
type Meta struct {
	Ui      cli.Ui
	Table   string
	KmsKey  string
	Version string
	Driver  *gcredstash.Driver
}
