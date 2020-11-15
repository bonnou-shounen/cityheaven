package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/bonnou-shounen/cityheaven"
)

type DumpShopCasts struct {
	NoFav bool   `help:"skip counts"`
	Area  string `help:"area name in URL (default: tokyo)"`
	Shop  string `xor:"shop-url" help:"shop name in URL"`
	URL   string `xor:"shop-url" help:"URL of shop page"`
}

func (d *DumpShopCasts) Run() error {
	strURL, err := d.getURL()
	if err != nil {
		return err
	}

	c := cityheaven.NewClient()

	casts, err := c.GetShopCasts(strURL)
	if err != nil {
		return err
	}

	for _, cast := range casts {
		var favCount int
		if !d.NoFav {
			favCount, _ = c.GetFavoriteCount(cast)
		}

		fmt.Fprintf(os.Stdout, "%d\t%d\t%d\t%s\t%s\n", cast.ShopID, cast.ID, favCount, cast.Name, cast.ShopName)
	}

	return nil
}

func (d *DumpShopCasts) getURL() (string, error) {
	var area, shop string

	if d.Shop != "" {
		shop = d.Shop

		if d.Area != "" {
			area = d.Area
		} else {
			area = "tokyo"
		}
	} else {
		if d.URL == "" {
			d.URL = d.readURL()
		}
		err := d.parseURL(d.URL, &area, &shop)
		if err != nil {
			return "", err
		}
	}

	c := cityheaven.NewClient()

	strURL, err := c.GetShopURL(area, shop)
	if err != nil {
		return "", err
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
	ParseError := fmt.Errorf("parse URL failed: [%s]", str)

	iAreaEnd := strings.Index(str, "/A")
	if iAreaEnd < 0 {
		return ParseError
	}

	iArea := strings.LastIndex(str[:iAreaEnd], "/")
	if iArea < 0 {
		return ParseError
	}

	*pArea = str[iArea+1 : iAreaEnd]

	iShop := iAreaEnd + 15

	iShopLen := strings.Index(str[iShop:], "/")
	if iShopLen < 0 {
		return ParseError
	}

	*pShop = str[iShop : iShop+iShopLen]

	return nil
}
