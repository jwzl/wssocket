package translator

import(
	"fmt"
	"encoding/json"

	"github.com/wssocket/model"	
	"github.com/wssocket/protoc/message"	
	"github.com/golang/protobuf/proto"
)

type MessageTransCoding struct {
}

func NewTransCoding() *MessageTransCoding {
	return &MessageTransCoding{}
}

func (mtc *MessageTransCoding) protomodel(proMsg *message.Message, modelMsg *model.Message) error {
	modelMsg.FillBody(proMsg.Content)
	modelMsg.Header.ID = proMsg.Header.ID
	modelMsg.Header.MessageType = proMsg.Header.MessageType

	return nil
}
func (mtc *MessageTransCoding) modeltopro(modelMsg *model.Message, proMsg *message.Message) error {
	proMsg.Header.ID = modelMsg.Header.ID
	proMsg.Header.MessageType = modelMsg.Header.MessageType

	if content := modelMsg.GetContent(); content!= nil {
		switch content.(type) {
		case []byte:
			proMsg.Content = content.([]byte)	
		case string:
			proMsg.Content = []byte(content.(string))
		default:
			bytes, err := json.Marshal(content)
			if err != nil {
				fmt.Errorf("failed to marshal")
				return err
			}
			proMsg.Content = bytes
		}
	}

	return nil
}


//Convert raw data to protocol buf message, then convert to model Message
func (mtc *MessageTransCoding) Decode(raw []byte, msg interface{}) error {
	modelMsg, ok := msg.(*model.Message)
	if !ok {
		return fmt.Errorf("bad msg type")
	}

	protoMsg := message.Message{}
	//unmarshall the data to  protocol buf message.
	err := proto.Unmarshal(raw, &protoMessage)
	if err != nil {
		fmt.Errorf("protocol buf Unmarshal err!")
		return err
	}

	mtc.protomodel(&protoMsg, modelMsg)
	return nil	
}

////Convert model message to protocol buf message, then convert to byte array.
func (mtc *MessageTransCoding) Encode(msg interface{}) ([]byte, error) {
	modelMsg, ok := msg.(*model.Message)
	if !ok {
		return fmt.Errorf("bad msg type")
	}

	protoMsg := message.Message{
		Header: &message.MessageHeader{},
	}

	err := mtc.modeltopro(modelMsg, &protoMsg)
	if err != nil {
		return nil, err
	}

	msgBytes, err :=proto.marshal(&protoMsg)
	if err != nil {
		fmt.Errorf("protocol buf marshal err!")
		return nil, err
	}

	return msgBytes, nil
}

