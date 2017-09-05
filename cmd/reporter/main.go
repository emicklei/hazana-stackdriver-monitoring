package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/emicklei/hazana"
	metrics "github.com/emicklei/hazana-stackdriver-monitoring"
)

func main() {
	log.Println("hazana stackdriver reporting")
	// pick up a report from the arg and send it to stackdriver
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal("reading report failed ", err)
	}
	report := hazana.RunReport{}
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&report)
	if err != nil {
		log.Fatal("decoding report failed ", err)
	}
	if report.Metrics == nil {
		log.Fatal("no metrics to report")
	}
	driver, err := metrics.NewStackDriver(report.Configuration.Metadata["project_id"])
	if err != nil {
		log.Fatal("failed to create driver ", err)
	}
	err = driver.SendReport(report)
	if err != nil {
		log.Fatal("failed to send metrics ", err)
	}
}
