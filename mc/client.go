package mc

import (
	"github.com/willroberts/minecraft-client"
)

type ClientOpts struct {
	Password string
	HostPort string
}

type Client struct {
	mc *minecraft.Client
}

func (c *Client) Command(cmd string) error {
	if _, err := c.mc.SendCommand(cmd); err != nil {
		return err
	}
	return nil
}

func NewClient(opts ClientOpts) (*Client, error) {
	mc, err := minecraft.NewClient(opts.HostPort)
	if err != nil {
		return nil, err
	}

	if err := mc.Authenticate(opts.Password); err != nil {
		return nil, err
	}

	return &Client{mc}, nil
}
