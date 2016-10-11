package utils

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricCollector func() error

func metricCounter(component, status string) prometheus.Counter {
	return prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "tap",
			Subsystem: "metrics",
			Name:      "collecting_count",
			Help:      "Number of metric collecting rounds",
			ConstLabels: map[string]string{
				"component": component,
				"status":    status,
			},
		})
}

var (
	metricsCollectingStarted    prometheus.Counter
	metricsCollectingSuccessful prometheus.Counter
	metricsCollectingFailed     prometheus.Counter
	metricsCollectingTime       prometheus.Summary
)

func RegisterMetrics(componentName string, cs ...prometheus.Collector) {
	prometheus.MustRegister(cs...)
	metricsCollectingStarted = metricCounter(componentName, "started")
	metricsCollectingSuccessful = metricCounter(componentName, "sucessfull")
	metricsCollectingFailed = metricCounter(componentName, "failed")
	metricsCollectingTime = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Namespace:   "tap",
			Subsystem:   "metrics",
			Name:        "collecting_duration_nanoseconds",
			Help:        "Time taken for collecting meitrcs",
			ConstLabels: map[string]string{"component": componentName},
		},
	)
	prometheus.MustRegister(
		metricsCollectingStarted,
		metricsCollectingSuccessful,
		metricsCollectingFailed,
		metricsCollectingTime,
	)
}

func EnableMetricsCollecting(delay time.Duration, collectors ...MetricCollector) chan<- struct{} {
	done := make(chan struct{})
	EnableMetricsCollectingWithEnd(delay, done, collectors...)
	return done
}

// Explicitly pass end signal
func EnableMetricsCollectingWithEnd(delay time.Duration, done <-chan struct{}, collectors ...MetricCollector) {
	go collectingLoop(done, delay, collectors)
}

func GetHandler() http.Handler {
	return promhttp.Handler()
}

func collectingLoop(done <-chan struct{}, delay time.Duration, collectors []MetricCollector) {
	log.Println("Started metrics collecting loop")
	for {
		select {
		case <-time.After(delay):
			go collectMetrics(collectors)
		case <-done:
			log.Println("Finishing collecting loop")
			return
		}
	}
}

func collectMetrics(collectors []MetricCollector) {
	metricsCollectingStarted.Inc()
	startTime := time.Now()
	wasError := false

	defer func() {
		if wasError {
			metricsCollectingFailed.Inc()
		} else {
			metricsCollectingSuccessful.Inc()
		}
		took := time.Now().UnixNano() - startTime.UnixNano()
		metricsCollectingTime.Observe(float64(took))
	}()

	for c := range collectors {
		cErr := collectors[c]()
		if cErr != nil {
			wasError = true
		}
	}
}
