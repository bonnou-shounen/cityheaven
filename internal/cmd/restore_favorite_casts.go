package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/bonnou-shounen/cityheaven"
	"github.com/bonnou-shounen/cityheaven/internal/util"
)

type RestoreFavoriteCasts struct{}

func (r *RestoreFavoriteCasts) Run(o *CLI) error {
	newCasts := r.readCasts(os.Stdin)

	ctx := context.Background()

	client, err := util.NewLoggedClient(ctx, o.Restore.Fav.LoginID, o.Restore.Fav.Password)
	if err != nil {
		return fmt.Errorf("on NewLoggedClient(): %w", err)
	}

	curCasts, err := client.GetFavoriteCasts(ctx)
	if err != nil {
		return fmt.Errorf("on GetFavoriteCasts(): %w", err)
	}

	if r.areSame(curCasts, newCasts) {
		return nil
	}

	delCasts, addCasts := r.castsDiff(curCasts, newCasts)

	if err := client.DeleteFavoriteCasts(ctx, delCasts); err != nil {
		return fmt.Errorf("on DeleteFavoriteCasts(): %w", err)
	}

	if err := client.AddFavoriteCasts(ctx, addCasts); err != nil {
		return fmt.Errorf("on AddFavoriteCasts(): %w", err)
	}

	if err := client.SortFavoriteCasts(ctx, newCasts); err != nil {
		return fmt.Errorf("on SortFavoriteCasts(): %w", err)
	}

	return nil
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

func (r *RestoreFavoriteCasts) castsDiff(curCasts, newCasts []*cityheaven.Cast) (delCasts, addCasts []*cityheaven.Cast) { //nolint:lll,nonamedreturns
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
