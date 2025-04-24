package services

import "sync/atomic"

type MetricsService struct {
	fileServerHits atomic.Int32
}

func NewMetricsService() *MetricsService {
	return &MetricsService{}
}

func (m *MetricsService) Increment() {
	m.fileServerHits.Add(1)
}

func (m *MetricsService) GetHits() int32 {
	return m.fileServerHits.Load()
}

func (m *MetricsService) Reset() {
	m.fileServerHits.Store(0)
}
