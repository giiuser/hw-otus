package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type client struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &client{address: address, timeout: timeout, in: in, out: out}
}

func (c *client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return fmt.Errorf("error connection: %w", err)
	}
	c.conn = conn
	return nil
}

func (c *client) Close() error {
	return c.conn.Close()
}

func (c *client) Send() error {
	if _, err := io.Copy(c.conn, c.in); err != nil {
		return fmt.Errorf("error sending: %w", err)
	}
	fmt.Fprintln(os.Stderr, "...EOF")
	return nil
}

func (c *client) Receive() error {
	if _, err := io.Copy(c.out, c.conn); err != nil {
		return fmt.Errorf("error receiving: %w", err)
	}
	fmt.Fprintln(os.Stderr, "...Connection was closed by peer")
	return nil
}
