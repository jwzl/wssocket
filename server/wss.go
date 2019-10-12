package server

import (
	"log"
	"io"	
	"os"
	"fmt"
	"net"
	"time"
	"net/http"

	"github.com/gorilla/websocket"
	wstype "github.com/jwzl/wssocket/types"	
	"github.com/jwzl/wssocket/model"	
	"github.com/jwzl/wssocket/packer"	
	"github.com/jwzl/wssocket/translator"		
)

type WSServer struct {
	WriteDeadline time.Time
	ReadDeadline  time.Time
	//websocket use 
	ConnUse	string
	options Server
	server  *http.Server
	conn	*websocket.Conn
}

func NewWSServer(opts Server) *WSServer {

	server := http.Server{
		Addr: 		opts.Addr,
		TLSConfig: 	opts.TLSConfig,
		ErrorLog:	log.New(os.Stderr, "", log.LstdFlags),
	}

	wsServer := &WSServer{
		options:	opts,
		server:		&server,
	}

	//register a http route handle.
	http.HandleFunc("/", wsServer.ServerHTTP)

	return wsServer
}


// Convert http server to websocket server
func (wss *WSServer) upgrade(w http.ResponseWriter, r *http.Request) *websocket.Conn {
	upgrader := websocket.Upgrader{
		HandshakeTimeout: wss.options.HandshakeTimeout,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("failed to upgrade to websocket")
		return nil
	}

	return conn
}

func (wss *WSServer) ServerHTTP(w http.ResponseWriter, r *http.Request){
	if wss.options.Filter != nil {
		if filtered := wss.options.Filter(w, r); filtered {	
			return
		}
	} 

	WsConn := wss.upgrade(w, r)
	if WsConn == nil {
		return	
	}
	wss.conn = WsConn
	wss.ConnUse = string(r.Header.Get("ConnectionUse"))

	// Connection callback.
	if wss.options.ConnNotify != nil {
		wss.options.ConnNotify(wss.options)
	}

	//Handle connection
	go wss.handleConn()
}

func (wss *WSServer) handleConn(){
	switch 	wss.ConnUse {
	case wstype.UseTypeMessage:
		go wss.handleMessage()
	case wstype.UseTypeStream:
		go wss.handleRawData()
	case wstype.UseTypeShare:	
		log.Println("don't support share in websocket")
	}
}

// unpack the package from websocket connection and Decode into model message. 
func (wss *WSServer) unpackPackageAndDecode(msg *model.Message) error {
	rawData, err := packer.NewReader(wss).Read()
	if err != nil {
		fmt.Errorf("failed to read, error: %+v", err)
		return err
	}

	// convert raw data to protocol buf message, then into model message.
	return translator.NewTransCoding().Decode(rawData, msg)
}

// let model message convert to protocol buf message, then package this msg. 
func (wss *WSServer) encodeAndPackPackage(msg *model.Message) error {
	rawData, err := translator.NewTransCoding().Encode(msg)
	if err != nil {
		fmt.Errorf("failed to Encode, error: %+v", err)
		return err
	}

	// pack the message and send by websocket.
	_, err = packer.NewWriter(wss).Write(rawData)
	return err
}

func (wss *WSServer) filterControlMessage (msg *model.Message) bool {
	//check control message
	//process control message
	// feedback the response
	return false
}

func (wss *WSServer) handleMessage(){
	msg := &model.Message{}
	for {
		// Read the message
		err := wss.unpackPackageAndDecode(msg)		
		if err != nil {
			if err != io.EOF {
				fmt.Errorf("failed to read message, error: %+v", err)
			}
			wss.conn.Close()
			return 
		}

		// filter control message
		if filtered := wss.filterControlMessage (msg); filtered {
			continue
		}

		// to check whether the message is a response or not

		// put the messages into fifo and wait for reading

		//let wss handler to process message.
		if wss.options.Handler != nil && wss.options.Handler.MessageProcess != nil {
			wss.options.Handler.MessageProcess(nil, msg)
		}
	}
}

// Read data from websocket connection. can MATCH io.Reader 
func (wss *WSServer) Read(p []byte) (int, error){
	_, msgData, err := wss.conn.ReadMessage()
	if err != nil {
		if err != io.EOF {
			fmt.Errorf("failed to read data, error: %+v", err)
		}
		return len(msgData), err
	}

	p = append(p[:0], msgData...)
	return len(msgData), err
}

// write data into websocket connection. can MATCH io.Writer 
func (wss *WSServer) Write(p []byte) (int, error) {
	err := wss.conn.WriteMessage(websocket.BinaryMessage, p)
	if err != nil {
		fmt.Errorf("write websocket message error: %+v", err)
		return len(p), err
	}

	return len(p), err
}

// WSS 's WriteMessage
func (wss *WSServer) WriteMessage(msg *model.Message) error {
	return wss.encodeAndPackPackage(msg)
}

// Set ReadDeadline 
func (wss *WSServer) SetReadDeadline(t time.Time) error {
	wss.ReadDeadline = t
	return wss.conn.SetReadDeadline(t)
}

// Set WriteDeadline 
func (wss *WSServer) SetWriteDeadline(t time.Time) error {
	wss.WriteDeadline = t
	return wss.conn.SetWriteDeadline(t)
}

// Get RemoteAddr 
func (wss *WSServer) RemoteAddr() net.Addr {
	return wss.conn.RemoteAddr()
}

// Get LocalAddr 
func (wss *WSServer) LocalAddr() net.Addr {
	return wss.conn.LocalAddr()
}

// Close connection 
func (wss *WSServer) CloseConnection() error {
	return wss.conn.Close()
}

func (wss *WSServer) handleRawData(){
	if wss.options.Consumer == nil {
		log.Println("bad consumer for raw data!")
		return 
	}

	if !wss.options.AutoRoute {
		return 
	}

	//Read the raw data
	_, err := io.Copy(wss.options.Consumer, wss)
	if err != nil {
		log.Println("failed to copy data, error:", err)
		wss.conn.Close()
		return
	}
}

func (wss *WSServer) ListenAndServeTLS() error {
	return wss.server.ListenAndServeTLS("", "")
}

func (wss *WSServer) Close() error {
	if wss.server != nil {
		return wss.server.Close()
	}

	return nil
}
