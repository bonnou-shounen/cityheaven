package cmd

import (
	"fmt"

	"github.com/alecthomas/kong"
)

type PrintVersion struct{}

func (PrintVersion) Run(vars kong.Vars) error {
	fmt.Println(vars["version"])

	return nil
}
