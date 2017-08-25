package monitoring

import (
	"context"
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

// Send will sends metrics to StackDriver using measurements of a samples.
func (s *StackDriver) Send(report hazana.RunReport) error {
	resourceLabels := map[string]string{"project_id": s.projectID}
	// collect labels from metadata
	for k, v := range report.Configuration.Metadata {
		if strings.HasPrefix(k, "resource.label.") {
			// Note: v must be a recognized resource label. TODO where is this documented?
			resourceLabels[k[len("resource.label."):]] = v
		}
	}
	resourceType, ok := report.Configuration.Metadata["resource.type"]
	if !ok {
		resourceType = "global"
	}
	resource := &monitoredrespb.MonitoredResource{
		Type:   resourceType,
		Labels: resourceLabels,
	}
	metricType, ok := report.Configuration.Metadata["metric.type"]
	if !ok {
		metricType = "custom.googleapis.com/missing-metric-type"
	}
	for sample, each := range report.Metrics {
		metric := &metricpb.Metric{
			Type: metricType,
			Labels: map[string]string{
				"sample": sample,
			},
		}
		dataPoint := newDatapoint(report.FinishedAt, (float64(each.Latencies.Mean.Nanoseconds()) / 1.0e6)) // ms
		if err := s.createTimeSeries(dataPoint, metric, resource); err != nil {
			return err
		}
	}
	return nil
}

func (s *StackDriver) createTimeSeries(
	dataPoint *monitoringpb.Point,
	metric *metricpb.Metric,
	resource *monitoredrespb.MonitoredResource) error {
	return s.client.CreateTimeSeries(s.ctx, &monitoringpb.CreateTimeSeriesRequest{
		Name: stackmoni.MetricProjectPath(s.projectID),
		TimeSeries: []*monitoringpb.TimeSeries{
			{
				Metric:   metric,
				Resource: resource,
				Points:   []*monitoringpb.Point{dataPoint},
			},
		},
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
