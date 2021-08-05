package cityheaven

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

func (c *Client) GetFavoriteCasts(ctx context.Context) ([]*Cast, error) {
	resp, err := c.getRaw(ctx, "https://www.cityheaven.net/tt/community/ABEditFavoriteGirl/", "")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error on NewDocumentFromReader(): %w", err)
	}

	var casts []*Cast

	doc.Find("div.myshopnew").Each(func(j int, div *goquery.Selection) {
		shopName := div.Find("p.myshopnew-head-name").Text()
		castName := div.Find("p.myshopnew-title").Text()
		strCastID, _ := div.Find(`input[name="girl"]`).Attr("value")
		strShopID, _ := div.Find(`input[name="commu"]`).Attr("value")

		castID, _ := strconv.Atoi(strCastID)
		shopID, _ := strconv.Atoi(strShopID)

		if castID == 0 || castName == "" || shopID == 0 || shopName == "" {
			return
		}

		casts = append(casts,
			&Cast{
				ID:       castID,
				Name:     castName,
				ShopID:   shopID,
				ShopName: shopName,
			},
		)
	})

	return casts, nil
}

func (c *Client) AddFavoriteCast(ctx context.Context, cast *Cast) error {
	values := url.Values{
		"girlId": []string{fmt.Sprint(cast.ID)},
	}

	return c.get(ctx, "https://www.cityheaven.net/tokyo/A0000/A000000/a/okiniiri/", values)
}

func (c *Client) DeleteFavoriteCast(ctx context.Context, cast *Cast) error {
	return c.DeleteFavoriteCasts(ctx, []*Cast{cast})
}

func (c *Client) AddFavoriteCasts(ctx context.Context, casts []*Cast) error {
	var anyErr error

	for i := len(casts) - 1; i >= 0; i-- {
		err := c.AddFavoriteCast(ctx, casts[i])
		if err != nil {
			anyErr = err
		}
	}

	return anyErr
}

func (c *Client) DeleteFavoriteCasts(ctx context.Context, casts []*Cast) error {
	if len(casts) == 0 {
		return nil
	}

	values := url.Values{}

	for _, cast := range casts {
		values.Add(fmt.Sprint("data_", cast.ID), "削除する")
	}

	return c.get(ctx, "https://www.cityheaven.net/tt/community/ABEditFavoriteGirl/", values)
}

func (c *Client) SortFavoriteCasts(ctx context.Context, casts []*Cast) error {
	if len(casts) == 0 {
		return nil
	}

	queryB := bytes.NewBufferString("update=変更を反映する")

	for _, cast := range casts {
		queryB.WriteString(fmt.Sprintf("&sort_girl[%d]=1", cast.ID))
	}

	resp, err := c.postRaw(ctx, "https://www.cityheaven.net/y/community/ABEditFavoriteGirl/", queryB.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (c *Client) GetFavoriteCount(ctx context.Context, cast *Cast) (int, error) {
	values := url.Values{
		"girl_id":  []string{fmt.Sprint(cast.ID)},
		"commu_id": []string{fmt.Sprint(cast.ShopID)},
	}

	resp, err := c.getRaw(ctx, "https://www.cityheaven.net/api/myheaven/v1/getgirlfavcnt/", values.Encode())
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	res := struct{ Cnt string }{}

	err = decoder.Decode(&res)
	if err != nil {
		return 0, fmt.Errorf("error on Decode(): %w", err)
	}

	count, err := strconv.Atoi(res.Cnt)
	if err != nil {
		return 0, fmt.Errorf(`error on Atoi("%s"): %w`, res.Cnt, err)
	}

	return count, nil
}
