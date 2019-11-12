package conn

import(
	"io"	
	"net"
	"time"
	"errors"
	"net/http"
	"k8s.io/klog"
	"crypto/x509"
	"encoding/json"
	"github.com/gorilla/websocket"
	wstype "github.com/jwzl/wssocket/types"
	"github.com/jwzl/wssocket/fifo"
	"github.com/jwzl/wssocket/model"
	"github.com/jwzl/wssocket/packer"	
	"github.com/jwzl/wssocket/translator"
)

type Handler interface {
	MessageProcess(Header http.Header, msg *model.Message, c *Connection)
}
// connection states
// TODO: add connection state filed
type ConnectionState struct {
	State            string
	Header          http.Header
	PeerCertificates []*x509.Certificate
}

type Connection struct {
	//Message revice handler. 
	Handler  Handler
	// auto route flag
	AutoRoute bool
	// client type
	ConnUse string
	// Consumer
	Consumer io.Writer
	// State
	State  *ConnectionState
	//websocket Connection
	Conn *websocket.Conn
	// message fifo.
	MessageFifo  *fifo.MessageFifo	
}

//Deep Copy http header.
func DeepCopyHeader(header http.Header) http.Header {
	headerByte, err := json.Marshal(header)
	if err != nil {
		klog.Errorf("faile to marshal header, error:%+v", err)
		return nil
	}

	dstHeader := make(http.Header)
	err = json.Unmarshal(headerByte, &dstHeader)
	if err != nil {
		klog.Errorf("failed to unmarshal header, error:%+v", err)
		return nil
	}
	return dstHeader
}

// start to recieve message from Connection
func (c *Connection) ConnRecieve(){
	switch c.ConnUse {
	case wstype.UseTypeMessage:
		go c.handleMessage()
	case wstype.UseTypeStream:
		go c.handleRawData()
	case wstype.UseTypeShare:	
		klog.Warningf("don't support share in websocket")	
	}
}

func (c *Connection) handleMessage(){
	msg := &model.Message{}
	for {
		// Read the message
		err := c.unpackPackageAndDecode(msg)		
		if err != nil {
			if err != io.EOF {
				klog.Errorf("failed to read message, error: %+v", err)
			}
			c.State.State = wstype.StatDisconnected
			c.Conn.Close()
			return 
		}

		// filter control message
		if filtered := c.filterControlMessage (msg); filtered {
			continue
		}

		// to check whether the message is a response or not

		// Put the messages into fifo and wait for reading
		if !c.AutoRoute {
			if c.MessageFifo != nil {
				c.MessageFifo.Write(msg)
				continue
			}
		}

		//let c handler to process message.
		if c.Handler != nil && c.Handler.MessageProcess != nil {
			c.Handler.MessageProcess(c.State.Header, msg, c)
		}
	}
}

// unpack the package from websocket connection and Decode into model message. 
func (c *Connection) unpackPackageAndDecode(msg *model.Message) error {
	rawData, err := packer.NewReader(c).Read()
	if err != nil {
		klog.Errorf("failed to read, error: %+v", err)
		return err
	}

	// convert raw data to protocol buf message, then into model message.
	return translator.NewTransCoding().Decode(rawData, msg)
}

// let model message convert to protocol buf message, then package this msg. 
func (c *Connection) encodeAndPackPackage(msg *model.Message) error {
	rawData, err := translator.NewTransCoding().Encode(msg)
	if err != nil {
		klog.Errorf("failed to Encode, error: %+v", err)
		return err
	}

	// pack the message and send by websocket.
	_, err = packer.NewWriter(c).Write(rawData)
	return err
}

// Read data from websocket connection. can MATCH io.Reader 
func (c *Connection) Read(p []byte) (int, error){
	_, msgData, err := c.Conn.ReadMessage()
	if err != nil {
		if err != io.EOF {
			klog.Errorf("failed to read data, error: %+v", err)
		}
		return len(msgData), err
	}

	p = append(p[:0], msgData...)
	return len(msgData), err
}

// write data into websocket connection. can MATCH io.Writer 
func (c *Connection) Write(p []byte) (int, error) {
	err := c.Conn.WriteMessage(websocket.BinaryMessage, p)
	if err != nil {
		klog.Errorf("write websocket message error: %+v", err)
		return len(p), err
	}

	return len(p), err
}

func (c *Connection) filterControlMessage (msg *model.Message) bool {
	//check control message
	//process control message
	// feedback the response
	return false
}

//Stream data from socket (raw data)
func (c *Connection) handleRawData(){
	if c.Consumer == nil {
		klog.Errorf("bad consumer for raw data!")
		return 
	}

	if !c.AutoRoute {
		return 
	}

	//Read the raw data
	_, err := io.Copy(c.Consumer, c)
	if err != nil {
		klog.Errorf("failed to copy data, error:", err)
		c.State.State = wstype.StatDisconnected
		c.Conn.Close()
		return
	}
}

//some API for user

// Connection 's WriteMessage
func (c *Connection) WriteMessage(msg *model.Message) error {
	return c.encodeAndPackPackage(msg)
}

// ReadMessage read the message from fifo.
func (c *Connection) ReadMessage() (*model.Message, error){
	if c.MessageFifo != nil {
		return c.MessageFifo.Read()
	}

	return nil, errors.New("message fifo is nil.")	
}
// Set ReadDeadline 
func (c *Connection) SetReadDeadline(t time.Time) error {
	return c.Conn.SetReadDeadline(t)
}

// Set WriteDeadline 
func (c *Connection) SetWriteDeadline(t time.Time) error {
	return c.Conn.SetWriteDeadline(t)
}

// Get RemoteAddr 
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// Get LocalAddr 
func (c *Connection) LocalAddr() net.Addr {
	return c.Conn.LocalAddr()
}

// Close connection 
func (c *Connection) Close() error {
	return c.Conn.Close()
}

// get Connection state
func (c *Connection) ConnectionState() *ConnectionState {
	return c.State
}
