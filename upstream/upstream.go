package upstream

type Upstream struct {
	Endpoint string
	Healthy  bool
	// client   interface{} // Placeholder for client providing eth data
}

func NewUpstream(endpoint string) *Upstream {
	return &Upstream{Endpoint: endpoint, Healthy: false}
}
