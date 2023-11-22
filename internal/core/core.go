package core

import "github.com/1F47E/go-shaihulud/internal/client"

type Core struct {
	Client *client.Client
}

func NewCore(cli *client.Client) *Core {
	return &Core{cli}
}
