package mc

import (
	"context"

	"github.com/willroberts/minecraft-client"
)

type ClientConfig struct {
	Password string
	HostPort string
}

type Client struct {
	cfg ClientConfig
}

var client *minecraft.Client

func getClient(opts ClientConfig) (*minecraft.Client, error) {
	if client != nil {
		return client, nil
	}

	var err error
	client, err := minecraft.NewClient(opts.HostPort)
	if err != nil {
		return nil, err
	}

	if err := client.Authenticate(opts.Password); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) Command(ctx context.Context, cmd string) (string, error) {
	mc, err := getClient(c.cfg)
	if err != nil {
		return "", err
	}

	response, err := mc.SendCommand(cmd)
	if err != nil {
		return "", err
	}

	return response.Body, err
}

func NewClient(cfg ClientConfig) *Client {
	return &Client{cfg}
}
