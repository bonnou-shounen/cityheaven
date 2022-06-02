package util

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/bonnou-shounen/cityheaven"
	libnetrc "github.com/jdxcode/netrc"
)

func NewLoggedClient(ctx context.Context) (*cityheaven.Client, error) {
	loginID, password := getCredential()
	if loginID == "" || password == "" {
		return nil, errors.New("missing credentials")
	}

	client := cityheaven.NewClient()

	err := client.Login(ctx, loginID, password)
	if err != nil {
		return nil, fmt.Errorf(`on Login("%s", "***"): %w`, loginID, err)
	}

	return client, nil
}

func getCredential() (string, string) {
	var loginID, password string

	getters := []func() (string, string){
		fromEnv,
		fromNetrc,
	}
	for _, getter := range getters {
		if loginID != "" && password != "" {
			return "", ""
		}

		id, pwd := getter()

		if loginID == "" {
			loginID = id
		}

		if password == "" {
			password = pwd
		}
	}

	return loginID, password
}

func fromEnv() (string, string) {
	loginID := os.Getenv("CITYHEAVEN_LOGIN")
	password := os.Getenv("CITYHEAVEN_PASSWORD")

	return loginID, password
}

func fromNetrc() (string, string) {
	netrc := getNetrc()
	if netrc == nil {
		return "", ""
	}

	machine := netrc.Machine("www.cityheaven.net")
	if machine == nil {
		return "", ""
	}

	loginID := machine.Get("login")
	password := machine.Get("password")

	return loginID, password
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
	for _, name := range []string{"NETRC", "CURLOPT_NETRC_FILE"} {
		if path := os.Getenv(name); path != "" {
			return path
		}
	}

	if dir := os.Getenv("HOME"); dir != "" {
		return dir + "/.netrc"
	}

	return ""
}
