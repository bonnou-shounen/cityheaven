package cityheaven

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/jesse0michael/errgroup"
)

func (c *Client) GetShopURL(ctx context.Context, area, shop string) (string, error) {
	req, err := http.NewRequest(
		http.MethodHead,
		fmt.Sprint("https://www.cityheaven.net/", area, "/A0000/A000000/", shop, "/"),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("on NewRequest(): %w", err)
	}

	resp, _ := http.DefaultTransport.RoundTrip(req.WithContext(ctx))
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusMovedPermanently {
		locations := resp.Header["Location"]
		if len(locations) > 0 {
			return "https:" + locations[0], nil
		}
	}

	return "", fmt.Errorf("shop not found: %s of %s", shop, area)
}

type castsPageInfo struct {
	LastPage int
	ShopID   int
	ShopName string
}

func (c *Client) GetShopCasts(ctx context.Context, strURL string) ([]*Cast, error) {
	var info castsPageInfo

	casts, err := c.getShopCastsOnPage(ctx, strURL, 1, &info)
	if err != nil {
		return nil, err
	}

	if info.LastPage >= 2 {
		seg, segCtx := errgroup.WithContext(ctx, 3)

		castsOnPage := make([][]*Cast, info.LastPage+1)

		for page := 2; page <= info.LastPage; page++ {
			page := page

			seg.Go(func() error {
				casts, err := c.getShopCastsOnPage(segCtx, strURL, page, nil)
				if err != nil {
					return fmt.Errorf("on getShopCastsOnPage(%d): %w", page, err)
				}

				castsOnPage[page] = casts

				return nil
			})
		}

		if err := seg.Wait(); err != nil {
			return nil, fmt.Errorf("on a goroutine: %w", err)
		}

		for page := 2; page <= info.LastPage; page++ {
			casts = append(casts, castsOnPage[page]...)
		}
	}

	for _, cast := range casts {
		cast.ShopID = info.ShopID
		cast.ShopName = info.ShopName
	}

	return casts, nil
}

func (c *Client) getShopCastsOnPage(ctx context.Context, strURL string, page int, pInfo *castsPageInfo) ([]*Cast, error) { //nolint:lll
	resp, err := c.getRaw(ctx, fmt.Sprint(strURL, "girllist/", page, "/"), "")
	if err != nil {
		return nil, fmt.Errorf("on getRaw(): %w", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("on NewDocumentFromReader(): %w", err)
	}

	div := doc.Find("div.girllistimg")
	if div.Length() == 0 {
		return c.getShopCastsOnOldPage(doc, pInfo)
	}

	var casts []*Cast

	div.Each(func(_ int, div *goquery.Selection) {
		castName, _ := div.Find("img.no_login").Attr("alt")
		href, _ := div.Find("a").Attr("href")
		castID := c.parseNumber(href, "girlid-", "/")

		if castID != 0 && castName != "" {
			casts = append(casts,
				&Cast{
					ID:   castID,
					Name: castName,
				},
			)
		}
	})

	if pInfo != nil {
		pInfo.LastPage, _ = strconv.Atoi(doc.Find(`ul.paging a:not([class="next"])`).Last().Text())

		doc.Find("head script").EachWithBreak(func(_ int, scr *goquery.Selection) bool {
			shopID := c.parseNumber(scr.Text(), `{'shop_id':'`, `'}`)
			if shopID > 0 {
				pInfo.ShopID = shopID

				return false
			}

			return true
		})

		pInfo.ShopName = doc.Find(`div#location span[itemprop="name"]`).Last().Text()
	}

	return casts, nil
}

func (c *Client) getShopCastsOnOldPage(doc *goquery.Document, pInfo *castsPageInfo) ([]*Cast, error) {
	var casts []*Cast

	doc.Find("a").Each(func(_ int, a *goquery.Selection) {
		castName, _ := a.Attr("title")
		href, _ := a.Attr("href")
		castID := c.parseNumber(href, "girlid-", "/")

		if castID != 0 && castName != "" {
			casts = append(casts,
				&Cast{
					ID:   castID,
					Name: castName,
				},
			)
		}
	})

	if pInfo != nil {
		pInfo.LastPage, _ = strconv.Atoi(doc.Find(`ul.paging a:not([class="next"])`).Last().Text())

		doc.Find("head script").EachWithBreak(func(_ int, scr *goquery.Selection) bool {
			shopID := c.parseNumber(scr.Text(), `{'shop_id':'`, `'}`)
			if shopID > 0 {
				pInfo.ShopID = shopID

				return false
			}

			return true
		})

		pInfo.ShopName = doc.Find(`div#location span[itemprop="name"]`).Last().Text()
	}

	return casts, nil
}
