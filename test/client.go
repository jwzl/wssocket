package main 

import(
	"time"
	"net/http"
	"k8s.io/klog"

	"github.com/jwzl/wssocket/client"
	"github.com/jwzl/wssocket/conn"
	"github.com/jwzl/wssocket/model"
	wstype "github.com/jwzl/wssocket/types"	
)

type MessageHandler struct {}
func ( mh *MessageHandler ) MessageProcess(Header http.Header, msg *model.Message,  c *conn.Connection){

}

func Connected (conn *conn.Connection, resp *http.Response){
	klog.Infof("Connected!")
}
func main() {
	tlsConfig, err := client.CreateTLSConfig("/etc/edgedev/certs/edgedev.crt", "/etc/edgedev/certs/edgedev.key")
	if err != nil {
		klog.Errorf("Create tlsconfig err")
		return
	}

	options := client.Options{
		ConnUse: wstype.UseTypeMessage,
		TLSConfig: tlsConfig,
		HandshakeTimeout:	45 * time.Second,
		AutoRoute:  false,
		Handler:	&MessageHandler{},
		Connected:	Connected,
	}
	httpHeader := make(http.Header)
	httpHeader.Set("nodeid", "123")
	wsClient := &client.Client{
		Options: options,
		RequestHeader: httpHeader,
	}

	klog.Infof("Start websocket client...")
	wsClient.Start()
	klog.Infof("Connect to the server...")
	err = wsClient.Connect("wss://127.0.0.1:443/")
	if err != nil {
		klog.Infof("Connect failed, %v", err)
		return
	}
	
        	
	for {
	    time.Sleep(time.Second * 1)
	}
}
