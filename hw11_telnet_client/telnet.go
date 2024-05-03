package main

import (
	"io"
	"log"
	"net"
	"time"
)

const networkType = "tcp"

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (t *telnetClient) Connect() error {
	conn, err := net.DialTimeout(networkType, t.address, t.timeout)
	if err != nil {
		return err
	}

	t.conn = conn
	log.Printf("Connected to %s", t.address)

	return nil
}

func (t *telnetClient) Close() error {
	if t.conn != nil {
		err := t.conn.Close()
		return err
	}
	return nil
}

func (t *telnetClient) Send() error {
	_, err := io.Copy(t.conn, t.in)
	return err
}

func (t *telnetClient) Receive() error {
	_, err := io.Copy(t.out, t.conn)
	return err
}
