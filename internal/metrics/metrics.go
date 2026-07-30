package metrics

import "sync/atomic"

type Metrics struct {
	requestsTotal       atomic.Uint64
	apiKeysCreatedTotal atomic.Uint64
	validationsTotal    atomic.Uint64
	rateLimitedTotal    atomic.Uint64
}

type Snapshot struct {
	RequestsTotal       uint64 `json:"requests_total"`
	APIKeysCreatedTotal uint64 `json:"api_keys_created_total"`
	ValidationsTotal    uint64 `json:"validations_total"`
	RateLimitedTotal    uint64 `json:"rate_limited_total"`
}

func New() *Metrics {
	return &Metrics{}
}

func (m *Metrics) IncRequests() {
	m.requestsTotal.Add(1)
}

func (m *Metrics) IncAPIKeysCreated() {
	m.apiKeysCreatedTotal.Add(1)
}

func (m *Metrics) IncValidations() {
	m.validationsTotal.Add(1)
}

func (m *Metrics) IncRateLimited() {
	m.rateLimitedTotal.Add(1)
}

func (m *Metrics) Snapshot() Snapshot {
	return Snapshot{
		RequestsTotal:       m.requestsTotal.Load(),
		APIKeysCreatedTotal: m.apiKeysCreatedTotal.Load(),
		ValidationsTotal:    m.validationsTotal.Load(),
		RateLimitedTotal:    m.rateLimitedTotal.Load(),
	}
}
