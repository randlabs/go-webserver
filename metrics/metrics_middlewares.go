// See the LICENSE file for license details.

package metrics

import (
	"net/http"

	webserver "github.com/mxmauro/go-webserver/v2"
)

// -----------------------------------------------------------------------------

func (mws *Controller) createAliveMiddleware() webserver.HandlerFunc {
	return func(req *webserver.RequestContext) error {
		// Process the request if we are not shutting down
		if !mws.rp.Acquire() {
			req.Error(http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return nil
		}
		defer mws.rp.Release()

		return req.Next()
	}
}
