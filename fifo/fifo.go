package fifo

import (
	"errors"
	"k8s.io/klog"
	"github.com/jwzl/wssocket/model"
)

const (
	defaultMaxFifoSize = 128
)

type MessageFifo struct {
	fifo chan *model.Message
}

// new message fifo.
func NewMessageFifo(fifoMaxSize int) *MessageFifo {
	if fifoMaxSize <= 0 {
		fifoMaxSize = defaultMaxFifoSize
	}

	return &MessageFifo{
		fifo: make(chan *model.Message, fifoMaxSize),
	}
}

// write the message into fifo. 
func (mf *MessageFifo) Write(msg *model.Message) {
	select {
	case mf.fifo <- msg:
	default:
		//discard a old message
		<- mf.fifo
		mf.fifo <- msg
		klog.Warningf("too many message, we discard a old message")
	}
}

// Read the message from fifo.
func (mf *MessageFifo) Read() (*model.Message, error) {
	msg, ok := <- mf.fifo
	if !ok {
		return nil, errors.New("this fifo is broken")
	}

	return msg, nil
}
