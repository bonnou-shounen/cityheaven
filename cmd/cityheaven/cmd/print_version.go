package cmd

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
)

type PrintVersion struct{}

func (PrintVersion) Run(vars kong.Vars) error {
	fmt.Fprintf(os.Stdout, "%s\n", vars["version"])

	return nil
}
