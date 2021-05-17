package main

import (
	"os"

	"github.com/alecthomas/kong"

	"github.com/bonnou-shounen/cityheaven/cmd/cityheaven/cmd"
)

func main() {
	arg := cmd.Arg{}
	ctx := kong.Parse(
		&arg,
		kong.Name("cityheaven"),
		kong.Vars{"version": "0.0.8"},
		kong.UsageOnError(),
	)

	if arg.Option.Login != "" {
		os.Setenv("CITYHEAVEN_LOGIN", arg.Option.Login)
	}

	if arg.Option.Password != "" {
		os.Setenv("CITYHEAVEN_PASSWORD", arg.Option.Password)
	}

	ctx.FatalIfErrorf(ctx.Run(&arg.Option))
}
