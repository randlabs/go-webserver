// See the LICENSE file for license details.

package middleware_test

import (
	"net/http"
	"testing"

	webserver "github.com/mxmauro/go-webserver/v2"
	"github.com/mxmauro/go-webserver/v2/internal/testcommon"
	"github.com/mxmauro/go-webserver/v2/middleware"
)

// -----------------------------------------------------------------------------

func TestMiddlewareEtag(t *testing.T) {
	var statusCode int

	//Create server
	srv := testcommon.RunWebServer(t, func(srv *webserver.Server) error {
		// Add some middlewares
		srv.Use(middleware.NewETag(false))

		// Done
		return nil
	})
	defer srv.Stop()

	// Query api for the first time
	_, headers, err := testcommon.QueryApiVersion(false, nil, nil, []int{200})
	if err != nil {
		t.Fatalf("unable to query api [%v]", err)
	}
	etag := headers.Get(http.CanonicalHeaderKey("ETag"))
	if len(etag) == 0 {
		t.Fatalf("ETag header not found")
	}

	// Query it again and expect not modified
	statusCode, _, err = testcommon.QueryApiVersion(false, nil, http.Header{
		"If-None-Match": []string{etag},
	}, []int{304})

	if err != nil {
		if statusCode == 0 {
			t.Fatalf("unable to query api [%v]", err)
		} else {
			t.Fatalf("unexpected status code while querying api [%d]", statusCode)
		}
	}
}
