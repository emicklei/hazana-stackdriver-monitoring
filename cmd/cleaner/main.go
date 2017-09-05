package main

// remove an obsolete (automatic) custom metric

import (
	"context"
	"log"
	"os"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	pb "google.golang.org/genproto/googleapis/monitoring/v3"
)

// go run main.go projects/<project>/metricDescriptors/custom.googleapis.com%2Fstepconnector%2Fv1

// DOES NOT WORK At the time of writing

// 2017/09/04 15:40:55 Get error: rpc error: code = PermissionDenied desc = User <user> does not have permission to see metric custom.googleapis.com/stepconnector/v1
// 2017/09/04 15:40:56 Delete error:rpc error: code = InvalidArgument desc = Field name had an invalid value of "custom.googleapis.com/stepconnector/v1": The metric type must be a URL-formatted string with a domain and non-empty path.

func main() {
	metric := os.Args[1]
	ctx := context.Background()
	client, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	{
		req := &pb.GetMetricDescriptorRequest{
			Name: metric,
		}
		if resp, err := client.GetMetricDescriptor(ctx, req); err != nil {
			log.Println("Get error:", err)
		} else {
			log.Println(resp.GetDisplayName())
		}

	}

	{
		req := &pb.DeleteMetricDescriptorRequest{
			Name: metric,
		}
		if err := client.DeleteMetricDescriptor(ctx, req); err != nil {
			log.Fatal("Delete error:", err)
		}
	}
	log.Println("deleted ", metric)
}
