package main

import (
	"gcredstash/command"
	"github.com/mitchellh/cli"
)

func Commands(meta *command.Meta) map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"delete": func() (cli.Command, error) {
			return &command.DeleteCommand{
				Meta: *meta,
			}, nil
		},
		"get": func() (cli.Command, error) {
			return &command.GetCommand{
				Meta: *meta,
			}, nil
		},
		"getall": func() (cli.Command, error) {
			return &command.GetallCommand{
				Meta: *meta,
			}, nil
		},
		"list": func() (cli.Command, error) {
			return &command.ListCommand{
				Meta: *meta,
			}, nil
		},
		"put": func() (cli.Command, error) {
			return &command.PutCommand{
				Meta: *meta,
			}, nil
		},
		"setup": func() (cli.Command, error) {
			return &command.SetupCommand{
				Meta: *meta,
			}, nil
		},
	}
}
