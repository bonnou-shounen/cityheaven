package cityheaven

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/jesse0michael/errgroup"
)

func (c *Client) GetFollowingCasts(ctx context.Context) ([]*Cast, error) {
	var lastPage int

	casts, err := c.getFollowingCastsOnPage(ctx, 1, &lastPage)
	if err != nil {
		return nil, fmt.Errorf("on getFollowingCastsOnPage(1): %w", err)
	}

	if lastPage >= 2 {
		eg, egCtx := errgroup.WithContext(ctx, 3)

		castsOnPage := make([][]*Cast, lastPage+1)

		for page := 2; page <= lastPage; page++ {
			page := page

			eg.Go(func() error {
				casts, err := c.getFollowingCastsOnPage(egCtx, page, nil)
				if err != nil {
					return fmt.Errorf("on getFollowingCastsOnPage(%d): %w", page, err)
				}

				castsOnPage[page] = casts

				return nil
			})
		}

		if err := eg.Wait(); err != nil {
			return nil, fmt.Errorf("on goroutine: %w", err)
		}

		for page := 2; page <= lastPage; page++ {
			casts = append(casts, castsOnPage[page]...)
		}
	}

	casts, err = c.fillShopInfo(ctx, casts)
	if err != nil {
		return nil, fmt.Errorf("on fillShopInfo(): %w", err)
	}

	return casts, nil
}

func (c *Client) getFollowingCastsOnPage(ctx context.Context, page int, pLastPage *int) ([]*Cast, error) {
	strURL := fmt.Sprint("https://www.cityheaven.net/tt/community/ABFollowGirlList/?start=", page)

	resp, err := c.get(ctx, strURL, "")
	if err != nil {
		return nil, fmt.Errorf(`on get("%s"): %w`, strURL, err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("on NewDocumentFromReader(): %w", err)
	}

	var casts []*Cast

	doc.Find("ul.follower_box li").Each(func(_ int, li *goquery.Selection) {
		castName := li.Find("p.name").Text()
		shopName := li.Find("p.age").Text()

		href, _ := li.Find("a").Attr("href")
		castID := c.parseNumber(href, "girlid-", "/")

		pathURL := href
		if i := strings.Index(href, "girlid-"); i > 0 {
			pathURL = href[:i]
		}

		mutualFollow := strings.HasPrefix(li.Find("div.btn_follo").Text(), "相互")

		if castID != 0 && castName != "" {
			casts = append(casts,
				&Cast{
					ID:           castID,
					Name:         castName,
					ShopName:     shopName,
					PathURL:      pathURL,
					MutualFollow: mutualFollow,
				},
			)
		}
	})

	if pLastPage != nil {
		num := c.parseNumber(doc.Find("span.fov-cnt").Text(), "全", "人")
		*pLastPage = (num + 29) / 30
	}

	return casts, nil
}

func (c *Client) fillShopInfo(ctx context.Context, casts []*Cast) ([]*Cast, error) {
	cache := map[string]*Cast{}

	for _, cast := range casts {
		var shop *Cast

		if cached, ok := cache[cast.PathURL]; ok {
			shop = cached
		} else {
			strURL := fmt.Sprint("https://www.cityheaven.net", cast.PathURL)

			got, err := c.getShopFromPage(ctx, strURL)
			if err != nil {
				return nil, fmt.Errorf("on getShopFromPage(): %w", err)
			}

			shop, cache[cast.PathURL] = got, got
		}

		cast.ShopID = shop.ShopID
		cast.ShopName = shop.ShopName
	}

	return casts, nil
}

func (c *Client) getShopFromPage(ctx context.Context, strURL string) (*Cast, error) {
	resp, err := c.get(ctx, strURL, "pcmode=sp")
	if err != nil {
		return nil, fmt.Errorf(`on get("%s"): %w`, strURL, err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("on NewDocumentFromReader(): %w", err)
	}

	shop := &Cast{}

	doc.Find("a.shopinfobox-button").EachWithBreak(func(_ int, a *goquery.Selection) bool {
		if shopID, ok := a.Attr("data-c_commu_id"); ok {
			shop.ShopID, _ = strconv.Atoi(shopID)
			shop.ShopName, _ = a.Attr("data-infoname")

			return false
		}

		return true
	})

	return shop, nil
}
