package api

import (
	"strconv"

	"github.com/Trendyol/go-dcp-client/stream"

	"github.com/Trendyol/go-dcp-client/couchbase"

	"github.com/Trendyol/go-dcp-client/logger"

	"github.com/Trendyol/go-dcp-client/helpers"
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"

	"github.com/prometheus/client_golang/prometheus"
)

type metricCollector struct {
	stream           stream.Stream
	client           couchbase.Client
	vBucketDiscovery stream.VBucketDiscovery

	mutation   *prometheus.Desc
	deletion   *prometheus.Desc
	expiration *prometheus.Desc

	currentSeqNo *prometheus.Desc
	startSeqNo   *prometheus.Desc
	endSeqNo     *prometheus.Desc

	averageProcessMs *prometheus.Desc
	rebalanceCount   *prometheus.Desc

	lag *prometheus.Desc

	totalMembers      *prometheus.Desc
	memberNumber      *prometheus.Desc
	membershipType    *prometheus.Desc
	vBucketCount      *prometheus.Desc
	vBucketRangeStart *prometheus.Desc
	vBucketRangeEnd   *prometheus.Desc
}

func (s *metricCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(s, ch)
}

//nolint:funlen
func (s *metricCollector) Collect(ch chan<- prometheus.Metric) {
	seqNoMap, err := s.client.GetVBucketSeqNos()

	observer := s.stream.GetObserver()

	observer.LockMetrics()

	for vbID, metric := range observer.GetMetrics() {
		ch <- prometheus.MustNewConstMetric(
			s.mutation,
			prometheus.CounterValue,
			metric.TotalMutations,
			strconv.Itoa(int(vbID)),
		)

		ch <- prometheus.MustNewConstMetric(
			s.deletion,
			prometheus.CounterValue,
			metric.TotalDeletions,
			strconv.Itoa(int(vbID)),
		)

		ch <- prometheus.MustNewConstMetric(
			s.expiration,
			prometheus.CounterValue,
			metric.TotalExpirations,
			strconv.Itoa(int(vbID)),
		)
	}

	observer.UnlockMetrics()

	s.stream.LockOffsets()

	offsets, _, _ := s.stream.GetOffsets()

	for vbID, offset := range offsets {
		ch <- prometheus.MustNewConstMetric(
			s.currentSeqNo,
			prometheus.GaugeValue,
			float64(offset.SeqNo),
			strconv.Itoa(int(vbID)),
		)

		ch <- prometheus.MustNewConstMetric(
			s.startSeqNo,
			prometheus.GaugeValue,
			float64(offset.StartSeqNo),
			strconv.Itoa(int(vbID)),
		)

		ch <- prometheus.MustNewConstMetric(
			s.endSeqNo,
			prometheus.GaugeValue,
			float64(offset.EndSeqNo),
			strconv.Itoa(int(vbID)),
		)

		var lag float64

		if seqNoMap[vbID] > offset.SeqNo {
			lag = float64(seqNoMap[vbID] - offset.SeqNo)
		}

		if err != nil {
			ch <- prometheus.NewInvalidMetric(
				s.lag,
				err,
			)
		} else {
			ch <- prometheus.MustNewConstMetric(
				s.lag,
				prometheus.GaugeValue,
				lag,
				strconv.Itoa(int(vbID)),
			)
		}
	}

	s.stream.UnlockOffsets()

	streamMetric := s.stream.GetMetric()

	ch <- prometheus.MustNewConstMetric(
		s.averageProcessMs,
		prometheus.CounterValue,
		streamMetric.AverageProcessMs.Value(),
		[]string{}...,
	)

	ch <- prometheus.MustNewConstMetric(
		s.rebalanceCount,
		prometheus.CounterValue,
		float64(streamMetric.RebalanceCount),
		[]string{}...,
	)

	vBucketDiscoveryMetric := s.vBucketDiscovery.GetMetric()

	ch <- prometheus.MustNewConstMetric(
		s.totalMembers,
		prometheus.GaugeValue,
		float64(vBucketDiscoveryMetric.TotalMembers),
		[]string{}...,
	)

	ch <- prometheus.MustNewConstMetric(
		s.memberNumber,
		prometheus.GaugeValue,
		float64(vBucketDiscoveryMetric.MemberNumber),
		[]string{}...,
	)

	ch <- prometheus.MustNewConstMetric(
		s.membershipType,
		prometheus.GaugeValue,
		0,
		[]string{vBucketDiscoveryMetric.Type}...,
	)

	ch <- prometheus.MustNewConstMetric(
		s.vBucketCount,
		prometheus.GaugeValue,
		float64(vBucketDiscoveryMetric.VBucketCount),
		[]string{}...,
	)

	ch <- prometheus.MustNewConstMetric(
		s.vBucketRangeStart,
		prometheus.GaugeValue,
		float64(vBucketDiscoveryMetric.VBucketRangeStart),
		[]string{}...,
	)

	ch <- prometheus.MustNewConstMetric(
		s.vBucketRangeEnd,
		prometheus.GaugeValue,
		float64(vBucketDiscoveryMetric.VBucketRangeEnd),
		[]string{}...,
	)
}

//nolint:funlen
func newMetricCollector(client couchbase.Client, stream stream.Stream, vBucketDiscovery stream.VBucketDiscovery) *metricCollector {
	return &metricCollector{
		stream:           stream,
		client:           client,
		vBucketDiscovery: vBucketDiscovery,

		mutation: prometheus.NewDesc(
			prometheus.BuildFQName(helpers.Name, "mutation", "total"),
			"Mutation count",
			[]string{"vbId"},
			nil,
		),
		deletion: prometheus.NewDesc(
			prometheus.BuildFQName(helpers.Name, "deletion", "total"),
			"Deletion count",
			[]string{"vbId"},
			nil,
		),
		expiration: prometheus.NewDesc(
			prometheus.BuildFQName(helpers.Name, "expiration", "total"),
			"Expiration count",
			[]string{"vbId"},
			nil,
		),
		currentSeqNo: prometheus.NewDesc(
			prometheus.BuildFQName(helpers.Name, "seq_no", "current"),
			"Current seq no",
			[]string{"vbId"},
			nil,
		),
		startSeqNo: prometheus.NewDesc(
			prometheus.BuildFQName(helpers.Name, "start_seq_no", "current"),
			"Start seq no",
			[]string{"vbId"},
			nil,
		),
		endSeqNo: prometheus.NewDesc(
			prometheus.BuildFQName(helpers.Name, "end_seq_no", "current"),
			"End seq no",
			[]string{"vbId"},
			nil,
		),
		lag: prometheus.NewDesc(
			prometheus.BuildFQName(helpers.Name, "lag", "current"),
			"Lag",
			[]string{"vbId"},
			nil,
		),
		averageProcessMs: prometheus.NewDesc(
			prometheus.BuildFQName(helpers.Name, "average_process_ms", "current"),
			"Average process ms at 10sec windows",
			[]string{},
			nil,
		),
		rebalanceCount: prometheus.NewDesc(
			prometheus.BuildFQName(helpers.Name, "rebalance_count", "current"),
			"Rebalance count",
			[]string{},
			nil,
		),
		totalMembers: prometheus.NewDesc(
			prometheus.BuildFQName(helpers.Name, "total_members", "current"),
			"Total members",
			[]string{},
			nil,
		),
		memberNumber: prometheus.NewDesc(
			prometheus.BuildFQName(helpers.Name, "member_number", "current"),
			"Member number",
			[]string{},
			nil,
		),
		membershipType: prometheus.NewDesc(
			prometheus.BuildFQName(helpers.Name, "membership_type", "current"),
			"Membership type",
			[]string{"type"},
			nil,
		),
		vBucketCount: prometheus.NewDesc(
			prometheus.BuildFQName(helpers.Name, "vbucket_count", "current"),
			"VBucket count",
			[]string{},
			nil,
		),
		vBucketRangeStart: prometheus.NewDesc(
			prometheus.BuildFQName(helpers.Name, "vbucket_range_start", "current"),
			"VBucket range start",
			[]string{},
			nil,
		),
		vBucketRangeEnd: prometheus.NewDesc(
			prometheus.BuildFQName(helpers.Name, "vbucket_range_end", "current"),
			"VBucket range end",
			[]string{},
			nil,
		),
	}
}

func NewMetricMiddleware(app *fiber.App,
	config *helpers.Config,
	stream stream.Stream,
	client couchbase.Client,
	vBucketDiscovery stream.VBucketDiscovery,
) (func(ctx *fiber.Ctx) error, error) {
	err := prometheus.DefaultRegisterer.Register(newMetricCollector(client, stream, vBucketDiscovery))
	if err != nil {
		return nil, err
	}

	fiberPrometheus := fiberprometheus.New(config.Dcp.Group.Name)
	fiberPrometheus.RegisterAt(app, config.Metric.Path)

	logger.Log.Printf("metric middleware registered on path %s", config.Metric.Path)

	return fiberPrometheus.Middleware, nil
}