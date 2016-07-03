package msg

const (
	NSE_QUERY    = 520
	NSE_ESTIMATE = 521
	// Reserved up to 539.
)

type NSEQuery struct {
	// This is empty.
}

type NSEEstimate struct {
	EstimatePeers        uint32
	EstimateStdDeviation uint32
}
