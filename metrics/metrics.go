package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"

	mutils "github.com/trustedanalytics/metrics/utils"
)

var tapCounts = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: "tap",
		Subsystem: "catalog",
		Name:      "counts",
		Help:      "Count of various TAP components",
	}, []string{"component"})

func collectApplicationsCount() error {
	// TODO
	tapCounts.WithLabelValues("applications").Set(43)
	return nil
}

func collectServicesCount() error {
	// TODO
	tapCounts.WithLabelValues("services").Set(43)
	return nil
}

func collectServiceInstancesCount() error {
	// TODO
	tapCounts.WithLabelValues("serviceInstances").Set(43)
	return nil
}

func EnableCollection(delay time.Duration) chan<- struct{} {
	mutils.RegisterMetrics("catalog", tapCounts)
	return mutils.EnableMetricsCollecting(delay,
		collectApplicationsCount,
		collectServicesCount,
		collectServiceInstancesCount,
	)
}
