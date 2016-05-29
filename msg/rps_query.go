package msg

type RPSQuery struct {
	// This is intended to be empty
}

func (m RPSQuery) TypeId() uint16 {
	return RPS_QUERY
}

func NewRPSQuery(data []byte) (RPSQuery, error) {

	m := RPSQuery{}
	return m, nil
}
