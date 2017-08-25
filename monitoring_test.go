package monitoring

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/emicklei/hazana"
)

var report = `
{
	"startedAt": "2017-08-25T15:54:41.43681012+02:00",
	"finishedAt": "2017-08-25T15:54:51.470388711+02:00",
	"configuration": {
		"rps": 10,
		"attackTimeSec": 20,
		"rampupTimeSec": 10,
		"maxAttackers": 10,
		"verbose": true,
		"metadata": {
			"metric.type": "custom.googleapis.com/myservice",
		}
	},
	"metrics": {
		"item.xml": {
			"latencies": {
				"total": 2174448417,
				"mean": 43488968,
				"50th": 38301367,
				"95th": 48652985,
				"99th": 141567896,
				"max": 153030573
			},
			"earliest": "2017-08-25T15:54:41.540237194+02:00",
			"latest": "2017-08-25T15:54:51.237186313+02:00",
			"end": "2017-08-25T15:54:51.277747013+02:00",
			"duration": 9696949119,
			"wait": 40560700,
			"requests": 50,
			"rate": 5.156260942117459,
			"success": 1,
			"status_codes": {
				"200": 50
			},
			"errors": null
		},
		"variant.xml": {
			"latencies": {
				"total": 2182407327,
				"mean": 42792300,
				"50th": 39192429,
				"95th": 58143997,
				"99th": 62634128,
				"max": 112791958
			},
			"earliest": "2017-08-25T15:54:41.436828733+02:00",
			"latest": "2017-08-25T15:54:51.437223555+02:00",
			"end": "2017-08-25T15:54:51.470282401+02:00",
			"duration": 10000394822,
			"wait": 33058846,
			"requests": 51,
			"rate": 5.099798648729791,
			"success": 1,
			"status_codes": {
				"200": 51
			},
			"errors": null
		}
	}
}
`

func TestSendReport(t *testing.T) {
	//t.Skip() // first change YOURPROJECT into a StackDriver enabled GCP project
	r := hazana.RunReport{}
	json.NewDecoder(strings.NewReader(report)).Decode(&r)
	//d, err := NewStackDriver("YOURPROJECT")
	d, err := NewStackDriver("kramp-hub")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(d.Send(r))
}
