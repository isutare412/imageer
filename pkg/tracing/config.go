package tracing

type Config struct {
	Enabled          bool
	ServiceName      string
	SamplingRatio    float64 // 0.0 - 1.0
	OTLPGRPCEndpoint string
}
