package cmd

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/bonnou-shounen/cityheaven"
	"github.com/bonnou-shounen/cityheaven/cmd/cityheaven/util"
)

type RestoreFavoriteCasts struct{}

func (r *RestoreFavoriteCasts) Run() error {
	newCasts := r.readCasts(os.Stdin)

	c, err := util.NewLoggedClient()
	if err != nil {
		return err
	}

	curCasts, err := c.GetFavoriteCasts()
	if err != nil {
		return err
	}

	if r.areSame(curCasts, newCasts) {
		return nil
	}

	delCasts, addCasts := r.castsDiff(curCasts, newCasts)
	c.DeleteFavoriteCasts(delCasts) //nolint:errcheck
	c.AddFavoriteCasts(addCasts)    //nolint:errcheck

	return c.SortFavoriteCasts(newCasts)
}

func (r *RestoreFavoriteCasts) readCasts(reader io.Reader) []*cityheaven.Cast {
	casts := make([]*cityheaven.Cast, 0)

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			continue
		}

		shopID, _ := strconv.Atoi(fields[0])
		castID, _ := strconv.Atoi(fields[1])

		if castID != 0 && shopID != 0 {
			casts = append(casts, &cityheaven.Cast{ID: castID, ShopID: shopID})
		}
	}

	return casts
}

func (r *RestoreFavoriteCasts) areSame(curCasts, newCasts []*cityheaven.Cast) bool {
	if len(curCasts) != len(newCasts) {
		return false
	}

	for i := range curCasts {
		if curCasts[i].ID != newCasts[i].ID {
			return false
		}
	}

	return true
}

//nolint:lll
func (r *RestoreFavoriteCasts) castsDiff(curCasts, newCasts []*cityheaven.Cast) (delCasts, addCasts []*cityheaven.Cast) {
LA:
	for _, newCast := range newCasts {
		for _, curCast := range curCasts {
			if newCast.ID == curCast.ID {
				continue LA
			}
		}
		addCasts = append(addCasts, newCast)
	}

LD:
	for _, curCast := range curCasts {
		for _, newCast := range newCasts {
			if curCast.ID == newCast.ID {
				continue LD
			}
		}
		delCasts = append(delCasts, curCast)
	}

	return delCasts, addCasts
}
