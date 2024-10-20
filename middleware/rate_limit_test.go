// See the LICENSE file for license details.

package middleware_test

import (
	"strconv"
	"testing"
	"time"

	webserver "github.com/mxmauro/go-webserver/v2"
	"github.com/mxmauro/go-webserver/v2/internal/testcommon"
	"github.com/mxmauro/go-webserver/v2/middleware"
)

// -----------------------------------------------------------------------------

func TestMiddlewareRateLimiter(t *testing.T) {
	//Create server
	srv := testcommon.RunWebServer(t, func(srv *webserver.Server) error {
		// Add some middlewares
		srv.Use(middleware.NewRateLimiter(middleware.RateLimiterOptions{
			Max:                5,
			Expiration:         1 * time.Second,
			KeyGenerator:       nil,
			LimitReached:       nil,
			SkipFailedRequests: false,
			ExternalStorage:    nil,
			MaxMemoryCacheSize: 0,
		}))

		// Done
		return nil
	})
	defer srv.Stop()

	for count := 1; count <= 5; count++ {
		_, headers, err := testcommon.QueryApiVersion(false, nil, nil, []int{200})
		if err != nil {
			t.Fatalf("unable to query api [%v]", err)
		}
		rateLimitLimit := headers.Get("X-Rate-Limit-Limit")
		if rateLimitLimit != "5" {
			t.Fatalf("unexpected X-Rate-Limit-Limit [got:%v / expected:%v]", rateLimitLimit, 5)
		}
		rateLimitRemaining := headers.Get("X-Rate-Limit-Remaining")
		expected := strconv.Itoa(5 - count)
		if rateLimitRemaining != expected {
			t.Fatalf("unexpected X-Rate-Limit-Remaining [got:%v / expected:%v]", rateLimitRemaining, expected)
		}
		// rateLimitReset := headers.Get("X-Rate-Limit-Reset")
	}

	statusCode, headers, err := testcommon.QueryApiVersion(false, nil, nil, []int{429})
	if err != nil {
		if statusCode == 0 {
			t.Fatalf("unable to query api with wrong header [%v]", err)
		} else {
			t.Fatalf("unexpected status code while querying api with wrong header [%d]", statusCode)
		}
	}
	retryAfter := headers.Get("Retry-After")
	if len(retryAfter) == 0 {
		t.Fatalf("missing Retry-After header")
	}
}
