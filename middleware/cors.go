package middleware

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	webserver "github.com/randlabs/go-webserver/v2"
	"github.com/randlabs/go-webserver/v2/util"
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
func DefaultCORS() webserver.HandlerFunc {
	return NewCORS(CORSOptions{})
}

// NewCORS creates a new CORS middleware based on the specified options
func NewCORS(opts CORSOptions) webserver.HandlerFunc {
	var allowMethods []byte
	var allowOrigins []string
	var allowOriginPatterns []string
	var exposeHeaders []byte
	var allowHeaders []byte
	var maxAge []byte

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

	if len(opts.AllowMethods) > 0 {
		allowMethods = util.UnsafeString2ByteSlice(strings.Join(opts.AllowMethods, ","))
	} else {
		allowMethods = util.UnsafeString2ByteSlice(
			http.MethodGet + "," + http.MethodPost + "," + http.MethodHead + "," +
				http.MethodPut + "," + http.MethodPatch + "," + http.MethodDelete,
		)
	}
	if len(opts.AllowHeaders) > 0 {
		allowHeaders = util.UnsafeString2ByteSlice(strings.Join(opts.AllowHeaders, ","))
	}
	if len(opts.ExposeHeaders) > 0 {
		exposeHeaders = util.UnsafeString2ByteSlice(strings.Join(opts.ExposeHeaders, ","))
	}
	if opts.MaxAge > 0 {
		maxAge = util.UnsafeString2ByteSlice(strconv.Itoa(opts.MaxAge))
	}

	// Setup middleware function
	return func(req *webserver.RequestContext) error {
		origin := util.UnsafeByteSlice2String(req.RequestHeaders().PeekBytes(util.HeaderOrigin))
		allowOrigin := ""

		if len(origin) > 0 {
			req.ResponseHeaders().AddBytesKV(util.HeaderVary, util.HeaderOrigin)
		}

		// See https://developer.mozilla.org/en-US/docs/Glossary/Preflight_request for details on how to handle
		// preflight request.
		preflight := req.IsOptions()

		// No origin provided?
		if len(origin) == 0 {
			if preflight {
				req.NoContent(http.StatusNoContent)
				return nil
			}
			// Go to next middleware
			return req.Next()
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
				req.NoContent(http.StatusNoContent)
				return nil
			}
			// Go to next middleware
			return req.Next()
		}

		respHdrs := req.ResponseHeaders()

		respHdrs.SetBytesK(util.HeaderAccessControlAllowOrigin, allowOrigin)
		if opts.AllowCredentials {
			respHdrs.SetBytesKV(util.HeaderAccessControlAllowCredentials, util.BytesTrue)
		}

		// Simple request
		if !preflight {
			if exposeHeaders != nil {
				respHdrs.SetBytesKV(util.HeaderAccessControlExposeHeaders, exposeHeaders)
			}
			// Go to next middleware
			return req.Next()
		}

		// Preflight request
		respHdrs.AddBytesKV(util.HeaderVary, util.HeaderAccessControlRequestMethod)
		respHdrs.AddBytesKV(util.HeaderVary, util.HeaderAccessControlRequestHeaders)
		respHdrs.SetBytesKV(util.HeaderAccessControlAllowMethods, allowMethods)
		if allowHeaders != nil {
			respHdrs.SetBytesKV(util.HeaderAccessControlAllowHeaders, allowHeaders)
		} else {
			header := req.RequestHeaders().PeekBytes(util.HeaderAccessControlRequestHeaders)
			if len(header) > 0 {
				respHdrs.SetBytesKV(util.HeaderAccessControlAllowHeaders, header)
			}
		}
		if maxAge != nil {
			respHdrs.SetBytesKV(util.HeaderAccessControlMaxAge, maxAge)
		}

		// Done
		req.NoContent(http.StatusNoContent)
		return nil
	}
}
