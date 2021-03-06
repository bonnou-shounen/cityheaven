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
		fields = append(fields, "", "", "", "", "")

		castID, _ := strconv.Atoi(fields[0])
		shopID, _ := strconv.Atoi(fields[1])
		// fields[2] is favCount
		castName := fields[3]
		shopName := fields[4]

		if castID == 0 || shopID == 0 {
			continue
		}

		casts = append(casts,
			&cityheaven.Cast{
				ID:       castID,
				Name:     castName,
				ShopID:   shopID,
				ShopName: shopName,
			},
		)
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
	oldCasts := map[int]*cityheaven.Cast{}

	for _, curCast := range curCasts {
		oldCasts[curCast.ID] = curCast
	}

	for _, newCast := range newCasts {
		if _, exists := oldCasts[newCast.ID]; exists {
			delete(oldCasts, newCast.ID)
		} else {
			addCasts = append(addCasts, newCast)
		}
	}

	for _, oldCast := range oldCasts {
		delCasts = append(delCasts, oldCast)
	}

	return delCasts, addCasts
}
