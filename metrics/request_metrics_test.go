package metrics_test

import (
	"time"

	"code.cloudfoundry.org/clock/fakeclock"
	mfakes "code.cloudfoundry.org/diego-logging-client/testhelpers"
	"code.cloudfoundry.org/lager/lagertest"
	"code.cloudfoundry.org/locket/metrics"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/ginkgomon"
)

var _ = Describe("RequestMetrics", func() {
	type FakeGauge struct {
		Name  string
		Value int
	}

	var (
		runner           *metrics.RequestMetricsNotifier
		process          ifrit.Process
		fakeMetronClient *mfakes.FakeIngressClient
		logger           *lagertest.TestLogger
		fakeClock        *fakeclock.FakeClock
		metricsInterval  time.Duration
		metricsChan      chan FakeGauge
	)

	metricsChan = make(chan FakeGauge, 100)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("request-metrics")
		fakeMetronClient = new(mfakes.FakeIngressClient)
		fakeClock = fakeclock.NewFakeClock(time.Now())
		metricsInterval = 10 * time.Second

		fakeMetronClient.SendMetricStub = func(name string, value int) error {
			defer GinkgoRecover()

			Eventually(metricsChan).Should(BeSent(FakeGauge{name, value}))
			return nil
		}
		fakeMetronClient.SendDurationStub = func(name string, value time.Duration) error {
			defer GinkgoRecover()

			Eventually(metricsChan).Should(BeSent(FakeGauge{name, int(value)}))
			return nil
		}
	})

	JustBeforeEach(func() {
		runner = metrics.NewRequestMetricsNotifier(
			logger,
			fakeClock,
			fakeMetronClient,
			metricsInterval,
		)

		process = ifrit.Background(runner)
		Eventually(process.Ready()).Should(BeClosed())
	})

	AfterEach(func() {
		ginkgomon.Interrupt(process)
	})

	Context("when there is traffic to the locket server", func() {

		JustBeforeEach(func() {
			fakeClock.Increment(metricsInterval)
		})

		It("periodically emits the number of requests started", func() {
			Eventually(metricsChan).Should(Receive(Equal(FakeGauge{"RequestsStarted", 0})))
			runner.IncrementRequestsStartedCounter(1)
			fakeClock.Increment(metricsInterval)
			Eventually(metricsChan).Should(Receive(Equal(FakeGauge{"RequestsStarted", 1})))
		})

		It("periodically emits the number of requests succeeded", func() {
			Eventually(metricsChan).Should(Receive(Equal(FakeGauge{"RequestsSucceeded", 0})))
			runner.IncrementRequestsSucceededCounter(1)
			fakeClock.Increment(metricsInterval)
			Eventually(metricsChan).Should(Receive(Equal(FakeGauge{"RequestsSucceeded", 1})))
		})

		It("periodically emits the number of requests failed", func() {
			Eventually(metricsChan).Should(Receive(Equal(FakeGauge{"RequestsFailed", 0})))
			runner.IncrementRequestsFailedCounter(2)
			fakeClock.Increment(metricsInterval)
			Eventually(metricsChan).Should(Receive(Equal(FakeGauge{"RequestsFailed", 2})))
		})

		It("periodically emits the number of requests in flight", func() {
			Eventually(metricsChan).Should(Receive(Equal(FakeGauge{"RequestsInFlight", 0})))
			runner.IncrementRequestsInFlightCounter(4)
			fakeClock.Increment(metricsInterval)
			Eventually(metricsChan).Should(Receive(Equal(FakeGauge{"RequestsInFlight", 4})))

			runner.DecrementRequestsInFlightCounter(2)
			fakeClock.Increment(metricsInterval)
			Eventually(metricsChan).Should(Receive(Equal(FakeGauge{"RequestsInFlight", 2})))
		})

		It("periodically emits the max latency", func() {
			Eventually(metricsChan).Should(Receive(Equal(FakeGauge{"RequestLatencyMax", 0})))
			runner.UpdateLatency(5 * time.Millisecond)
			fakeClock.Increment(metricsInterval)
			Eventually(metricsChan).Should(Receive(Equal(FakeGauge{"RequestLatencyMax", int(5 * time.Millisecond)})))
			fakeClock.Increment(metricsInterval)
			Eventually(metricsChan).Should(Receive(Equal(FakeGauge{"RequestLatencyMax", 0})))
		})
	})
})
