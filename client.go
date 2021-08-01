package cityheaven

import (
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
	jar, _ := cookiejar.New(nil)

	return &Client{
		http: &http.Client{Jar: jar},
	}
}

func (c *Client) Login(id, password string) error {
	values := url.Values{
		"user": []string{id},
		"pass": []string{password},
	}

	u, _ := url.Parse("https://www.cityheaven.net/")
	c.http.Jar.SetCookies(u, []*http.Cookie{{}})

	resp, err := c.http.PostForm("https://www.cityheaven.net/tokyo/loginajax/", values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	res := struct{ IsLogin bool }{}

	err = decoder.Decode(&res)
	if err != nil || !res.IsLogin {
		return fmt.Errorf("login failed")
	}

	return nil
}

func (c *Client) get(strURL string, values url.Values) error {
	resp, err := c.getRaw(strURL, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (c *Client) getRaw(strURL string, values url.Values) (*http.Response, error) {
	return c.http.Get(fmt.Sprint(strURL, "?", values.Encode()))
}

func (c *Client) post(strURL, body string) error {
	resp, err := c.http.Post(strURL, "application/x-www-form-urlencoded", strings.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
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
