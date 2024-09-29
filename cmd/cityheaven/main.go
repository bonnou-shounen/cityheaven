package main

import (
	"github.com/alecthomas/kong"

	"github.com/bonnou-shounen/cityheaven/internal/cmd"
)

func main() {
	cli := cmd.CLI{}
	ctx := kong.Parse(&cli,
		kong.ShortUsageOnError(),
	)

	ctx.FatalIfErrorf(ctx.Run(&cli))
}
