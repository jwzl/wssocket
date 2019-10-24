package model

import(
	"github.com/satori/go.uuid"
)

type MessageHeader struct {
	// the message id
	ID string `json:"ID,omitempty"`
	// message type
	Type string `json:"type,omitempty"`
	// the time of creating
	Timestamp int64 `json:"timestamp,omitempty"`
	//tag for other need
	Tag	 map[string]string `json:"tag,omitempty"`	
}

//MessageRoute contains structure of message
type MessageRoute struct {
	// where the message come from
	Source string `json:"source,omitempty"`
	// where the message will broadcasted to
	Group string `json:"group, omitempty"`
	// where the message come to
	Target string `json:"target,omitempty"`

	// what's the operation on resource
	Operation string `json:"operation,omitempty"`
	// what's the resource want to operate
	Resource string `json:"resource,omitempty"`
}

// Message struct
type Message struct {
	Header  MessageHeader `json:"header"`
	Router  MessageRoute  `json:"route,omitempty"`
	Content interface{}   `json:"content,omitempty"`
}

//GetID returns message ID
func (msg *Message) GetID() string {
	return msg.Header.ID
}

//BuildRouter sets route and resource operation in message
func (msg *Message) BuildRouter(source, group, target, res, opr string) *Message {
	msg.SetRoute(source, group, target)
	msg.SetResourceOperation(res, opr)
	return msg
}

//SetResourceOperation sets router resource and operation in message
func (msg *Message) SetResourceOperation(res, opr string) *Message {
	msg.Router.Resource = res
	msg.Router.Operation = opr
	return msg
}
//SetRoute sets router source and group in message
func (msg *Message) SetRoute(source, group, target string) *Message {
	msg.Router.Source = source
	msg.Router.Group = group
	msg.Router.Target = target
	return msg
}

//GetTimestamp returns message timestamp
func (msg *Message) GetTimestamp() int64 {
	return msg.Header.Timestamp
}

//GetContent returns message content
func (msg *Message) GetContent() interface{} {
	return msg.Content
}

//GetResource returns message route resource
func (msg *Message) GetResource() string {
	return msg.Router.Resource
}

//GetOperation returns message route operation string
func (msg *Message) GetOperation() string {
	return msg.Router.Operation
}

//GetSource returns message route source string
func (msg *Message) GetSource() string {
	return msg.Router.Source
}

//GetGroup returns message route group
func (msg *Message) GetGroup() string {
	return msg.Router.Group
}

//GetTarget returns message route Target
func (msg *Message) GetTarget() string {
	return msg.Router.Target
}

//UpdateID returns message object updating its ID
func (msg *Message) UpdateID() *Message {
	msg.Header.ID = uuid.NewV4().String()
	return msg
}

// BuildHeader builds message header. You can also use for updating message header
func (msg *Message) BuildHeader(ID, typ string, timestamp int64) *Message {
	msg.Header.ID = ID
	msg.Header.Type = typ
	msg.Header.Timestamp = timestamp
	return msg
}

//FillBody fills message  content that you want to send
func (msg *Message) FillBody(content interface{}) *Message {
	msg.Content = content
	return msg
}

// NewRawMessage returns a new raw message:
// model.NewRawMessage().BuildHeader().BuildRouter().FillBody()
func NewRawMessage() *Message {
	return &Message{}
}

// NewMessage returns a new basic message:
// model.NewMessage().BuildRouter().FillBody()
func NewMessage(parentID string) *Message {
	msg := &Message{}
	msg.Header.ID = uuid.NewV4().String()
	return msg
}
