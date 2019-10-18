package server

import (
	"log"
	"io"	
	"os"
	"fmt"
	"net"
	"time"
	"net/http"
	"k8s.io/klog"

	"github.com/gorilla/websocket"
	wstype "github.com/jwzl/wssocket/types"	
	"github.com/jwzl/wssocket/model"	
	"github.com/jwzl/wssocket/packer"	
	"github.com/jwzl/wssocket/translator"	
)

type WSServer struct {
	options Server
	server  *http.Server
	cmgr	*ConnectionManager
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
		cmgr:  &ConnectionManager{
			ConnKey: getDefaultConnKey,
		}
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
		klog.Errorf("failed to upgrade to websocket")
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
	klog.Infof("request coming....%v", r)
	WsConn := wss.upgrade(w, r)
	if WsConn == nil {
		return	
	}
	// Create a connection
	connection := &conn.Connection{
		ConnUse: string(r.Header.Get("ConnectionUse")),
		Consumer:  wss.options.Consumer,
		Handler:  wss.options.Handler,  
		AutoRoute:  wss.options.AutoRoute,
		State: &conn.ConnectionState{
			State:  wstype.StatConnected,
			Header: conn.DeepCopyHeader(r.Header),	
		},
		Conn: WsConn, 
	}
	
	// Connection callback.
	if wss.options.ConnNotify != nil {
		wss.options.ConnNotify(connection)
	}

	//Connection manager
	if wss.cmgr	!= nil {
		wss.cmgr.AddConnection(connection)
	} 
	//Handle connection
	go connection.ConnRecieve()
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

//Simple Connection manager
type ConnectionKey  func(connection conn.Connection)string
type ConnectionManager struct{
	ConnKey  ConnectionKey
	Connections sync.Map
}

func getDefaultConnKey(conn *conn.Connection) string {
	return conn.RemoteAddr().String()
}

//Add connection
func (cm* ConnectionManager) AddConnection(conn *conn.Connection){
	cm.Connections.Store(cm.ConnKey(conn), conn)
}

//Del connection
func (cm* ConnectionManager) DelConnection(conn *conn.Connection){
	cm.Connections.Delete(cm.ConnKey(conn))
}

//Get connection
func (cm* ConnectionManager) GetConnection(key string) *conn.Connection {
	obj, ok := cm.Connections.Load(key)
	if ok {
		return obj.(*conn.Connection)
	}

	return nil
}

//Range
func (cm* ConnectionManager) Range(f func(key, value interface{}) bool){
	cm.Connections.Range(f)
}
