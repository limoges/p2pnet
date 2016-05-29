package msg

import (
	"fmt"
)

type AuthSessionStart struct {
	Hostkey []byte
}

func (m AuthSessionStart) String() string {
	return fmt.Sprintf("AUTH_SESSION_START:{%v}", m.Hostkey)
}

func (m AuthSessionStart) TypeId() uint16 {
	return AUTH_SESSION_START
}

func (m AuthSessionStart) LenInBytes() int {
	return len(m.Hostkey)
}

func NewAuthSessionStart(data []byte) (AuthSessionStart, error) {

	m := AuthSessionStart{}
	m.Hostkey = make([]byte, len(data))
	copy(m.Hostkey, data)
	return m, nil
}
