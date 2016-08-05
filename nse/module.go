package nse

import (
	"net"
	"time"

	"github.com/limoges/p2pnet"
	"github.com/limoges/p2pnet/cfg"
	"github.com/limoges/p2pnet/msg"
)

type NSE struct {
	EstimatedPeers     int
	EstimatedDeviation float32

	HistoryLength int
	Period        int
	APIAddr       string
}

const (
	ApiAddrToken         = "api_address"
	DefaultApiAddr       = "127.0.0.1:6001"
	HistoryLengthToken   = "history_length"
	DefaultHistoryLength = 10
	PeriodToken          = "period"
	DefaultPeriod        = 60
	ModuleToken          = "NSE"
)

func New(conf *cfg.Configurations) (*NSE, error) {

	var module *NSE

	module = &NSE{}
	conf.Init(&module.APIAddr, ModuleToken, ApiAddrToken, DefaultApiAddr)
	conf.Init(&module.HistoryLength, ModuleToken, HistoryLengthToken, DefaultHistoryLength)
	conf.Init(&module.Period, ModuleToken, PeriodToken, DefaultPeriod)

	return module, nil
}

func (m *NSE) Name() string {
	return ModuleToken
}

func (m *NSE) Addresses() (APIAddr, P2PAddr string) {
	return m.APIAddr, ""
}

func (m *NSE) Run() error {

	var c <-chan time.Time

	m.Calculate()

	c = time.Tick(time.Duration(m.Period) * time.Second)

	for {
		select {
		case <-c:
			m.Calculate()
		}
	}
	return nil
}

func (m *NSE) Handle(source net.Conn, message msg.Message) error {

	switch message.(type) {
	case msg.NSEQuery:
		return m.Reply(source)
	default:
		return p2pnet.ErrModuleDoesNotHandle
	}
}

func (m *NSE) Reply(conn net.Conn) error {

	var reply msg.NSEEstimate

	reply = msg.NSEEstimate{
		Peers:     uint32(m.EstimatedPeers),
		Deviation: uint32(m.EstimatedDeviation),
	}

	return msg.WriteMessage(conn, reply)
}

func (m *NSE) Calculate() {
	// Implement the update network estimation algorithm.
}
