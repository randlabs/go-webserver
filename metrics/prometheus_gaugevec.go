// See the LICENSE file for license details.

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"google.golang.org/protobuf/proto"
)

// -----------------------------------------------------------------------------

type gaugeVecWithCallbackCollector struct {
	desc    *prometheus.Desc
	metrics []prometheus.Metric
}

type gaugeVecWithCallbackMetric struct {
	self       prometheus.Metric
	desc       *prometheus.Desc
	labelPairs []*dto.LabelPair
	handler    ValueHandler
}

// -----------------------------------------------------------------------------

func (c *gaugeVecWithCallbackCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.desc
}

func (c *gaugeVecWithCallbackCollector) Collect(ch chan<- prometheus.Metric) {
	for _, v := range c.metrics {
		ch <- v
	}
}

func (v *gaugeVecWithCallbackMetric) Desc() *prometheus.Desc {
	return v.desc
}

func (v *gaugeVecWithCallbackMetric) Write(out *dto.Metric) error {
	out.Label = v.labelPairs
	out.Gauge = &dto.Gauge{
		Value: proto.Float64(v.handler()),
	}
	return nil
}
