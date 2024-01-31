// Copyright 2024 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package schedulerlatency

import (
	"fmt"
	"runtime/metrics"
	"time"
)

type metricState interface {
	setPeriodAndDuration(period, duration time.Duration)
	processUpdate(value metrics.Value)
}

type runtimeMetrics struct {
	samples          []metrics.Sample
	namedMetricState map[string]metricState
}

func initRuntimeMetrics(namedMetricState map[string]metricState) runtimeMetrics {
	samples := []metrics.Sample{}
	for name, _ := range namedMetricState {
		samples = append(samples, metrics.Sample{
			Name: name,
		})
	}
	return runtimeMetrics{
		samples,
		namedMetricState,
	}
}

func (r runtimeMetrics) SampleRuntimeMetric() {
	metrics.Read(r.samples)
	for _, s := range r.samples {
		ms, exists := r.namedMetricState[s.Name]
		if !exists {
			panic(fmt.Sprintf("uninitialized go runtime metric: %s", s.Name))
		}
		ms.processUpdate(s.Value)
	}
}
