package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/bonnou-shounen/cityheaven/cmd/cityheaven/util"
)

type DumpFavoriteCasts struct {
	NoFav bool `help:"skip counting favorites"`
}

func (d *DumpFavoriteCasts) Run() error {
	ctx := context.Background()

	client, err := util.NewLoggedClient(ctx)
	if err != nil {
		return fmt.Errorf("on NewLoggedClient(): %w", err)
	}

	casts, err := client.GetFavoriteCasts(ctx)
	if err != nil {
		return fmt.Errorf("on GetFavoriteCasts(): %w", err)
	}

	for _, cast := range casts {
		var favCount int
		if !d.NoFav {
			favCount, _ = client.GetFavoriteCount(ctx, cast)
		}

		fmt.Fprintf(os.Stdout, "%d\t%d\t%d\t%s\t%s\n", cast.ID, cast.ShopID, favCount, cast.Name, cast.ShopName)
	}

	return nil
}
