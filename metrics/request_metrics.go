package metrics

import (
	"os"
	"sync"
	"time"

	"code.cloudfoundry.org/clock"
	loggingclient "code.cloudfoundry.org/diego-logging-client"
	loggregator "code.cloudfoundry.org/go-loggregator"
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

	requestsStarted       map[string]uint64
	requestsSucceeded     map[string]uint64
	requestsFailed        map[string]uint64
	requestsInFlight      map[string]uint64
	latencyMax            map[string]time.Duration
	latencyLock           *sync.Mutex
	requestsStartedLock   *sync.Mutex
	requestsSucceededLock *sync.Mutex
	requestsFailedLock    *sync.Mutex
	requestsInFlightLock  *sync.Mutex
}

func NewRequestMetricsNotifier(logger lager.Logger, ticker clock.Clock, metronClient loggingclient.IngressClient, metricsInterval time.Duration) *RequestMetricsNotifier {
	return &RequestMetricsNotifier{
		logger:                logger,
		ticker:                ticker,
		metricsInterval:       metricsInterval,
		metronClient:          metronClient,
		latencyLock:           &sync.Mutex{},
		requestsStartedLock:   &sync.Mutex{},
		requestsSucceededLock: &sync.Mutex{},
		requestsFailedLock:    &sync.Mutex{},
		requestsInFlightLock:  &sync.Mutex{},
		requestsStarted:       map[string]uint64{},
		requestsSucceeded:     map[string]uint64{},
		requestsFailed:        map[string]uint64{},
		requestsInFlight:      map[string]uint64{},
		latencyMax:            map[string]time.Duration{},
	}
}

func (notifier *RequestMetricsNotifier) IncrementRequestsStartedCounter(requestType string, delta int) {
	notifier.requestsStartedLock.Lock()
	defer notifier.requestsStartedLock.Unlock()

	notifier.requestsStarted[requestType] += uint64(delta)
}

func (notifier *RequestMetricsNotifier) IncrementRequestsSucceededCounter(requestType string, delta int) {
	notifier.requestsSucceededLock.Lock()
	defer notifier.requestsSucceededLock.Unlock()

	notifier.requestsSucceeded[requestType] += uint64(delta)
}

func (notifier *RequestMetricsNotifier) IncrementRequestsFailedCounter(requestType string, delta int) {
	notifier.requestsFailedLock.Lock()
	defer notifier.requestsFailedLock.Unlock()

	notifier.requestsFailed[requestType] += uint64(delta)
}

func (notifier *RequestMetricsNotifier) IncrementRequestsInFlightCounter(requestType string, delta int) {
	notifier.requestsInFlightLock.Lock()
	defer notifier.requestsInFlightLock.Unlock()

	notifier.requestsInFlight[requestType] += uint64(delta)
}

func (notifier *RequestMetricsNotifier) DecrementRequestsInFlightCounter(requestType string, delta int) {
	notifier.latencyLock.Lock()
	defer notifier.latencyLock.Unlock()

	notifier.requestsInFlight[requestType] -= uint64(delta)
}

func (notifier *RequestMetricsNotifier) UpdateLatency(requestType string, dur time.Duration) {
	notifier.latencyLock.Lock()
	defer notifier.latencyLock.Unlock()

	if dur > notifier.latencyMax[requestType] {
		notifier.latencyMax[requestType] = dur
	}
}

func (notifier *RequestMetricsNotifier) getRequestsStarted(requestType string) uint64 {
	notifier.requestsStartedLock.Lock()
	defer notifier.requestsStartedLock.Unlock()

	return notifier.requestsStarted[requestType]
}

func (notifier *RequestMetricsNotifier) getRequestsSucceeded(requestType string) uint64 {
	notifier.requestsSucceededLock.Lock()
	defer notifier.requestsSucceededLock.Unlock()

	return notifier.requestsSucceeded[requestType]
}

func (notifier *RequestMetricsNotifier) getRequestsFailed(requestType string) uint64 {
	notifier.requestsFailedLock.Lock()
	defer notifier.requestsFailedLock.Unlock()

	return notifier.requestsFailed[requestType]
}

func (notifier *RequestMetricsNotifier) getRequestsInFlight(requestType string) uint64 {
	notifier.requestsInFlightLock.Lock()
	defer notifier.requestsInFlightLock.Unlock()

	return notifier.requestsInFlight[requestType]
}

func (notifier *RequestMetricsNotifier) readAndResetMaxLatency(requestType string) time.Duration {
	notifier.latencyLock.Lock()
	defer notifier.latencyLock.Unlock()

	max := notifier.latencyMax[requestType]
	notifier.latencyMax[requestType] = 0

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

			var err error
			for requestType, _ := range notifier.requestsStarted {
				opt := loggregator.WithEnvelopeTag("request-type", requestType)
				value := notifier.getRequestsStarted(requestType)
				err = notifier.metronClient.SendMetric(requestsStartedMetric, int(value), opt)
				if err != nil {
					logger.Error("failed-to-emit-requests-started-metric", err)
				}
			}

			for requestType, _ := range notifier.requestsSucceeded {
				opt := loggregator.WithEnvelopeTag("request-type", requestType)
				value := notifier.getRequestsSucceeded(requestType)
				err = notifier.metronClient.SendMetric(requestsSucceededMetric, int(value), opt)
				if err != nil {
					logger.Error("failed-to-emit-requests-succeeded-metric", err)
				}
			}

			for requestType, _ := range notifier.requestsFailed {
				opt := loggregator.WithEnvelopeTag("request-type", requestType)
				value := notifier.getRequestsFailed(requestType)
				err = notifier.metronClient.SendMetric(requestsFailedMetric, int(value), opt)
				if err != nil {
					logger.Error("failed-to-emit-requests-failed-metric", err)
				}
			}

			for requestType, _ := range notifier.requestsInFlight {
				opt := loggregator.WithEnvelopeTag("request-type", requestType)
				value := notifier.getRequestsInFlight(requestType)
				err = notifier.metronClient.SendMetric(requestsInFlightMetric, int(value), opt)
				if err != nil {
					logger.Error("failed-to-emit-requests-in-flight-metric", err)
				}
			}

			for requestType, _ := range notifier.latencyMax {
				opt := loggregator.WithEnvelopeTag("request-type", requestType)
				max := notifier.readAndResetMaxLatency(requestType)
				err = notifier.metronClient.SendDuration(requestLatencyMaxDuration, max, opt)
				if err != nil {
					logger.Error("failed-to-emit-requests-latency-max-metric", err)
				}
			}
		}
	}

	return nil
}
