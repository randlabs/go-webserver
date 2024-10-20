// See the LICENSE file for license details.

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"google.golang.org/protobuf/proto"
)

// -----------------------------------------------------------------------------

type counterVecWithCallbackCollector struct {
	desc    *prometheus.Desc
	metrics []prometheus.Metric
}

type counterVecWithCallbackMetric struct {
	self       prometheus.Metric
	desc       *prometheus.Desc
	labelPairs []*dto.LabelPair
	handler    ValueHandler
}

// -----------------------------------------------------------------------------

func (c *counterVecWithCallbackCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.desc
}

func (c *counterVecWithCallbackCollector) Collect(ch chan<- prometheus.Metric) {
	for _, v := range c.metrics {
		ch <- v
	}
}

func (v *counterVecWithCallbackMetric) Desc() *prometheus.Desc {
	return v.desc
}

func (v *counterVecWithCallbackMetric) Write(out *dto.Metric) error {
	out.Label = v.labelPairs
	out.Counter = &dto.Counter{
		Value:    proto.Float64(v.handler()),
		Exemplar: nil,
	}
	return nil
}
