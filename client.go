package cityheaven

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
)

type Client struct {
	http *http.Client
}

func NewClient() *Client {
	u, _ := url.Parse("https://www.cityheaven.net/")
	jar, _ := cookiejar.New(nil)
	jar.SetCookies(u, []*http.Cookie{{}})

	return &Client{
		http: &http.Client{Jar: jar},
	}
}

func (c *Client) Login(ctx context.Context, id, password string) error {
	if err := c.getSimple(ctx, "https://www.cityheaven.net/", nil); err != nil {
		return err
	}

	values := url.Values{
		"user": []string{id},
		"pass": []string{password},
	}

	resp, err := c.post(ctx, "https://www.cityheaven.net/tokyo/loginajax/", values.Encode())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	res := struct{ IsLogin bool }{}

	err = decoder.Decode(&res)
	if err != nil {
		return fmt.Errorf("on Decode(): %w", err)
	}

	if !res.IsLogin {
		return fmt.Errorf("login failed: %s", id)
	}

	return nil
}

func (c *Client) getSimple(ctx context.Context, strURL string, values url.Values) error {
	resp, err := c.get(ctx, strURL, values.Encode())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (c *Client) get(ctx context.Context, strURL string, query string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, strURL+"?"+query, nil)
	if err != nil {
		return nil, fmt.Errorf("on NewRequest(): %w", err)
	}

	resp, err := c.http.Do(req.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("on http.Do(): %w", err)
	}

	return resp, nil
}

func (c *Client) post(ctx context.Context, strURL string, form string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, strURL, strings.NewReader(form))
	if err != nil {
		return nil, fmt.Errorf("on NewRequest(): %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.http.Do(req.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("on http.Do(): %w", err)
	}

	return resp, nil
}

func (c *Client) parseNumber(str, prefix, suffix string) int {
	if i := strings.Index(str, prefix); i >= 0 {
		str = str[i+len(prefix):]
		if j := strings.Index(str, suffix); j >= 0 {
			str = str[:j]
		}
	}

	num, _ := strconv.Atoi(str)

	return num
}
