package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/bonnou-shounen/cityheaven"
)

type DumpShopCasts struct {
	NoFav     bool   `help:"skip counting favorites"`
	Attendees bool   `help:"only attendees"`
	Area      string `help:"area part in shop URL"`
	Shop      string `xor:"shop-url" help:"name part in shop URL"`
	URL       string `xor:"shop-url" help:"the shop URL"`
}

func (d *DumpShopCasts) Run() error {
	strURL, err := d.getURL()
	if err != nil {
		return err
	}

	ctx := context.Background()
	client := cityheaven.NewClient()

	var casts []*cityheaven.Cast

	if d.Attendees {
		casts, err = client.GetShopAttendees(ctx, strURL)
		if err != nil {
			return fmt.Errorf(`on GetShopAtendees("%s"): %w`, strURL, err)
		}
	} else {
		casts, err = client.GetShopCasts(ctx, strURL)
		if err != nil {
			return fmt.Errorf(`on GetShopCasts("%s"): %w`, strURL, err)
		}
	}

	for _, cast := range casts {
		line := fmt.Sprint(cast.ID, "\t", cast.ShopID)

		if !d.NoFav {
			favCount, _ := client.GetFavoriteCount(ctx, cast)
			line += fmt.Sprint("\t", favCount)
		}

		if d.Attendees {
			line += fmt.Sprint("\t", cast.NextStart)
		}

		line += fmt.Sprint("\t", cast.Name, "\t", cast.ShopName)

		fmt.Fprintln(os.Stdout, line)
	}

	return nil
}

func (d *DumpShopCasts) getURL() (string, error) {
	area := d.Area
	if area == "" {
		area = os.Getenv("CITYHEAVEN_AREA")
	}

	if area == "" {
		area = "tokyo"
	}

	shop := d.Shop
	if shop == "" {
		url := d.URL
		if url == "" {
			url = d.readURL()
		}

		if err := d.parseURL(url, &area, &shop); err != nil {
			return "", fmt.Errorf(`on parseURL("%s"): %w`, url, err)
		}
	}

	ctx := context.Background()
	c := cityheaven.NewClient()

	strURL, err := c.GetShopURL(ctx, area, shop)
	if err != nil {
		return "", fmt.Errorf(`on GetShopURL("%s", "%s"): %w`, area, shop, err)
	}

	return strURL, nil
}

func (d *DumpShopCasts) readURL() string {
	fmt.Fprint(os.Stderr, "paste shop URL: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	return scanner.Text()
}

func (d *DumpShopCasts) parseURL(str string, pArea, pShop *string) error {
	iAreaEnd := strings.Index(str, "/A")
	if iAreaEnd < 0 {
		return fmt.Errorf(`not found area end "/A"`)
	}

	iArea := strings.LastIndex(str[:iAreaEnd], "/")
	if iArea < 0 {
		return fmt.Errorf(`not found area start "/"`)
	}

	*pArea = str[iArea+1 : iAreaEnd]

	iShop := iAreaEnd + 15

	iShopLen := strings.Index(str[iShop:], "/")
	if iShopLen < 0 {
		return fmt.Errorf(`not found shop end "/"`)
	}

	*pShop = str[iShop : iShop+iShopLen]

	return nil
}
