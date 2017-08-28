# hazana-stackdriver-monitoring
extension to the hazana package that will send metrics information to Stackdriver

# requirements

The configuration of hazana must have the following meta data entry:

    "metric.type": "custom.googleapis.com/YOUR-CUSTOM-PATH"

This package requires go version 1.9+

# using custom metrics
see https://cloud.google.com/monitoring/custom-metrics/creating-metrics