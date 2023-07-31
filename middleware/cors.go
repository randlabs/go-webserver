package middleware

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	webserver "github.com/randlabs/go-webserver"
	"github.com/randlabs/go-webserver/request"
	"github.com/randlabs/go-webserver/util"
)

// -----------------------------------------------------------------------------

// CORSOptions defines the behavior on how CORS requests should be handled.
type CORSOptions struct {
	// AllowOrigins defines a list of origins that may access the resource.
	// Optional. Defaults to "*".
	AllowOrigins []string `json:"allow-origins,omitempty"`

	// AllowMethods defines a list methods allowed when accessing the resource.
	// If defined as an empty list, the preflight `Allow` request header value will be used.
	AllowMethods []string `json:"allow-methods,omitempty"`

	// AllowHeaders defines a list of request headers that can be used when
	// making the actual request.
	AllowHeaders []string `json:"allow-headers,omitempty"`

	// AllowCredentials indicates whether the response to the request
	// can be exposed when the credentials flag is true.
	// Do not set to true if allow origins is "*".
	// See: http://blog.portswigger.net/2016/10/exploiting-cors-misconfigurations-for.html
	AllowCredentials bool `json:"allow-credentials,omitempty"`

	// ExposeHeaders defines a whitelist headers that clients are allowed to access.
	ExposeHeaders []string `json:"expose-headers,omitempty"`

	// MaxAge indicates how many seconds the results of a preflight request can be cached. Defaults to 0.
	MaxAge int `json:"max-age,omitempty"`
}

// -----------------------------------------------------------------------------

// DefaultCORS creates a default CORS middleware that allows requests from anywhere
func DefaultCORS() webserver.MiddlewareFunc {
	return NewCORS(CORSOptions{})
}

// NewCORS creates a new CORS middleware based on the specified options
func NewCORS(opts CORSOptions) webserver.MiddlewareFunc {
	var allowOrigins []string
	var allowOriginPatterns []string

	// Parse options
	hasWildCardOrigin := true
	if len(opts.AllowOrigins) > 0 {
		allowOrigins = make([]string, 0)
		allowOriginPatterns = make([]string, 0)
		hasWildCardOrigin = false

		for _, allowOrigin := range opts.AllowOrigins {
			if allowOrigin == "*" {
				hasWildCardOrigin = true
			} else {
				allowOrigins = append(allowOrigins, allowOrigin)

				pattern := regexp.QuoteMeta(allowOrigin)
				pattern = strings.Replace(pattern, "\\*", ".*", -1)
				pattern = strings.Replace(pattern, "\\?", ".", -1)
				pattern = "^" + pattern + "$"
				allowOriginPatterns = append(allowOriginPatterns, pattern)
			}
		}
	}

	hasCustomAllowMethods := true
	var allowMethods string
	if len(opts.AllowMethods) > 0 {
		allowMethods = strings.Join(opts.AllowMethods, ",")
	} else {
		hasCustomAllowMethods = false
		allowMethods = http.MethodGet + "," + http.MethodHead + "," + http.MethodPut + "," + http.MethodPatch +
			"," + http.MethodPost + "," + http.MethodDelete
	}
	allowHeaders := strings.Join(opts.AllowHeaders, ",")
	exposeHeaders := strings.Join(opts.ExposeHeaders, ",")
	maxAge := strconv.Itoa(opts.MaxAge)

	// Setup middleware function
	return func(next webserver.HandlerFunc) webserver.HandlerFunc {
		return func(req *request.RequestContext) error {
			origin := req.RequestHeader("origin")
			allowOrigin := ""

			if len(origin) > 0 {
				req.AddResponseHeader("Vary", origin)
			}

			// See https://developer.mozilla.org/en-US/docs/Glossary/Preflight_request for details on how to handle
			// preflight request.
			preflight := req.IsOptions()

			// If the router set an allow methods
			routerAllowMethods := ""
			if preflight {
				var ok bool

				routerAllowMethods, ok = req.UserValue("routerAllow").(string)
				if ok && len(routerAllowMethods) > 0 {
					req.SetResponseHeader("Allow", routerAllowMethods)
				} else {
					routerAllowMethods = ""
				}
			}

			// No origin provided?
			if len(origin) == 0 {
				if preflight {
					return req.NoContent(http.StatusNoContent)
				}
				// Go to next middleware
				return next(req)
			}

			// Check allowed origins
			if hasWildCardOrigin {
				if opts.AllowCredentials {
					allowOrigin = origin
				} else {
					allowOrigin = "*"
				}
			}

			if len(allowOrigin) == 0 {
				for _, o := range allowOrigins {
					if util.DoesSubdomainMatch(origin, o) {
						allowOrigin = origin
						break
					}
				}
			}

			if len(allowOrigin) == 0 && len(origin) <= (253+3+5) && strings.Contains(origin, "://") {
				for _, re := range allowOriginPatterns {
					if match, _ := regexp.MatchString(re, origin); match {
						allowOrigin = origin
						break
					}
				}
			}

			// Origin not allowed
			if len(allowOrigin) == 0 {
				if preflight {
					return req.NoContent(http.StatusNoContent)
				}
				// Go to next middleware
				return next(req)
			}

			req.SetResponseHeader("Access-Control-Allow-Origin", allowOrigin)
			if opts.AllowCredentials {
				req.SetResponseHeader("Access-Control-Allow-Credentials", "true")
			}

			// Simple request
			if !preflight {
				if len(exposeHeaders) > 0 {
					req.SetResponseHeader("Access-Control-Expose-Headers", exposeHeaders)
				}
				// Go to next middleware
				return next(req)
			}

			// Preflight request
			req.AddResponseHeader("Vary", "Access-Control-Request-Method")
			req.AddResponseHeader("Vary", "Access-Control-Request-Headers")

			if !hasCustomAllowMethods && len(routerAllowMethods) > 0 {
				req.SetResponseHeader("Access-Control-Allow-Methods", routerAllowMethods)
			} else {
				req.SetResponseHeader("Access-Control-Allow-Methods", allowMethods)
			}

			if len(allowHeaders) > 0 {
				req.SetResponseHeader("Access-Control-Allow-Headers", allowHeaders)
			} else {
				header := req.RequestHeader("Access-Control-Request-Headers")
				if len(header) > 0 {
					req.SetResponseHeader("Access-Control-Allow-Headers", header)
				}
			}
			if len(maxAge) > 0 {
				req.SetResponseHeader("Access-Control-Max-Age", maxAge)
			}

			return req.NoContent(http.StatusNoContent)
		}
	}
}
