package msg

import (
	"errors"
	"fmt"
	"github.com/limoges/p2pnet"
	"io"
	"net"
	"reflect"
)

const (
	HeaderLength = 4
	IPLength     = 16
)

var (
	ErrDataTooShort = errors.New("unmarshal: data is shorter than expected")
)

func createMessage(messageType int, payloadBuf []byte) (data []byte) {

	header := Header{
		Size: uint16(len(payloadBuf) + HeaderLength),
		Type: uint16(messageType),
	}
	headerBuf, _ := header.MarshalBinary()
	data = append(data, headerBuf...)
	data = append(data, payloadBuf...)
	return data
}

type Message interface {
	MinimumLength() int
	PayloadLength() int
	MarshalBinary() (data []byte, err error)
	UnmarshalBinary(data []byte) error
}

func MarshalBinary(m Message) {

	value := reflect.ValueOf(m)

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	//fmt.Println(reflect.TypeOf(value))

	length := uintptr(0)

	for i := 0; i < value.NumField(); i++ {
		fieldvalue := value.Field(i)
		fieldkind := value.Field(i).Kind()
		//fieldtype := value.Type().Field(i).Type
		fieldsize := value.Type().Field(i).Type.Size()
		fieldname := value.Type().Field(i).Name

		if fieldkind == reflect.Array || fieldkind == reflect.Slice {
			fieldlength := value.Field(i).Len()
			fieldsize = value.Field(i).Index(0).Type().Size() * uintptr(fieldlength)
		}

		if value.Field(i).CanInterface() {
			fieldinterface := value.Field(i).Interface()
			switch v := fieldinterface.(type) {
			case uint8:
			case uint16:
			case uint32:
			case uint64:
			case []byte:
			case bool:
			default:
				panic(fmt.Sprintf("type:%v not supported\n", v))
			}
			fmt.Printf("exported: %v, %v\n", fieldname, fieldvalue)
		} else {

			fmt.Printf("unexported:%v, %v\n", fieldname, fieldvalue)
		}
		length = length + fieldsize
	}

	fmt.Printf("Encoded message length:%v\n", length)

}

func NewMessage(h Header, payload []byte) Message {

	var m Message
	switch h.Type {
	case GOSSIP_ANNOUNCE:
		m = &GossipAnnounce{}
	case GOSSIP_NOTIFY:
		m = &GossipNotify{}
	case GOSSIP_NOTIFICATION:
		m = &GossipNotification{}
	case GOSSIP_VALIDATION:
		m = &GossipValidation{}
	case NSE_QUERY:
		m = &NSEQuery{}
	case NSE_ESTIMATE:
		m = &NSEEstimate{}
	case RPS_QUERY:
		m = &RPSQuery{}
	case RPS_PEER:
		m = &RPSPeer{}
	case ONION_TUNNEL_BUILD:
		m = &OnionTunnelBuild{}
	case ONION_TUNNEL_READY:
		m = &OnionTunnelReady{}
	case ONION_TUNNEL_INCOMING:
		m = &OnionTunnelIncoming{}
	case ONION_TUNNEL_DESTROY:
		m = &OnionTunnelDestroy{}
	case ONION_TUNNEL_DATA:
		m = &OnionTunnelData{}
	case ONION_ERROR:
		m = &OnionError{}
	case ONION_COVER:
		m = &OnionCover{}
	// case AUTH_SESSION_START:
	// case AUTH_SESSION_HS1:
	// case AUTH_SESSION_INCOMING_HS1:
	// case AUTH_SESSION_HS2:
	// case AUTH_SESSION_INCOMING_HS2:
	// case AUTH_LAYER_ENCRYPT:
	// case AUTH_LAYER_ENCRYPT_RESP:
	// case AUTH_LAYER_DECRYPT:
	// case AUTH_LAYER_DECRYPT_RESP:
	// case AUTH_SESSION_CLOSE:
	default:
		fmt.Printf("Unhandled message type: %v\n", h.Type)
	}
	m.UnmarshalBinary(payload)
	fmt.Println(m)
	return m
}

func Handle(conn net.Conn) {

	headerBuf := make([]byte, HeaderLength)

	// Build the token we use for logging purposes
	token := p2pnet.BuildConnexionIdentityToken(conn)
	fmt.Printf("[%v] New connexion opened by peer\n", token)

	// Run the handling loop, which is supposed to read packets, until
	// the connexion receives the EOF token, signifying that the peer
	// has closed the connexion.
	for {

		// First, we read the message's header.
		bytesRead, err := conn.Read(headerBuf)

		// Check read errors and connexion status.
		if err != nil {

			// EOF error is put up when the connexion has been closed by peer.
			if err == io.EOF {
				fmt.Printf("[%v] Connexion closed by peer.\n", token)
				// We simply stop handling this connexion.
				return
			}

			fmt.Println(err)
			continue
		}

		// Don't bother with the rest if the message is not long enough.
		header := Header{}
		err = header.UnmarshalBinary(headerBuf[:bytesRead])
		if err != nil {
			// We have a non-conforming header
			fmt.Println(err)
			continue
		}

		// Next we read the payload data.
		payloadBuf := make([]byte, header.PayloadSize())

		bytesRead, err = conn.Read(payloadBuf)
		if bytesRead != len(payloadBuf) {
			continue
		}

		packetBuf := make([]byte, 0)
		packetBuf = append(packetBuf, headerBuf...)
		packetBuf = append(packetBuf, payloadBuf...)

		// We now have the header and payload, just build the message.
		message := NewMessage(header, payloadBuf)
		fmt.Println(message)
	}
}
