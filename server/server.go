package server

import (
	"crypto/tls"
	"net"
	"net/http"

	wstype "github.com/wssocket/types"	
	"github.com/wssocket/model"			
)

// the filter function before upgrading the http to websocket
type WSFilterFunc func(w http.ResponseWriter, r *http.Request) bool
// Server 
type Server struct {
	// Addr optionally specifies the TCP address for the server to listen on,
	// in the form "host:port". If empty, ":http" (port 80) is used.
	// The service names are defined in RFC 6335 and assigned by IANA.
	// See net.Dial for details of the address format.
	Addr             string
	// When http recieive the request from client, then it 
	// means that connection is created.
	ConnNotify       ConnNotify		//optional.
	AutoRoute        bool
	//HandShake timeout
	HandshakeTimeout time.Duration

	//Message revice handler. 
	Handler          wstype.Handler

	// this is for stream message
	Consumer         io.Writer		////optional.

	// the necessary processing before upgrading
	Filter 			 WSFilterFunc  //optional.

	// TLSConfig optionally provides a TLS configuration for use
	// by ServeTLS and ListenAndServeTLS.
	TLSConfig        *tls.Config
	server			 *WSServer
}


// get tls config
func (s *Server) getTLSConfig(cert, key string) (*tls.Config, error) {
	var tlsConfig *tls.Config

	if s.tlsConfig == nil {
		tlsConfig = &tls.Config{}
	} else {
		tlsConfig = s.TLSConfig.Clone()
	}

	hasCert := false
	if len(tlsConfig.Certificates) > 0 ||
		tlsConfig.GetCertificate != nil {
		hasCert = true
	}
	if !hasCert || cert != "" || key != "" {
		var err error
		tlsConfig.Certificates = make([]tls.Certificate, 1)
		tlsConfig.Certificates[0], err = tls.LoadX509KeyPair(cert, key)
		if err != nil {
			return nil, err
		}
	}

	return tlsConfig, nil
}

// Start Server.
func (s *Server) StartServer(cert, key string) error {
	tlsConfig, err := s.getTLSConfig(cert, key)
	if err != nil {
		return err
	}

	s.TLSConfig = tlsConfig
	s.server = NewWSServer(*s)
	// Start server 
	err = s.server.ListenAndServeTLS() 

	return err
}

//WriteMessage
func (s *Server) WriteMessage(msg *model.Message)  error {
	if s.server != nil {
		return s.server.WriteMessage(msg)
	}

	return nil
}

//SetReadDeadline
func (s *Server) SetReadDeadline(t time.Time) error {
	if s.server != nil {
		return s.server.SetReadDeadline(t)
	}

	return nil
}

//SetWriteDeadline
func (s *Server) SetWriteDeadline(t time.Time) error {
	if s.server != nil {
		return s.server.SetWriteDeadline(t)
	}

	return nil
}

// RemoteAddr
func (s *Server) RemoteAddr() net.Addr {
	if s.server != nil {
		return s.server.RemoteAddr()
	}

	return nil
}

// LocalAddr
func (s *Server) LocalAddr() net.Addr {
	if s.server != nil {
		return s.server.LocalAddr()
	}

	return nil
}

//CloseConnection
func (s *Server) CloseConnection() error {
	if s.server != nil {
		return s.server.CloseConnection()
	}

	return nil
}

// Close
func (s *Server) Close() error {
	if s.server != nil {
		return s.server.Close()
	}

	return nil
}