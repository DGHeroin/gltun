package gltun

import (
	"net"
	"errors"
)

var (
	errNotConnect = errors.New("client not connect")
)

type Client interface {
	Connect(Addr string) (Session, error)

	HandleDataFunc(func(Session, []byte))
	HandleClosedFunc(func(Session))

	Wait() error
	Close()
}

type client struct {
	conn net.Conn
	session Session

	dataFunc func(Session, []byte)
	closedFunc func(Session)
}

func NewClient() (Client) {
	return &client {session:nil}
}

func (c *client) Connect(Addr string) (Session, error) {
	if conn, err := net.Dial("tcp", Addr); err != nil {
		return nil, err
	} else {
		c.conn = conn
		c.session = NewSession(conn)
	}
	return c.session, nil
}

func (c *client) HandleDataFunc(cb func(Session, []byte)) { c.dataFunc = cb }
func (c *client) HandleClosedFunc(cb func(Session)) { c.closedFunc = cb }

func (c *client) Wait() (error) {
	buf := make([]byte, 4096)
	defer func() {
		if c.closedFunc != nil {
			c.closedFunc(c.session)
		}
		c.Close()
	}()

	for {
		if n, err := c.conn.Read(buf); err != nil {
			return err
		} else {
			data := buf[:n]
			if pkg := c.session.read(data); pkg != nil {
				if c.dataFunc != nil {
					c.dataFunc(c.session, pkg.GetPayload())
				}
			}
		}
	}
}


func (c *client) Close() {
	if c.session != nil {
		c.session.Close()
		c.session = nil
	}
}
