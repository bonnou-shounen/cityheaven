package util

import (
	"context"
	"errors"
	"fmt"
	"os"

	libnetrc "github.com/jdxcode/netrc"

	cityheaven "github.com/bonnou-shounen/cityheaven"
)

func NewLoggedClient(ctx context.Context) (*cityheaven.Client, error) {
	id, password := getCredential()
	if id == "" || password == "" {
		return nil, errors.New("missing credentials")
	}

	client := cityheaven.NewClient()

	err := client.Login(ctx, id, password)
	if err != nil {
		return nil, fmt.Errorf(`error on Login("%s", "***"): %w`, id, err)
	}

	return client, nil
}

func getCredential() (id, password string) {
	getters := []func() (string, string){
		fromEnv,
		fromNetrc,
	}

	for _, getter := range getters {
		if id != "" && password != "" {
			return
		}

		i, p := getter()

		if id == "" {
			id = i
		}

		if password == "" {
			password = p
		}
	}

	return
}

func fromEnv() (id, password string) {
	id = os.Getenv("CITYHEAVEN_LOGIN")
	password = os.Getenv("CITYHEAVEN_PASSWORD")

	return
}

func fromNetrc() (id, password string) {
	netrc := getNetrc()
	if netrc == nil {
		return
	}

	machine := netrc.Machine("www.cityheaven.net")
	if machine == nil {
		return
	}

	id = machine.Get("login")
	password = machine.Get("password")

	return
}

func getNetrc() *libnetrc.Netrc {
	netrcPath := getNetrcPath()
	if netrcPath == "" {
		return nil
	}

	netrc, err := libnetrc.Parse(netrcPath)
	if err != nil {
		return nil
	}

	return netrc
}

func getNetrcPath() string {
	path := os.Getenv("NETRC")
	if path != "" {
		return path
	}

	path = os.Getenv("CURLOPT_NETRC_FILE")

	if path != "" {
		return path
	}

	dir := os.Getenv("HOME")
	if dir != "" {
		return dir + "/.netrc"
	}

	return ""
}
