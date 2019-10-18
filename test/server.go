package main

import (
	"log"
	"time"
	"net/http"
	"github.com/jwzl/wssocket/conn"
	"github.com/jwzl/wssocket/model"	
	"github.com/jwzl/wssocket/server"
)

const (
	CA string ="/etc/edgedev/ca/rootCA.crt"
	Cert string ="/etc/edgedev/certs/edgedev.crt"
	key string ="/etc/edgedev/certs/edgedev.key"
)

type MessageHandler struct {}
func ( mh *MessageHandler ) MessageProcess(Header http.Header, msg *model.Message, c *conn.Connection){

}

func main() {
	mh := &MessageHandler{}
	server := &server.Server{
		Addr: "0.0.0.0:443",
		AutoRoute: true,
		HandshakeTimeout: 45 * time.Second,
		Handler: mh,
	}

	tlsConfig,err := server.CreateTLSConfig(CA, Cert, key)
	if err != nil {
		log.Println("Create tlsconfig err, %v", err)
		return 
	}
	server.TLSConfig = tlsConfig
	log.Println("Start the websocket server, listen: 127.0.0.1:443.....")
	server.StartServer("", "")
}
