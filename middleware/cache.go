package middleware

import (
	"fmt"
	"strings"

	webserver "github.com/randlabs/go-webserver"
	"github.com/randlabs/go-webserver/request"
)

// -----------------------------------------------------------------------------

// CacheControlOptions defines the behavior on how Cache-Control headers are sent.
type CacheControlOptions struct {
	Public                        bool
	Private                       bool
	NoCache                       bool
	NoStore                       bool
	NoTransform                   bool
	MustRevalidate                bool
	ProxyRevalidate               bool
	MaxAgeInSeconds               *uint32
	SharedMaxAgeInSeconds         *uint32
	StaleWhileRevalidateInSeconds *uint32
	StaleIfErrorInSeconds         *uint32
}

// -----------------------------------------------------------------------------

// DisableClientCache creates a default cache control middleware that disables the client's cache
func DisableClientCache() webserver.MiddlewareFunc {
	var zero uint32

	return NewCacheControl(CacheControlOptions{
		Private:               true,
		NoCache:               true,
		NoStore:               true,
		MustRevalidate:        true,
		ProxyRevalidate:       true,
		MaxAgeInSeconds:       &zero,
		SharedMaxAgeInSeconds: &zero,
	})
}

// NewCacheControl creates a new client cache control middleware based on the specified options
func NewCacheControl(opts CacheControlOptions) webserver.MiddlewareFunc {
	cacheValue := make([]string, 0)

	if opts.Public {
		cacheValue = append(cacheValue, "public")
	}
	if opts.Private {
		cacheValue = append(cacheValue, "private")
	}

	if opts.NoCache {
		cacheValue = append(cacheValue, "no-cache")
	}
	if opts.NoStore {
		cacheValue = append(cacheValue, "no-store")
	}
	if opts.NoTransform {
		cacheValue = append(cacheValue, "no-transform")
	}

	if opts.MustRevalidate {
		cacheValue = append(cacheValue, "must-revalidate")
	}
	if opts.ProxyRevalidate {
		cacheValue = append(cacheValue, "proxy-revalidate")
	}

	if opts.MaxAgeInSeconds != nil {
		cacheValue = append(cacheValue, fmt.Sprintf("max-age=%v", *opts.MaxAgeInSeconds))
	}
	if opts.SharedMaxAgeInSeconds != nil {
		cacheValue = append(cacheValue, fmt.Sprintf("s-maxage=%v", *opts.SharedMaxAgeInSeconds))
	}
	if opts.StaleWhileRevalidateInSeconds != nil {
		cacheValue = append(cacheValue, fmt.Sprintf("stale-if-error=%v", *opts.StaleWhileRevalidateInSeconds))
	}
	if opts.StaleIfErrorInSeconds != nil {
		cacheValue = append(cacheValue, fmt.Sprintf("stale-while-revalidate=%v", *opts.StaleIfErrorInSeconds))
	}

	finalCacheValue := strings.Join(cacheValue, ",")

	return func(next webserver.HandlerFunc) webserver.HandlerFunc {
		return func(req *request.RequestContext) error {
			// Set cache control header
			req.SetResponseHeader("Cache-Control", finalCacheValue)

			// Done
			return next(req)
		}
	}
}

func CacheControlMaxAge(seconds uint32) *uint32 {
	return &seconds
}
