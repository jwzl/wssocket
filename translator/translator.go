package translator

import(
	"fmt"
	"encoding/json"

	"github.com/jwzl/wssocket/model"	
	"github.com/jwzl/wssocket/protoc/message"	
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
	modelMsg.Header.Type = proMsg.Header.Type
	modelMsg.Header.Timestamp = proMsg.Header.Timestamp
	modelMsg.Header.Tag  = proMsg.Header.Tag 

	modelMsg.Router.Source = proMsg.Router.Source
	modelMsg.Router.Group = proMsg.Router.Group
	modelMsg.Router.Target = proMsg.Router.Target
	modelMsg.Router.Operation = proMsg.Router.Operation
	modelMsg.Router.Resource = proMsg.Router.Resource 
	return nil
}
func (mtc *MessageTransCoding) modeltopro(modelMsg *model.Message, proMsg *message.Message) error {
	proMsg.Header.ID = modelMsg.Header.ID
	proMsg.Header.Type = modelMsg.Header.Type
	proMsg.Header.Timestamp = modelMsg.Header.Timestamp
	proMsg.Header.Tag  = modelMsg.Header.Tag

	//Router
	proMsg.Router.Source = modelMsg.Router.Source
	proMsg.Router.Group = modelMsg.Router.Group
	proMsg.Router.Target = modelMsg.Router.Target
	proMsg.Router.Operation = modelMsg.Router.Operation
	proMsg.Router.Resource = modelMsg.Router.Resource 

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
	err := proto.Unmarshal(raw, &protoMsg)
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
		return nil, fmt.Errorf("bad msg type")
	}

	protoMsg := message.Message{
		Header: &message.MessageHeader{},
	}

	err := mtc.modeltopro(modelMsg, &protoMsg)
	if err != nil {
		return nil, err
	}

	msgBytes, err :=proto.Marshal(&protoMsg)
	if err != nil {
		fmt.Errorf("protocol buf marshal err!")
		return nil, err
	}

	return msgBytes, nil
}

