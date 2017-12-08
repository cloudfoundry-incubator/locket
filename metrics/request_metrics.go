package metrics

import (
	"os"
	"sync"
	"sync/atomic"
	"time"

	"code.cloudfoundry.org/clock"
	loggingclient "code.cloudfoundry.org/diego-logging-client"
	"code.cloudfoundry.org/lager"
)

const (
	requestsStartedMetric        = "RequestsStarted"
	requestsSucceededMetric      = "RequestsSucceeded"
	requestsFailedMetric         = "RequestsFailed"
	requestsInFlightMetric       = "RequestsInFlight"
	requestLatencyMaxDuration    = "RequestLatencyMax"
	requestLatencyMedianDuration = "RequestLatencyMedian"

	latencyBufferCapacity = 1024
)

//go:generate counterfeiter . RequestMetrics
type RequestMetrics interface {
	IncrementRequestsStartedCounter(delta int)
	IncrementRequestsSucceededCounter(delta int)
	IncrementRequestsFailedCounter(delta int)
	IncrementRequestsInFlightCounter(delta int)
	DecrementRequestsInFlightCounter(delta int)
	UpdateLatency(dur time.Duration)

	GetRequestsInFlight() uint64
}

type RequestMetricsNotifier struct {
	logger          lager.Logger
	ticker          clock.Clock
	metricsInterval time.Duration
	metronClient    loggingclient.IngressClient

	requestsStarted   uint64
	requestsSucceeded uint64
	requestsFailed    uint64
	requestsInFlight  uint64
	latencyMax        time.Duration
	latencyLock       *sync.Mutex
}

func NewRequestMetricsNotifier(logger lager.Logger, ticker clock.Clock, metronClient loggingclient.IngressClient, metricsInterval time.Duration) *RequestMetricsNotifier {
	return &RequestMetricsNotifier{
		logger:          logger,
		ticker:          ticker,
		metricsInterval: metricsInterval,
		metronClient:    metronClient,
		latencyLock:     &sync.Mutex{},
	}
}

func (notifier *RequestMetricsNotifier) IncrementRequestsStartedCounter(delta int) {
	atomic.AddUint64(&notifier.requestsStarted, uint64(delta))
}

func (notifier *RequestMetricsNotifier) IncrementRequestsSucceededCounter(delta int) {
	atomic.AddUint64(&notifier.requestsSucceeded, uint64(delta))
}

func (notifier *RequestMetricsNotifier) IncrementRequestsFailedCounter(delta int) {
	atomic.AddUint64(&notifier.requestsFailed, uint64(delta))
}

func (notifier *RequestMetricsNotifier) IncrementRequestsInFlightCounter(delta int) {
	atomic.AddUint64(&notifier.requestsInFlight, uint64(delta))
}

func (notifier *RequestMetricsNotifier) DecrementRequestsInFlightCounter(delta int) {
	atomic.AddUint64(&notifier.requestsInFlight, ^uint64(delta-1))
}

func (notifier *RequestMetricsNotifier) UpdateLatency(dur time.Duration) {
	notifier.latencyLock.Lock()
	defer notifier.latencyLock.Unlock()

	if dur > notifier.latencyMax {
		notifier.latencyMax = dur
	}
}

func (notifier *RequestMetricsNotifier) getRequestsStarted() uint64 {
	return atomic.LoadUint64(&notifier.requestsStarted)
}

func (notifier *RequestMetricsNotifier) getRequestsSucceeded() uint64 {
	return atomic.LoadUint64(&notifier.requestsSucceeded)
}

func (notifier *RequestMetricsNotifier) getRequestsFailed() uint64 {
	return atomic.LoadUint64(&notifier.requestsFailed)
}

func (notifier *RequestMetricsNotifier) GetRequestsInFlight() uint64 {
	return atomic.LoadUint64(&notifier.requestsInFlight)
}

func (notifier *RequestMetricsNotifier) readAndResetMaxLatency() time.Duration {
	notifier.latencyLock.Lock()
	defer notifier.latencyLock.Unlock()

	max := notifier.latencyMax
	notifier.latencyMax = 0

	return max
}

func (notifier *RequestMetricsNotifier) Run(signals <-chan os.Signal, ready chan<- struct{}) error {
	logger := notifier.logger.Session("request-metrics-notifier")
	logger.Info("starting", lager.Data{"interval": notifier.metricsInterval})
	defer logger.Info("completed")
	close(ready)

	tick := notifier.ticker.NewTicker(notifier.metricsInterval)

	for {
		select {
		case <-signals:
			return nil
		case <-tick.C():
			logger.Debug("emitting-metrics")

			err := notifier.metronClient.SendMetric(requestsStartedMetric, int(notifier.getRequestsStarted()))
			if err != nil {
				logger.Error("failed-to-emit-requests-started-metric", err)
			}

			err = notifier.metronClient.SendMetric(requestsSucceededMetric, int(notifier.getRequestsSucceeded()))
			if err != nil {
				logger.Error("failed-to-emit-requests-succeeded-metric", err)
			}

			err = notifier.metronClient.SendMetric(requestsFailedMetric, int(notifier.getRequestsFailed()))
			if err != nil {
				logger.Error("failed-to-emit-requests-failed-metric", err)
			}

			err = notifier.metronClient.SendMetric(requestsInFlightMetric, int(notifier.GetRequestsInFlight()))
			if err != nil {
				logger.Error("failed-to-emit-requests-in-flight-metric", err)
			}

			max := notifier.readAndResetMaxLatency()
			err = notifier.metronClient.SendDuration(requestLatencyMaxDuration, max)
			if err != nil {
				logger.Error("failed-to-emit-requests-latency-max-metric", err)
			}
		}
	}

	return nil
}
