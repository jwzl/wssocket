package client

import (

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
	Handler          wstype.Handler
	// this is for stream message
	Consumer         io.Writer		////optional.
	// Connected callback
	Connected	func(*websocket.Conn, *http.Response)
	
}
