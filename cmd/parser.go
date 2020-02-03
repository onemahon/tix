package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"tix/env"
)

const versionUsage = "prints tix version"

type Parser struct {
	env     []string
	flag    *flag.FlagSet
	version string
}


func NewParser(env []string, version string) *Parser {
	name := os.Args[0]
	flagSet := flag.NewFlagSet(name, flag.ExitOnError)
	flagSet.Usage = func() {
		_, _ = fmt.Fprintf(flagSet.Output(),
			"Usage: %s [OPTIONS] <markdown file> \n",
			name)
		flagSet.PrintDefaults()
	}
	return &Parser{env: env, flag: flagSet, version: version}
}


func (parser *Parser) Parse() Command {
	version := parser.setupVersion()

	_ = parser.flag.Parse(os.Args[1:])
	if *version {
		return NewVersionCommand(parser.version)
	} else {
		if parser.flag.NArg() < 1 {
			parser.printUsageAndExit()
		}
		path, _ := parser.localPath()
		return NewTixCommand(env.Map(), path)
	}
}

func (parser *Parser) setupVersion() *bool {
	var verbose bool
	parser.flag.BoolVar(&verbose, "version", false, versionUsage)

	return &verbose
}

func (parser *Parser) printUsageAndExit() {
	parser.flag.Usage()
	os.Exit(2)
}

func (parser *Parser) localPath() (string, error) {
	relativePath := parser.flag.Arg(0)
	return filepath.Abs(relativePath)
}