package monitoring

import (
	"context"
	"sync"
	"time"

	"github.com/emicklei/hazana"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

// Monitor is an Attack decorator that send latency metrics to Stackdriver
type Monitor struct {
	hazana.Attack
	mutex      *sync.Mutex
	dataPoints map[string][]*monitoringpb.Point
}

// NewMonitor returns a new Monitor decoration on an Attack
func NewMonitor(a hazana.Attack) *Monitor {
	return &Monitor{Attack: a, dataPoints: map[string][]*monitoringpb.Point{}, mutex: new(sync.Mutex)}
}

// Do is part of hazana.Attack
func (m *Monitor) Do(ctx context.Context) hazana.DoResult {
	before := time.Now()
	result := m.Attack.Do(ctx)
	after := time.Now()
	m.mutex.Lock()
	points, ok := m.dataPoints[result.RequestLabel]
	if !ok {
		points = []*monitoringpb.Point{}
	}
	points = append(points, newDatapoint(after, float64(after.Sub(before).Nanoseconds())))
	m.dataPoints[result.RequestLabel] = points
	m.mutex.Unlock()
	return result
}

// Setup is part of hazana.Attack
func (m *Monitor) Setup(c hazana.Config) error {
	return m.Attack.Setup(c)
}

// Clone is part of hazana.Attack
func (m *Monitor) Clone() hazana.Attack {
	return &Monitor{
		Attack: m.Attack.Clone(),
		// share the rest
		mutex:      m.mutex,
		dataPoints: m.dataPoints,
	}
}
