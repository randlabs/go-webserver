package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

// -----------------------------------------------------------------------------

func (mws *Controller) createPrometheusRegistry() error {
	// Create registry
	registry := prometheus.NewRegistry()

	// Add Golang specific collectors
	err := registry.Register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	if err != nil {
		return err
	}
	err = registry.Register(collectors.NewGoCollector())
	if err != nil {
		return err
	}

	// Done
	mws.registry = registry
	return nil
}
