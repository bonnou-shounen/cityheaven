package cityheaven

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/remeh/sizedwaitgroup"
)

func (c *Client) GetShopURL(area, shop string) (string, error) {
	req, err := http.NewRequest(
		"HEAD",
		fmt.Sprint("https://www.cityheaven.net/", area, "/A0000/A000000/", shop, "/"),
		nil,
	)
	if err != nil {
		return "", err
	}

	resp, _ := http.DefaultTransport.RoundTrip(req)
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

func (c *Client) GetShopCasts(strURL string) ([]*Cast, error) {
	var info castsPageInfo

	casts, err := c.getShopCastsOnPage(strURL, 1, &info)
	if err != nil {
		return nil, err
	}

	if info.LastPage <= 1 {
		return casts, nil
	}

	castsOnPage := make([][]*Cast, info.LastPage+1)
	swg := sizedwaitgroup.New(3)

	for page := 2; page <= info.LastPage; page++ {
		swg.Add()

		go func(page int) {
			defer swg.Done()

			castsOnPage[page], _ = c.getShopCastsOnPage(strURL, page, nil)
		}(page)
	}
	swg.Wait()

	for page := 2; page <= info.LastPage; page++ {
		casts = append(casts, castsOnPage[page]...)
	}

	for _, cast := range casts {
		cast.ShopID = info.ShopID
		cast.ShopName = info.ShopName
	}

	return casts, nil
}

func (c *Client) getShopCastsOnPage(strURL string, page int, pInfo *castsPageInfo) ([]*Cast, error) {
	resp, err := c.http.Get(fmt.Sprint(strURL, "girllist/", page, "/"))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var casts []*Cast

	doc.Find("div.girllistimg").Each(func(_ int, div *goquery.Selection) {
		name, _ := div.Find("img").Attr("alt")
		href, _ := div.Find("a").Attr("href")
		castID := c.parseNumber(href, "girlid-", "/")

		if name != "" && castID != 0 {
			casts = append(casts,
				&Cast{
					CastID: castID,
					Name:   name,
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
