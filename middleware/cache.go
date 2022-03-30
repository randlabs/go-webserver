package middleware

import (
	"fmt"
	"github.com/randlabs/go-webserver/models"
	"github.com/randlabs/go-webserver/request"
	"strings"
)

// -----------------------------------------------------------------------------

type CacheControlOptions struct {
	Public                         bool
	Private                        bool
	NoCache                        bool
	NoStore                        bool
	NoTransform                    bool
	MustRevalidate                 bool
	ProxyRevalidate                bool
	MaxAgeInSeconds                *uint32
	SharedMaxAgeInSeconds          *uint32
	StaleWhileRevalidateInSeconds  *uint32
	StaleIfErrorInSeconds          *uint32
}

// -----------------------------------------------------------------------------

func DisableCacheControl() MiddlewareFunc {
	var zero uint32

	return NewCacheControl(CacheControlOptions{
		Private:                true,
		NoCache:                true,
		NoStore:                true,
		MustRevalidate:         true,
		ProxyRevalidate:        true,
		MaxAgeInSeconds:        &zero,
		SharedMaxAgeInSeconds:  &zero,
	})
}

func NewCacheControl(opts CacheControlOptions) MiddlewareFunc {
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

	return func(next models.HandlerFunc) models.HandlerFunc {
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
