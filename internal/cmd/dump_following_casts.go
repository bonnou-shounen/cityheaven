package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/bonnou-shounen/cityheaven"
	"github.com/bonnou-shounen/cityheaven/internal/util"
)

type DumpFollowingCasts struct {
	Mutual bool `help:"only mutual followees"`
}

func (d *DumpFollowingCasts) Run(o *CLI) error {
	ctx := context.Background()

	client, err := util.NewLoggedClient(ctx, o.Dump.Follow.LoginID, o.Dump.Follow.Password)
	if err != nil {
		return fmt.Errorf("on NewLoggedClient(): %w", err)
	}

	casts, err := client.GetFollowingCasts(ctx)
	if err != nil {
		return fmt.Errorf("on GetFavoriteCasts(): %w", err)
	}

	withFavCount := !o.Dump.Follow.Casts.NoFav
	if withFavCount {
		for _, cast := range casts {
			if d.Mutual && !cast.MutualFollow {
				continue
			}

			cast.FavCount, _ = client.GetFavoriteCount(ctx, cast)
		}
	}

	if o.Dump.JSON {
		return d.dumpJSON(casts, d.Mutual, withFavCount)
	}

	return d.dumpTSV(casts, d.Mutual, withFavCount)
}

func (d *DumpFollowingCasts) dumpJSON(casts []*cityheaven.Cast, onlyMutual, withFavCount bool) error {
	dumpCasts := make([]map[string]interface{}, len(casts))

	for i, cast := range casts {
		if onlyMutual && !cast.MutualFollow {
			continue
		}

		dumpCast := map[string]interface{}{
			"id":       cast.ID,
			"name":     cast.Name,
			"shopId":   cast.ShopID,
			"shopName": cast.ShopName,
		}

		if !onlyMutual {
			dumpCast["isMutual"] = cast.MutualFollow
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

func (d *DumpFollowingCasts) dumpTSV(casts []*cityheaven.Cast, onlyMutual, withFavCount bool) error {
	for _, cast := range casts {
		if onlyMutual && !cast.MutualFollow {
			continue
		}

		line := fmt.Sprint(cast.ID, "\t", cast.ShopID)

		if !onlyMutual {
			var followType string
			if cast.MutualFollow {
				followType = "M"
			} else {
				followType = "F"
			}

			line += fmt.Sprint("\t", followType)
		}

		if withFavCount {
			line += fmt.Sprint("\t", cast.FavCount)
		}

		line += fmt.Sprint("\t", cast.Name, "\t", cast.ShopName)

		fmt.Fprintln(os.Stdout, line)
	}

	return nil
}
