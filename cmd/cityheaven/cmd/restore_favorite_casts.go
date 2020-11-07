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

		if shopID == 0 || castID == 0 {
			continue
		}

		casts = append(casts, &cityheaven.Cast{ShopID: shopID, CastID: castID})
	}

	return casts
}

func (r *RestoreFavoriteCasts) areSame(curCasts, newCasts []*cityheaven.Cast) bool {
	if len(curCasts) != len(newCasts) {
		return false
	}

	for i := range curCasts {
		if curCasts[i].CastID != newCasts[i].CastID {
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
			if newCast.CastID == curCast.CastID {
				continue LA
			}
		}
		addCasts = append(addCasts, newCast)
	}

LD:
	for _, curCast := range curCasts {
		for _, newCast := range newCasts {
			if curCast.CastID == newCast.CastID {
				continue LD
			}
		}
		delCasts = append(delCasts, curCast)
	}

	return delCasts, addCasts
}
