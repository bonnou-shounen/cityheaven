package cmd

import (
	"fmt"
	"os"

	"github.com/bonnou-shounen/cityheaven/cmd/cityheaven/util"
)

type DumpFavoriteCasts struct {
	NoFav bool `help:"skip counts"`
}

func (d *DumpFavoriteCasts) Run() error {
	c, err := util.NewLoggedClient()
	if err != nil {
		return err
	}

	casts, err := c.GetFavoriteCasts()
	if err != nil {
		return err
	}

	for _, cast := range casts {
		var favCount int
		if !d.NoFav {
			favCount, _ = c.GetFavoriteCount(cast)
		}

		fmt.Fprintf(os.Stdout, "%d\t%d\t%d\t%s\t%s\n", cast.ShopID, cast.ID, favCount, cast.Name, cast.ShopName)
	}

	return nil
}
