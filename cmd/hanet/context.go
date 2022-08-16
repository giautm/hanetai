package main

import (
	"context"
	"net/http"
	"time"

	"giautm.dev/hanetai"
	"golang.org/x/oauth2"
)

type CliContext struct {
	AccessToken string
	Context     context.Context
	Debug       bool
}

func (c *CliContext) NewClient() *hanetai.Client {
	source := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: c.AccessToken,
		TokenType:   "Bearer",
	})
	return hanetai.NewClient(&http.Client{
		Timeout: 60 * time.Second,
	}, source)
}
