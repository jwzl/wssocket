package client

import (

	"fmt"
	"io/ioutil"

	"github.com/jwzl/wssocket/conn"
	"github.com/gorilla/websocket"
	wstype "github.com/jwzl/wssocket/types"	
	"github.com/kubeedge/beehive/pkg/common/log"
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

	wsConn, response, err := wsc.dialer.Dial(serverAddr, header)
	if err == nil {
		log.Infof("dialer connect %s successful", serverAddr)

		//call onconnect callback
		if wsc.options.Connected {
			wsc.options.Connected(wsConn, response)	
		}
		// return the connection 
		return &conn.Connection{
			ConnUse: wsc.options.ConnUse,
			Consumer:  wsc.options.Consumer,
			Handler:  wsc.options.Handler,  
			AutoRoute:  wsc.options.AutoRoute,
			State: &conn.ConnectionState{
				State:  wstype.StatConnected,
				Headers: conn.DeepCopyHeader(header),	
			},
			Conn: wsConn, 
		}, nil
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

	log.Errorf("dial websocket error(%+v), response message: %s", err, resp)
	return nil, err
}
