package mc

import (
	"github.com/willroberts/minecraft-client"
)

type ClientOpts struct {
	Password string
	HostPort string
}

type Client struct {
	opts ClientOpts
}

var client *minecraft.Client

func getClient(opts ClientOpts) (*minecraft.Client, error) {
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

func (c *Client) Command(cmd string) error {
	mc, err := getClient(c.opts)
	if err != nil {
		return err
	}

	if _, err := mc.SendCommand(cmd); err != nil {
		return err
	}

	return nil
}

func NewClient(opts ClientOpts) (*Client, error) {
	return &Client{opts}, nil
}
