package conn

import(
	"crypto/x509"
	"encoding/json"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
	wstype "github.com/jwzl/wssocket/types"	
	"github.com/kubeedge/beehive/pkg/common/log"
)

// connection states
// TODO: add connection state filed
type ConnectionState struct {
	State            string
	Headers          http.Header
	PeerCertificates []*x509.Certificate
}

type Connection struct {
	//Message revice handler. 
	Handler          wstype.Handler
	// auto route flag
	AutoRoute bool
	// client type
	ConnUse string
	// Consumer
	Consumer io.Writer
	// State
	State  *ConnectionState
	//websocket Connection
	Conn *websocket.Conn
}

//Deep Copy http header.
func DeepCopyHeader(header http.Header) http.Header {
	headerByte, err := json.Marshal(header)
	if err != nil {
		log.LOGGER.Errorf("faile to marshal header, error:%+v", err)
		return nil
	}

	dstHeader := make(http.Header)
	err = json.Unmarshal(headerByte, &dstHeader)
	if err != nil {
		log.LOGGER.Errorf("failed to unmarshal header, error:%+v", err)
		return nil
	}
	return dstHeader
}
