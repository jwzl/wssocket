package server

import (
	"log"
	"os"
	"net/http"

	"github.com/gorilla/websocket"
	wstype "github.com/wssocket/types"	
	"github.com/wssocket/model"		
)

// the filter function before upgrading the http to websocket
type WSFilterFunc func(w http.ResponseWriter, r *http.Request) bool
// server options
type Options struct {
	// Addr optionally specifies the TCP address for the server to listen on,
	// in the form "host:port". If empty, ":http" (port 80) is used.
	// The service names are defined in RFC 6335 and assigned by IANA.
	// See net.Dial for details of the address format.
	Addr             string
	//websocket use 
	ConnUse 		 string
	// TLSConfig optionally provides a TLS configuration for use
	// by ServeTLS and ListenAndServeTLS.
	TLS              *tls.Config
	ConnNotify       ConnNotify
	ConnMgr          *cmgr.ConnectionManager
	ConnNumMax       int
	AutoRoute        bool
	HandshakeTimeout time.Duration
	Handler          mux.Handler
	Consumer         io.Writer
	// the necessary processing before upgrading
	Filter 			 WSFilterFunc
}

type WSServer struct {
	options Options
	server  *http.Server
	conn	*websocket.Conn
}

func NewWSServer(ops Options) *WSServer {

	server := http.Server{
		Addr: 		ops.Addr,
		TLSConfig: 	opts.TLS,
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
		wss.options.ConnNotify()
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

func (wss *WSServer) handleMessage(){
	msg := &model.Message{}
	for {
		

	}
}

func (wss *WSServer) handleRawData(){
	if wss.options.Consumer == nil {
		log.Println("bad consumer for raw data!")
		return 
	}

	if wss.options.AutoRoute != nil {
		return 
	}

	//Read the raw data
	_, err := io.Copy(wss.options.Consumer, wss.conn)
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
