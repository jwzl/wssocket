package types

import (
	"net/http"
	"github.com/jwzl/wssocket/model"	
)

const (
	// connection stat
	StatConnected    = "connected"
	StatDisconnected = "disconnected"
	// connection use type
	// connection only for message
	UseTypeMessage string = "msg"
	// connection only for stream
	UseTypeStream string = "str"
	// connection only can be used for message and stream
	UseTypeShare string = "shr"	
)

type Handler interface {
	MessageProcess(Header http.Header, msg *model.Message)
}

type ConnNotify func (interface{}) 
