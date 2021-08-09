package main

import (
	"github.com/ctrsploit/ctrsploit/auto"
	"github.com/urfave/cli/v2"
)

var autoCommand = &cli.Command{
	Name:    auto.ExpName,
	Aliases: []string{"a"},
	Usage:   "auto gathering information and exploit // TODO",
}
