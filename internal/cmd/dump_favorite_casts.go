package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/bonnou-shounen/cityheaven"
	"github.com/bonnou-shounen/cityheaven/internal/util"
)

type DumpFavoriteCasts struct{}

func (d *DumpFavoriteCasts) Run(o *CLI) error {
	ctx := context.Background()

	client, err := util.NewLoggedClient(ctx, o.Dump.Fav.LoginID, o.Dump.Fav.Password)
	if err != nil {
		return fmt.Errorf("on NewLoggedClient(): %w", err)
	}

	casts, err := client.GetFavoriteCasts(ctx)
	if err != nil {
		return fmt.Errorf("on GetFavoriteCasts(): %w", err)
	}

	withFavCount := !o.Dump.Fav.Casts.NoFav
	if withFavCount {
		for _, cast := range casts {
			cast.FavCount, _ = client.GetFavoriteCount(ctx, cast)
		}
	}

	if o.Dump.JSON {
		return d.dumpJSON(casts, withFavCount)
	}

	return d.dumpTSV(casts, withFavCount)
}

func (d *DumpFavoriteCasts) dumpJSON(casts []*cityheaven.Cast, withFavCount bool) error {
	dumpCasts := make([]map[string]interface{}, len(casts))

	for i, cast := range casts {
		dumpCast := map[string]interface{}{
			"id":       cast.ID,
			"name":     cast.Name,
			"shopId":   cast.ShopID,
			"shopName": cast.ShopName,
		}

		if withFavCount {
			dumpCast["favCount"] = cast.FavCount
		}

		dumpCasts[i] = dumpCast
	}

	b, err := json.Marshal(dumpCasts)
	if err != nil {
		return fmt.Errorf("on json.Marshal(): %w", err)
	}

	fmt.Fprintf(os.Stdout, "%s", b)

	return nil
}

func (d *DumpFavoriteCasts) dumpTSV(casts []*cityheaven.Cast, withFavCount bool) error {
	for _, cast := range casts {
		line := fmt.Sprint(cast.ID, "\t", cast.ShopID)

		if withFavCount {
			line += fmt.Sprint("\t", cast.FavCount)
		}

		line += fmt.Sprint("\t", cast.Name, "\t", cast.ShopName)

		fmt.Fprintln(os.Stdout, line)
	}

	return nil
}
