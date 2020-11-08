package cityheaven

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

func (c *Client) GetFavoriteCasts() ([]*Cast, error) {
	resp, err := c.http.Get("https://www.cityheaven.net/tt/community/ABFavoriteGirlList/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var casts []*Cast

	doc.Find("div.myshopnew").Each(func(j int, div *goquery.Selection) {
		shop := div.Find("p.myshopnew-head-name").Text()
		name := div.Find("p.myshopnew-title").Text()
		girl, _ := div.Find(`input[name="girl"]`).Attr("value")
		commu, _ := div.Find(`input[name="commu"]`).Attr("value")

		castID, _ := strconv.Atoi(girl)
		shopID, _ := strconv.Atoi(commu)

		if name != "" && shopID != 0 && castID != 0 {
			casts = append(casts,
				&Cast{
					ShopID:   shopID,
					CastID:   castID,
					Name:     name,
					ShopName: shop,
				},
			)
		}
	})

	return casts, nil
}

func (c *Client) AddFavoriteCast(cast *Cast) error {
	values := url.Values{
		"girlId": []string{fmt.Sprint(cast.CastID)},
	}

	return c.get("https://www.cityheaven.net/tokyo/A0000/A000000/a/okiniiri/", values)
}

func (c *Client) DeleteFavoriteCast(cast *Cast) error {
	return c.DeleteFavoriteCasts([]*Cast{cast})
}

func (c *Client) AddFavoriteCasts(casts []*Cast) error {
	var anyErr error

	for i := len(casts) - 1; i >= 0; i-- {
		err := c.AddFavoriteCast(casts[i])
		if err != nil {
			anyErr = err
		}
	}

	return anyErr
}

func (c *Client) DeleteFavoriteCasts(casts []*Cast) error {
	if len(casts) == 0 {
		return nil
	}

	values := url.Values{}

	for _, cast := range casts {
		values.Add(fmt.Sprint("data_", cast.CastID), "削除する")
	}

	return c.get("https://www.cityheaven.net/tt/community/ABEditFavoriteGirl/", values)
}

func (c *Client) SortFavoriteCasts(casts []*Cast) error {
	if len(casts) == 0 {
		return nil
	}

	queryB := bytes.NewBufferString("update=変更を反映する")

	for _, cast := range casts {
		queryB.WriteString(fmt.Sprintf("&sort_girl[%d]=1", cast.CastID))
	}

	return c.post("https://www.cityheaven.net/y/community/ABEditFavoriteGirl/", queryB.String())
}

func (c *Client) GetFavoriteCount(cast *Cast) (int, error) {
	values := url.Values{
		"commu_id": []string{fmt.Sprint(cast.ShopID)},
		"girl_id":  []string{fmt.Sprint(cast.CastID)},
	}
	resp, err := c.getRaw("https://www.cityheaven.net/api/myheaven/v1/getgirlfavcnt/", values)

	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	res := struct{ Cnt string }{}

	err = decoder.Decode(&res)
	if err != nil {
		return 0, err
	}

	count, err := strconv.Atoi(res.Cnt)
	if err != nil {
		return 0, err
	}

	return count, nil
}
