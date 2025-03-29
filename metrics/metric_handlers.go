// See the LICENSE file for license details.

package metrics

import (
	"strconv"

	webserver "github.com/mxmauro/go-webserver/v2"
	"github.com/mxmauro/go-webserver/v2/util"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// -----------------------------------------------------------------------------

func (mws *Controller) getHealthHandler() webserver.HandlerFunc {
	return func(req *webserver.RequestContext) error {
		// Get current health status from callback
		status := mws.healthCallback()

		// Send output
		respHdrs := req.ResponseHeaders()
		if isJSON(status) {
			respHdrs.SetBytesKV(util.HeaderContentType, util.ContentTypeApplicationJSON)
		} else {
			respHdrs.SetBytesKV(util.HeaderContentType, util.ContentTypeTextPlain)
		}

		if !req.IsHead() {
			_, _ = req.WriteString(status)
		} else {
			respHdrs.SetBytesK(util.HeaderContentLength, strconv.FormatUint(uint64(int64(len(status))), 10))
		}

		// Done
		req.Success()
		return nil
	}
}

func (mws *Controller) getMetricsHandler() webserver.HandlerFunc {
	return webserver.NewHandlerFromHttpHandler(promhttp.HandlerFor(
		mws.registry,
		promhttp.HandlerOpts{
			EnableOpenMetrics:   true,
			MaxRequestsInFlight: 5,
		},
	))
}
