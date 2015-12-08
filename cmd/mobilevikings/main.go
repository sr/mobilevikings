package main

import (
	"fmt"
	"os"

	"github.com/sr/mobilevikings"
	"github.com/sr/mobilevikings/dumper"
	"go.pedge.io/env"
)

type cmdEnv struct {
	Token     string `env:"MOBILE_VIKINGS_TOKEN,required"`
	Directory string `env:"MOBILE_VIKINGS_DIRECTORY,required"`
}

func run() error {
	cmdEnv := &cmdEnv{}
	if err := env.Populate(cmdEnv); err != nil {
		return err
	}

	client := mobilevikings.NewClient(cmdEnv.Token)
	dumper := dumper.NewDumper(client, cmdEnv.Directory)
	return dumper.Dump()
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
