package cmd

import (
	"fmt"
	"os"

	"github.com/bonnou-shounen/cityheaven"
)

type PrintVersion struct{}

func (*PrintVersion) Run() error {
	fmt.Fprintf(os.Stdout, "%s\n", cityheaven.Version)

	return nil
}
