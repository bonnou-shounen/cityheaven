package cityheaven

import (
	"context"
	"fmt"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

func (c *Client) GetShopAttendees(ctx context.Context, strURL string) ([]*Cast, error) {
	resp, err := c.getRaw(ctx, fmt.Sprint(strURL, "attend/soon/"), "")
	if err != nil {
		return nil, fmt.Errorf("on getRaw(): %w", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("on NewDocumentFromReader(): %w", err)
	}

	var shopID int

	doc.Find("head script").EachWithBreak(func(_ int, scr *goquery.Selection) bool {
		shopID = c.parseNumber(scr.Text(), `{'shop_id':'`, `'}`)

		return shopID == 0
	})

	shopName := doc.Find(`div#location span[itemprop="name"]`).Last().Text()

	var casts []*Cast

	re := regexp.MustCompile(`\d{2}:\d{2}`) //nolint:varnamelen

	doc.Find("div.sugunavi_wrapper").Each(func(_ int, div *goquery.Selection) {
		href, _ := div.Find("a").Attr("href")
		castID := c.parseNumber(href, "girlid-", "/")

		castName := div.Find("p.name_font_size").Text()

		nextStart := string(re.Find(
			[]byte(div.Find("div.title").Text()),
		))

		if castID != 0 && castName != "" && nextStart != "" {
			casts = append(casts,
				&Cast{
					ID:        castID,
					Name:      castName,
					ShopID:    shopID,
					ShopName:  shopName,
					NextStart: nextStart,
				},
			)
		}
	})

	return casts, nil
}
