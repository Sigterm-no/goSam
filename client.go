package goSam

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

// A Client represents a single Connection to the SAM bridge
type Client struct {
	SamConn net.Conn
	verbose bool
}

// NewDefaultClient creates a new client, connecting to the default host:port at localhost:7656
func NewDefaultClient() (*Client, error) {
	return NewClient("localhost:7656")
}

// NewClient creates a new client, connecting to a specified port
func NewClient(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	c := &Client{conn, false}
	return c, c.hello()
}

// ToggleVerbose switches logging on or off.
// (also passed to new clients inside Dial.)
func (c *Client) ToggleVerbose() {
	c.verbose = !c.verbose
}

// send the initial handshake command and check that the reply is ok
func (c *Client) hello() (err error) {
	var r *Reply

	r, err = c.sendCmd("HELLO VERSION MIN=3.0 MAX=3.0")
	if err != nil {
		return err
	}

	if r.Topic != "HELLO" {
		return fmt.Errorf("Unknown Reply: %+v\n", r)
	}

	if r.Pairs["RESULT"] != "OK" || r.Pairs["VERSION"] != "3.0" {
		return fmt.Errorf("Handshake did not succeed\nReply:%+v\n", r)
	}

	return nil
}

// helper to send one command and parse the reply by sam
func (c *Client) sendCmd(cmd string) (r *Reply, err error) {
	if _, err = fmt.Fprintln(c.SamConn, cmd); err != nil {
		return
	}

	if c.verbose {
		log.Printf(">Send>'%s'\n", cmd)
	}

	reader := bufio.NewReader(c.SamConn)
	line, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	if c.verbose {
		log.Printf("<Rcvd<'%s'\n", line)
	}

	r, err = parseReply(line)
	return
}

// Close the underlying socket to SAM
func (c *Client) Close() error {
	return c.SamConn.Close()
}
