package monitoring

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/emicklei/hazana"

	stackmoni "cloud.google.com/go/monitoring/apiv3"
	googlepb "github.com/golang/protobuf/ptypes/timestamp"
	metricpb "google.golang.org/genproto/googleapis/api/metric"
	monitoredrespb "google.golang.org/genproto/googleapis/api/monitoredres"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

// StackDriver provides the api to send a hazana.RunReport
type StackDriver struct {
	client    *stackmoni.MetricClient
	projectID string
	ctx       context.Context
}

// NewStackDriver create a connected StackDriver for a given project.
func NewStackDriver(projectID string) (*StackDriver, error) {
	ctx := context.Background()
	client, err := stackmoni.NewMetricClient(ctx)
	if err != nil {
		return nil, err
	}
	return &StackDriver{client: client, projectID: projectID, ctx: ctx}, nil
}

// SendReport will sends metrics to StackDriver using measurements of a samples.
func (s *StackDriver) SendReport(report hazana.RunReport) error {
	if report.Metrics == nil || len(report.Metrics) == 0 {
		return nil
	}
	metricType := s.metricType(report.Configuration)
	resource := s.newResource(report.Configuration)

	timeSeries := []*monitoringpb.TimeSeries{}
	for sample, each := range report.Metrics {
		for _, point := range []struct {
			key   string
			value float64
		}{
			{key: "mean", value: float64(each.Latencies.Mean.Nanoseconds()) / 1.0e6}, // ms
			{key: "max", value: float64(each.Latencies.Max.Nanoseconds()) / 1.0e6},   // ms
			{key: "99th", value: float64(each.Latencies.P99.Nanoseconds()) / 1.0e6},  // ms
			{key: "success", value: each.Success * 100},
			{key: "count", value: float64(each.Requests)},
			{key: "rate", value: float64(each.Rate)},
		} {
			metric := &metricpb.Metric{
				Type: metricType,
				Labels: map[string]string{
					"requestLabel": sample,
					"field":        point.key,
				},
			}
			dataPoint := newDatapoint(report.FinishedAt, point.value)
			timeSeries = append(timeSeries, newTimeSeries(dataPoint, metric, resource))
		}
	}
	if err := s.createTimeSeries(timeSeries); err != nil {
		return err
	}
	return nil
}

// SendMonitor sends the datapoints as timeseries to Stackdriver
func (s *StackDriver) SendMonitor(monitor *Monitor, config hazana.Config) error {
	timeSeries := []*monitoringpb.TimeSeries{}
	resource := s.newResource(config)
	for label, points := range monitor.dataPoints {
		metric := &metricpb.Metric{
			Type: s.metricType(config),
			Labels: map[string]string{
				"requestLabel": label,
				"field":        "duration",
			},
		}
		series := &monitoringpb.TimeSeries{
			Metric:   metric,
			Resource: resource,
			Points:   points,
		}
		timeSeries = append(timeSeries, series)
		if config.Verbose {
			log.Printf("collected [%d] datapoints for label [%s]\n", len(points), label)
		}
	}
	if err := s.createTimeSeries(timeSeries); err != nil {
		return err
	}
	return nil
}

func (s *StackDriver) metricType(config hazana.Config) string {
	metricType, ok := config.Metadata["metric.type"]
	if !ok {
		metricType = "custom.googleapis.com/missing-metric-type"
	}
	return metricType
}

func (s *StackDriver) newResource(config hazana.Config) *monitoredrespb.MonitoredResource {
	resourceLabels := map[string]string{"project_id": s.projectID}
	// collect labels from metadata
	for k, v := range config.Metadata {
		if strings.HasPrefix(k, "resource.label.") {
			// Note: v must be a recognized resource label. https://cloud.google.com/monitoring/custom-metrics/creating-metrics
			resourceLabels[k[len("resource.label."):]] = v
		}
	}
	resourceType, ok := config.Metadata["resource.type"]
	if !ok {
		resourceType = "global"
	}
	return &monitoredrespb.MonitoredResource{
		Type:   resourceType,
		Labels: resourceLabels,
	}
}

func newTimeSeries(dataPoint *monitoringpb.Point,
	metric *metricpb.Metric,
	resource *monitoredrespb.MonitoredResource) *monitoringpb.TimeSeries {
	return &monitoringpb.TimeSeries{
		Metric:   metric,
		Resource: resource,
		Points:   []*monitoringpb.Point{dataPoint},
	}
}

func (s *StackDriver) createTimeSeries(timeSeries []*monitoringpb.TimeSeries) error {
	return s.client.CreateTimeSeries(s.ctx, &monitoringpb.CreateTimeSeriesRequest{
		Name:       stackmoni.MetricProjectPath(s.projectID),
		TimeSeries: timeSeries,
	})
}

func newDatapoint(when time.Time, d float64) *monitoringpb.Point {
	return &monitoringpb.Point{
		Interval: &monitoringpb.TimeInterval{
			// for gauge metric StartTime must be the same as EndTime or zero
			EndTime: &googlepb.Timestamp{Seconds: when.Unix()},
		},
		Value: &monitoringpb.TypedValue{
			Value: &monitoringpb.TypedValue_DoubleValue{DoubleValue: d},
		},
	}
}
