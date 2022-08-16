package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"time"

	"giautm.dev/hanetai"
	"golang.org/x/oauth2"
)

type CliContext struct {
	AccessToken string
	Context     context.Context
	Debug       bool
	JSON        bool
	NoHeader    bool
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

func (c *CliContext) Writer() io.Writer {
	return os.Stdout
}
