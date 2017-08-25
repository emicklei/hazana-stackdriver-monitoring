package monitoring

import (
	"testing"
	"time"

	"github.com/emicklei/hazana"
)

func TestSendReport(t *testing.T) {
	t.Skip() // first change YOURPROJECT into a StackDriver enabled GCP project
	hm := new(hazana.Metrics)
	hm.Latencies.Mean = time.Duration(42000000)
	when := time.Now()
	r := hazana.RunReport{
		FinishedAt: when,
		Configuration: hazana.Config{
			Metadata: map[string]string{
				"metric.type": "custom.googleapis.com/myservice",
			},
		},
		Metrics: map[string]*hazana.Metrics{"testsample": hm},
	}
	d, err := NewStackDriver("YOURPROJECT")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(d.Send(r))
}
