package middleware

import (
	"fmt"
	"strings"

	webserver "github.com/randlabs/go-webserver"
	"github.com/randlabs/go-webserver/request"
)

// -----------------------------------------------------------------------------

// SafeBrowsingOptions defines how common response headers for safe browsing are added.
type SafeBrowsingOptions struct {
	// XXSSProtection sets the `X-XSS-Protection` header to stops pages from loading when they detect reflected
	// cross-site scripting (XSS) attacks.
	// Optional. Defaults to "1; mode=block".
	XXSSProtection string `json:"x-xss-protection,omitempty"`

	// XContentTypeNoSniff sets the `X-Content-Type-Options` header to indicate that the MIME types advertised
	// in the Content-Type headers should be followed and not be changed.
	// Optional. Defaults to "nosniff".
	XContentTypeNoSniff string `json:"x-content-type-options,omitempty"`

	// XFrameOptions can be used to indicate whether a browser should be allowed to render a page in a <frame>,
	// <iframe> or <object>.
	// Optional. Defaults to "sameorigin".
	// Possible values: "sameorigin", "deny", "allow-from uri"
	XFrameOptions string `json:"x-frame-options,omitempty"`

	// HSTS controls the `Strict-Transport-Security` header to inform browsers that the site should only be
	// accessed using HTTPS, and that any future attempts to access it using HTTP should automatically be
	// converted to HTTPS.
	HSTS struct {
		// MaxAge establishes the time, in seconds, that the browser should remember that a site is only to be
		// accessed using HTTPS.
		// Optional. Defaults to 0.
		MaxAge uint `json:"max-age,omitempty"`

		// IncludeSubdomains is used to apply the HSTS settings to all of the site's subdomains as well.
		// Optional.
		IncludeSubdomains bool `json:"include-subdomains,omitempty"`

		// Preload will add the preload tag in the HSTS header. See https://hstspreload.org/ for details.
		// Optional.
		Preload bool `json:"preload,omitempty"`
	} `json:"hsts,omitempty"`

	// ContentSecurityPolicy sets the `Content-Security-Policy` header to enhance security against XSS.
	// Optional.
	ContentSecurityPolicy string `json:"content-security-policy,omitempty"`

	// ContentSecurityPolicyReportOnly would use the `Content-Security-Policy-Report-Only` header instead
	// of the `Content-Security-Policy` header. Used to report violations instead of blocking resources.
	// Optional.
	ContentSecurityPolicyReportOnly bool `json:"csp-report-only,omitempty"`

	// ReferrerPolicy sets the `Referrer-Policy` header providing security against leaking potentially sensitive
	// request paths to third parties.
	// Optional.
	ReferrerPolicy string `json:"referrer-policy,omitempty"`
}

// DefaultSafeBrowsing creates a default SafeBrowsing middleware with commonly used options
func DefaultSafeBrowsing() webserver.MiddlewareFunc {
	return NewSafeBrowsing(SafeBrowsingOptions{
		XXSSProtection:      "1; mode=block",
		XContentTypeNoSniff: "nosniff",
		XFrameOptions:       "sameorigin",
	})
}

// NewSafeBrowsing creates a new SafeBrowsing middleware based on the specified options
func NewSafeBrowsing(opts SafeBrowsingOptions) webserver.MiddlewareFunc {
	// Setup middleware function
	return func(next webserver.HandlerFunc) webserver.HandlerFunc {
		return func(req *request.RequestContext) error {
			// Set X-XSS-Protection header
			if len(opts.XXSSProtection) > 0 {
				req.SetResponseHeader("X-XSS-Protection", opts.XXSSProtection)
			}

			// Set X-Content-Type-Options header
			if len(opts.XContentTypeNoSniff) > 0 {
				req.SetResponseHeader("X-Content-Type-Options", opts.XContentTypeNoSniff)
			}

			// Set X-Frame-Options header
			if len(opts.XFrameOptions) > 0 {
				req.SetResponseHeader("X-Frame-Options", opts.XFrameOptions)
			}

			// Set Strict-Transport-Security header
			if opts.HSTS.MaxAge > 0 && (req.IsTLS() || req.RequestHeader("X-Forwarded-Proto") == "https") {
				subdomains := make([]string, 1)
				subdomains[0] = fmt.Sprintf("max-age=%d", opts.HSTS.MaxAge)
				if opts.HSTS.IncludeSubdomains {
					subdomains = append(subdomains, "includeSubdomains")
				}
				if opts.HSTS.Preload {
					subdomains = append(subdomains, "preload")
				}
				req.SetResponseHeader("Strict-Transport-Security", strings.Join(subdomains, "; "))
			}

			// Set Content-Security-Policy/Content-Security-Policy-Report-Only header
			if len(opts.ContentSecurityPolicy) > 0 {
				if opts.ContentSecurityPolicyReportOnly {
					req.SetResponseHeader("Content-Security-Policy-Report-Only", opts.ContentSecurityPolicy)
				} else {
					req.SetResponseHeader("Content-Security-Policy", opts.ContentSecurityPolicy)
				}
			}

			// Set Referrer-Policy header
			if len(opts.ReferrerPolicy) > 0 {
				req.SetResponseHeader("Referrer-Policy", opts.ReferrerPolicy)
			}

			// Go to next middleware
			return next(req)
		}
	}
}
