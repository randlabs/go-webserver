package middleware_test

import (
	"net/http"
	"testing"

	webserver "github.com/randlabs/go-webserver/v2"
	"github.com/randlabs/go-webserver/v2/helpers_test"
	"github.com/randlabs/go-webserver/v2/middleware"
)

// -----------------------------------------------------------------------------

func TestEtag(t *testing.T) {
	var statusCode int

	//Create server
	srv := helpers_test.RunWebServer(t, func(srv *webserver.Server) error {
		// Add some middlewares
		srv.Use(middleware.NewETag(false))

		// Done
		return nil
	})
	defer srv.Stop()

	// Query api for the first time
	_, headers, err := helpers_test.QueryApiVersion(false, nil, nil, []int{200})
	if err != nil {
		t.Fatalf("unable to query api [%v]", err)
	}
	etag := headers.Get(http.CanonicalHeaderKey("ETag"))
	if len(etag) == 0 {
		t.Fatalf("ETag header not found")
	}

	// Query it again and expect not modified
	statusCode, _, err = helpers_test.QueryApiVersion(false, nil, http.Header{
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
