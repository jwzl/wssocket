package client

import (
	
	"fmt"
	"net/http"	
	"io/ioutil"

	"k8s.io/klog"
	"github.com/jwzl/wssocket/conn"
	"github.com/gorilla/websocket"
	wstype "github.com/jwzl/wssocket/types"
)

type WSClient struct {
	options Options
	//include all info to connect websocket server.
	dialer  *websocket.Dialer
}

// new websocket client 
func NewWSClient(opts Options) *WSClient {

	return &WSClient{
		options: opts,
		dialer: &websocket.Dialer{
			TLSClientConfig:  opts.TLSConfig,
			HandshakeTimeout: opts.HandshakeTimeout,
		},
	}
}

// try to connect remote server
// Use requestHeader to specify the
// origin (Origin), subprotocols (Sec-WebSocket-Protocol) and cookies (Cookie).
func (wsc *WSClient) Connect(serverAddr string, requestHeader http.Header)(*conn.Connection, error){
	header := requestHeader
	header.Add("ConnectionUse", string(wsc.options.ConnUse))

	klog.Infof("Connect %s", serverAddr)
	wsConn, response, err := wsc.dialer.Dial(serverAddr, header)
	if err == nil {
		klog.Infof("dialer connect %s successful", serverAddr)

		connection := &conn.Connection{
			ConnUse: wsc.options.ConnUse,
			Consumer:  wsc.options.Consumer,
			Handler:  wsc.options.Handler,  
			AutoRoute:  wsc.options.AutoRoute,
			State: &conn.ConnectionState{
				State:  wstype.StatConnected,
				Header: conn.DeepCopyHeader(header),	
			},
			Conn: wsConn, 
		}
		//call onconnect callback
		if wsc.options.Connected != nil {
			wsc.options.Connected(connection, response)	
		}
		// return the connection 
		return connection, nil
	}

	//failed.
	var resp string
	if response != nil {
		body, errRead := ioutil.ReadAll(response.Body)
		if errRead != nil {
			resp = fmt.Sprintf("response code: %d, response body: %s",response.StatusCode, string(body))
		}else {
			resp = fmt.Sprintf("response code: %d", response.StatusCode)
		}
		response.Body.Close()
		return nil, err
	} 

	klog.Errorf("dial websocket error(%+v), response message: %s", err, resp)
	return nil, err
}
