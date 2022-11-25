package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/bonnou-shounen/cityheaven"
)

func (o *optShop) GetURL() (string, error) {
	area := o.Area

	shop := o.Shop
	if shop == "" {
		url := o.URL
		if url == "" {
			url = o.readURL()
		}

		if err := o.parseURL(url, &area, &shop); err != nil {
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

func (o *optShop) readURL() string {
	fmt.Fprint(os.Stderr, "paste shop URL: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	return scanner.Text()
}

func (o *optShop) parseURL(str string, pArea, pShop *string) error {
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
