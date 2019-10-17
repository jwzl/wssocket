package client

import (
	"io"
	"net"
	"sync"
	"time"
	"net/http"
	"crypto/tls"
	
	"github.com/jwzl/wssocket/conn"
	"github.com/jwzl/wssocket/model"
)


type Options struct {
	// client type
	ConnUse string
	// tls config
	TLSConfig *tls.Config
	// HandshakeTimeout is the maximum duration that the cryptographic handshake may take.
	HandshakeTimeout time.Duration
	// auto route flag
	AutoRoute bool
	//Message revice handler. 
	Handler          conn.Handler
	// this is for stream message
	Consumer         io.Writer		////optional.
	// Connected callback
	Connected	func(*conn.Connection, *http.Response)
}

//Client
type Client struct {
	Options
	//http reuqest
	RequestHeader http.Header
	//Connection	
	conn	*conn.Connection
	// WSClient
	wsc *WSClient
	//Connection look
	connLock sync.Mutex
}

//Start the ws client
func (c *Client) Start() {
	c.wsc = NewWSClient(c.Options)
}

//Connect the remote server
func (c *Client) Connect(serverAddr string) error {
	c.connLock.Lock()
	conn, err := c.wsc.Connect(serverAddr, c.RequestHeader)
	if err != nil {
		c.connLock.Unlock()
		return err
	}

	c.conn = conn
	c.connLock.Unlock()

	go c.conn.ConnRecieve()
	return nil
}

func CreateTLSConfig(certFile, keyFile string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	return tlsConfig, nil
}
// WriteMessage
func (c *Client) WriteMessage(msg *model.Message) error {
	return c.conn.WriteMessage(msg)
}
// SetReadDeadline
func (c *Client) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}
// SetWriteDeadline
func (c *Client) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}
// RemoteAddr
func (c *Client) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}
// RemoteAddr
func (c *Client) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

//Close connection
func (c *Client) Close() error {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	if c.conn != nil {
		return c.conn.Close()
	}
	
	return nil
}
