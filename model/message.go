package model

import(
	"github.com/satori/go.uuid"
)

type MessageHeader struct {
	// the message id
	ID string `json:"ID,omitempty"`
	// message type
	MessageType string `json:"MessageType,omitempty"`
}

// Message struct
type Message struct {
	Header  MessageHeader `json:"header"`
	Content interface{}   `json:"content"`
}

//GetID returns message ID
func (msg *Message) GetID() string {
	return msg.Header.ID
}

//GetContent returns message content
func (msg *Message) GetContent() interface{} {
	return msg.Content
}

//UpdateID returns message object updating its ID
func (msg *Message) UpdateID() *Message {
	msg.Header.ID = uuid.NewV4().String()
	return msg
}

// BuildHeader builds message header. You can also use for updating message header
func (msg *Message) BuildHeader(ID string) *Message {
	msg.Header.ID = ID
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
