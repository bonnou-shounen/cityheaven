package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/bonnou-shounen/cityheaven"
)

type DumpShopCasts struct {
	Attendees bool `help:"only attendees"`
}

func (d *DumpShopCasts) Run(o *CLI) error {
	strURL, err := o.Dump.Shop.GetURL()
	if err != nil {
		return err
	}

	ctx := context.Background()
	client := cityheaven.NewClient()

	var casts []*cityheaven.Cast

	if d.Attendees {
		casts, err = client.GetShopAttendees(ctx, strURL)
		if err != nil {
			return fmt.Errorf(`on GetShopAtendees("%s"): %w`, strURL, err)
		}
	} else {
		casts, err = client.GetShopCasts(ctx, strURL)
		if err != nil {
			return fmt.Errorf(`on GetShopCasts("%s"): %w`, strURL, err)
		}
	}

	withFavCount := !o.Dump.Shop.Casts.NoFav
	if withFavCount {
		for _, cast := range casts {
			cast.FavCount, _ = client.GetFavoriteCount(ctx, cast)
		}
	}

	if o.Dump.JSON {
		return d.dumpJSON(casts, withFavCount, d.Attendees)
	}

	return d.dumpTSV(casts, withFavCount, d.Attendees)
}

func (d *DumpShopCasts) dumpJSON(casts []*cityheaven.Cast, withFavCount, withNextStart bool) error {
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

		if withNextStart {
			dumpCast["nextStart"] = cast.NextStart
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

func (d *DumpShopCasts) dumpTSV(casts []*cityheaven.Cast, withFavCount, withNextStart bool) error {
	for _, cast := range casts {
		line := fmt.Sprint(cast.ID, "\t", cast.ShopID)

		if withFavCount {
			line += fmt.Sprint("\t", cast.FavCount)
		}

		if withNextStart {
			line += fmt.Sprint("\t", cast.NextStart)
		}

		line += fmt.Sprint("\t", cast.Name, "\t", cast.ShopName)

		fmt.Fprintln(os.Stdout, line)
	}

	return nil
}
