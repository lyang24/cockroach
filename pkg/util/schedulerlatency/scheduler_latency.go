package schedulerlatency

import (
	"github.com/cockroachdb/cockroach/pkg/util/ring"
	"github.com/cockroachdb/cockroach/pkg/util/syncutil"
	"runtime/metrics"
	"time"
)

type SchedulerLatencyMetric struct {
	listener LatencyObserver
	mu       struct {
		syncutil.Mutex
		ringBuffer            ring.Buffer[*metrics.Float64Histogram]
		lastIntervalHistogram *metrics.Float64Histogram
	}
}

func (s *SchedulerLatencyMetric) lastIntervalHistogram() *metrics.Float64Histogram {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.mu.lastIntervalHistogram
}

func (s *SchedulerLatencyMetric) setPeriodAndDuration(period, duration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.mu.ringBuffer.Discard()
	numSamples := int(duration / period)
	if numSamples < 1 {
		numSamples = 1 // we need at least one sample to compare (also safeguards against integer division)
	}
	s.mu.ringBuffer.Resize(numSamples)
	s.mu.lastIntervalHistogram = nil
}
