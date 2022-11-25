package util

import (
	"context"
	"errors"
	"fmt"

	"github.com/bonnou-shounen/cityheaven"
)

func NewLoggedClient(ctx context.Context, loginID, password string) (*cityheaven.Client, error) {
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
