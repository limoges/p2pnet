package msg

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
)

// Header represents the protocol's common message header format.
// It is defined to contain message's size in the first 2 bytes,
// followed by the message's type in the following 2 bytes.
type GenericMessage struct {
	Size    uint16
	Type    uint16
	Content []byte
}

func (m GenericMessage) MessageType() uint16 {
	return m.Type
}

func (m GenericMessage) String() string {
	return fmt.Sprintf("%v#%v[%v]", Identifier(m.Type), m.Type, m.Size)
}

func Identifier(messageType uint16) string {

	switch messageType {
	case GOSSIP_ANNOUNCE:
		return "GOSSIP_ANNOUNCE"
	case GOSSIP_NOTIFY:
		return "GOSSIP_NOTIFY"
	case GOSSIP_NOTIFICATION:
		return "GOSSIP_NOTIFICATION"
	case GOSSIP_VALIDATION:
		return "GOSSIP_VALIDATION"
	case NSE_QUERY:
		return "NSE_QUERY"
	case NSE_ESTIMATE:
		return "NSE_ESTIMATE"
	case RPS_QUERY:
		return "RPS_QUERY"
	case RPS_PEER:
		return "RPS_PEER"
	case ONION_TUNNEL_BUILD:
		return "ONION_TUNNEL_BUILD"
	case ONION_TUNNEL_READY:
		return "ONION_TUNNEL_READY"
	case ONION_TUNNEL_INCOMING:
		return "ONION_TUNNEL_INCOMING"
	case ONION_TUNNEL_DESTROY:
		return "ONION_TUNNEL_DESTROY"
	case ONION_TUNNEL_DATA:
		return "ONION_TUNNEL_DATA"
	case ONION_ERROR:
		return "ONION_ERROR"
	case ONION_COVER:
		return "ONION_COVER"
	case AUTH_SESSION_START:
		return "AUTH_SESSION_START"
	case AUTH_SESSION_HS1:
		return "AUTH_SESSION_HS1"
	case AUTH_SESSION_INCOMING_HS1:
		return "AUTH_SESSION_INCOMING_HS1"
	case AUTH_SESSION_HS2:
		return "AUTH_SESSION_HS2"
	case AUTH_SESSION_INCOMING_HS2:
		return "AUTH_SESSION_INCOMING_HS2"
	case AUTH_LAYER_ENCRYPT:
		return "AUTH_LAYER_ENCRYPT"
	case AUTH_LAYER_ENCRYPT_RESP:
		return "AUTH_LAYER_ENCRYPT_RESP"
	case AUTH_LAYER_DECRYPT:
		return "AUTH_LAYER_DECRYPT"
	case AUTH_LAYER_DECRYPT_RESP:
		return "AUTH_LAYER_DECRYPT_RESP"
	case AUTH_SESSION_CLOSE:
		return "AUTH_SESSION_CLOSE"
	case AUTH_HANDSHAKE1:
		return "AUTH_HANDSHAKE1"
	case AUTH_HANDSHAKE2:
		return "AUTH_HANDSHAKE2"
	case AUTH_SESSION_CONFIRMED:
		return "AUTH_SESSION_CONFIRMED"
	case AUTH_SESSION_DECLINED:
		return "AUTH_SESSION_DECLINED"
	default:
		return "UNKNOWN_MESSAGE"
	}
}

func WriteGenericMessage(writer io.Writer, message GenericMessage) error {

	// Write the message's field to the buffer.
	value := reflect.Indirect(reflect.ValueOf(message))
	buf := new(bytes.Buffer)
	for i := 0; i < value.NumField(); i++ {

		field := value.Field(i).Interface()
		if err := binary.Write(buf, binary.BigEndian, field); err != nil {
			return err
		}
	}

	// Write the buffer to the writer.
	if _, err := writer.Write(buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func ReadGenericMessage(reader *bufio.Reader) (GenericMessage, error) {

	m := GenericMessage{}

	// Read the message size
	if err := binary.Read(reader, binary.BigEndian, &m.Size); err != nil {
		return m, err
	}

	// Read the message type
	if err := binary.Read(reader, binary.BigEndian, &m.Type); err != nil {
		return m, err
	}

	// Calculate the length left to read.
	mustRead := m.Size - HeaderLength
	buf := make([]byte, mustRead)

	// Copy the content into the generic message
	if _, err := io.ReadFull(reader, buf); err != nil {
		return m, err
	}
	m.Content = buf

	return m, nil
}
func ConvertToGeneric(message Message) (GenericMessage, error) {

	m := GenericMessage{}
	// Write the content of the message to a buffer
	content := new(bytes.Buffer)
	value := reflect.Indirect(reflect.ValueOf(message))
	for i := 0; i < value.NumField(); i++ {

		field := value.Field(i).Interface()
		if err := binary.Write(content, binary.BigEndian, field); err != nil {
			fmt.Printf("Could not write field %v\n", i)
			return m, err
		}
	}

	length := content.Len() + HeaderLength
	// Calculate the complete message's length
	if length > math.MaxUint16 {
		return m, errors.New(
			fmt.Sprintf("Message lenth must be smaller than %v. Current message length is %v\n",
				math.MaxUint16, length))
	}
	messageLength := uint16(content.Len() + HeaderLength)
	messageType := message.TypeId()

	m.Size = messageLength
	m.Type = messageType
	m.Content = content.Bytes()
	return m, nil
}

func ConvertFromGeneric(generic GenericMessage) (Message, error) {

	var m Message
	var err error

	switch generic.Type {
	// case GOSSIP_ANNOUNCE:
	// 	m = &GossipAnnounce{}
	// case GOSSIP_NOTIFY:
	// 	m = &GossipNotify{}
	// case GOSSIP_NOTIFICATION:
	// 	m = &GossipNotification{}
	// case GOSSIP_VALIDATION:
	// 	m = &GossipValidation{}
	// case NSE_QUERY:
	// 	m = &NSEQuery{}
	// case NSE_ESTIMATE:
	// 	m = &NSEEstimate{}
	// case RPS_QUERY:
	// 	m = &RPSQuery{}
	// case RPS_PEER:
	// 	m = &RPSPeer{}
	case ONION_TUNNEL_BUILD:
		m, err = NewOnionTunnelBuild(generic.Content)
	case ONION_TUNNEL_READY:
		m, err = NewOnionTunnelReady(generic.Content)
	// case ONION_TUNNEL_INCOMING:
	// 	m = &OnionTunnelIncoming{}
	// case ONION_TUNNEL_DESTROY:
	// 	m = &OnionTunnelDestroy{}
	// case ONION_TUNNEL_DATA:
	// 	m = &OnionTunnelData{}
	// case ONION_ERROR:
	// 	m = &OnionError{}
	// case ONION_COVER:
	//	m = &OnionCover{}
	case AUTH_SESSION_START:
		m, err = NewAuthSessionStart(generic.Content)
	case AUTH_SESSION_HS1:
		m, err = NewAuthSessionHS1(generic.Content)
	case AUTH_SESSION_INCOMING_HS1:
		m, err = NewAuthSessionIncomingHS1(generic.Content)
	case AUTH_SESSION_HS2:
		m, err = NewAuthSessionHS2(generic.Content)
	case AUTH_SESSION_INCOMING_HS2:
		m, err = NewAuthSessionIncomingHS2(generic.Content)
	case AUTH_LAYER_ENCRYPT:
		m, err = NewAuthLayerEncrypt(generic.Content)
	case AUTH_LAYER_ENCRYPT_RESP:
		m, err = NewAuthLayerEncryptResp(generic.Content)
	case AUTH_LAYER_DECRYPT:
		m, err = NewAuthLayerDecrypt(generic.Content)
	case AUTH_LAYER_DECRYPT_RESP:
		m, err = NewAuthLayerDecryptResp(generic.Content)
	case AUTH_SESSION_CLOSE:
		m, err = NewAuthSessionClose(generic.Content)
	case AUTH_HANDSHAKE1:
		m, err = NewAuthHandshake1(generic.Content)
	case AUTH_HANDSHAKE2:
		m, err = NewAuthHandshake2(generic.Content)
	case AUTH_SESSION_CONFIRMED:
		m, err = NewAuthSessionConfirmed(generic.Content)
	case AUTH_SESSION_DECLINED:
		m, err = NewAuthSessionDeclined(generic.Content)
	default:
		fmt.Printf("Type cannot be converted from generic: %v\n", generic.Type)
		panic("Need to implement in msg/messages.go")
	}
	if err != nil {
		fmt.Println(err)
	}
	return m, err
}
