package types

import (
	"net/http"
	"github.com/wssocket/model"	
)

const (

	// connection use type
	// connection only for message
	UseTypeMessage UseType = "msg"
	// connection only for stream
	UseTypeStream UseType = "str"
	// connection only can be used for message and stream
	UseTypeShare UseType = "shr"	
)

type Handler interface {
	MessageProcess(Header http.Header, msg *model.Message)
}

