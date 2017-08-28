package monitoring

import (
	"sync"
	"time"

	"github.com/emicklei/hazana"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

// Monitor is a Attack decorator that send metrics to graphite
type Monitor struct {
	hazana.Attack
	driver     *StackDriver
	mutex      *sync.Mutex
	dataPoints map[string][]*monitoringpb.Point
}

// NewMonitor returns a new Monitor decoration on an Attack
func NewMonitor(a hazana.Attack, s *StackDriver) *Monitor {
	return &Monitor{Attack: a, driver: s, dataPoints: []*monitoringpb.Point{}, mutex: new(sync.Mutex)}
}

// Do is part of hazana.Attack
func (m *Monitor) Do() hazana.DoResult {
	before := time.Now()
	result := m.Attack.Do()
	after := time.Now()
	m.mutex.Lock()
	points, ok := m.dataPoints[result.RequestLabel]
	if !ok {
		points = []*monitoringpb.Point{}
	}
	points = append(points, newDatapoint(after, after.Sub(before).Nanoseconds()))
	m.dataPoints[result.RequestLabel] = points
	m.mutex.UnLock()
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
		driver:     m.driver,
		mutext:     m.mutex,
		dataPoints: m.dataPoints,
	}
}
