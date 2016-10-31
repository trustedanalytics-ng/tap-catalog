/**
 * Copyright (c) 2016 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"

	mutils "github.com/trustedanalytics/tap-metrics/utils"
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
