package client

import (
	"fmt"
	"net"
)

type Client struct {
	Conn net.Conn
}

func New(serverAddr string) (*Client, error) {
	Conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return nil, err
	}

	return &Client{Conn: Conn}, nil
}

func (c *Client) Start() *Client {
	fmt.Fprintf(c.Conn, "start\n")
	return c
}

func (c *Client) Stop() *Client {
	fmt.Fprintf(c.Conn, "stop\n")
	return c
}

func (c *Client) Refresh() *Client {
	fmt.Fprintf(c.Conn, "refresh\n")
	return c
}

func (c *Client) Message(s string) *Client {
	fmt.Fprintf(c.Conn, "%s\n", s)
	return c
}

func (c *Client) Close() {
	c.Conn.Close()
}
