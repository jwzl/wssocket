package server

import (
	"io/ioutil"
	"io"
	"fmt"	
	"crypto/tls"
	"crypto/x509"
	"time"
	"net/http"
	
	"github.com/jwzl/wssocket/conn"			
)

// the filter function before upgrading the http to websocket
type WSFilterFunc func(w http.ResponseWriter, r *http.Request) bool
type ConnNotifyFunc func (connection *conn.Connection) 
// Server 
type Server struct {
	// Addr optionally specifies the TCP address for the server to listen on,
	// in the form "host:port". If empty, ":http" (port 80) is used.
	// The service names are defined in RFC 6335 and assigned by IANA.
	// See net.Dial for details of the address format.
	Addr             string
	// When http recieive the request from client, then it 
	// means that connection is created.
	ConnNotify       ConnNotifyFunc		//optional.
	AutoRoute        bool
	//HandShake timeout
	HandshakeTimeout time.Duration

	//Message revice handler. 
	Handler          conn.Handler

	// this is for stream message
	Consumer         io.Writer		////optional.

	// the necessary processing before upgrading
	Filter 			 WSFilterFunc  //optional.

	// TLSConfig optionally provides a TLS configuration for use
	// by ServeTLS and ListenAndServeTLS.
	TLSConfig        *tls.Config
	server			 *WSServer
}

// create tls config
func (s *Server) CreateTLSConfig(caFile, certFile, keyFile string) (*tls.Config, error) {
	pool := x509.NewCertPool()
	rootCA, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, err
	}	
	ok := pool.AppendCertsFromPEM(rootCA)
	if !ok {
		return nil, fmt.Errorf("fail to load ca content")
	}
	
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		ClientCAs:    pool,
		ClientAuth:   tls.RequestClientCert,
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
		CipherSuites: []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256},
	}

	return tlsConfig, nil
}

// get tls config
func (s *Server) getTLSConfig(cert, key string) (*tls.Config, error) {
	var tlsConfig *tls.Config

	if s.TLSConfig == nil {
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

//For each connection.
func (s *Server) RangeConnection(f func(key, value interface{}) bool){
	s.server.cmgr.Range(f)
}
// Close
func (s *Server) Close() error {
	if s.server != nil {
		return s.server.Close()
	}

	return nil
}
