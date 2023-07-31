package middleware

import (
	"strings"

	webserver "github.com/randlabs/go-webserver"
	"github.com/randlabs/go-webserver/request"
	"github.com/randlabs/go-webserver/util"
)

// -----------------------------------------------------------------------------

// TrailingSlashOptions defines a middleware that adds or removes trailing slashes in paths.
type TrailingSlashOptions struct {
	// Remove tells the middleware to remove trailing slashes if present.
	// If this setting is false, then the trailing slash is added if absent.
	Remove bool `json:"remove,omitempty"`

	// RedirectCode, if not zero, will make the middleware to return a redirect response.
	RedirectCode uint `json:"redirectCode,omitempty"`
}

// -----------------------------------------------------------------------------

// NewTrailingSlash creates a new middleware to handle trailing slashes in request's paths
func NewTrailingSlash(opts TrailingSlashOptions) webserver.MiddlewareFunc {
	// Setup middleware function
	return func(next webserver.HandlerFunc) webserver.HandlerFunc {
		return func(req *request.RequestContext) error {
			uri := req.URI()

			// Check if the path contains (or not) the trailing slash
			modified := false
			path := string(uri.Path())
			if opts.Remove {
				if strings.HasSuffix(path, "/") {
					path = strings.TrimRight(path, "/\\")
					modified = true
				}
			} else {
				if !strings.HasSuffix(path, "/") {
					path += "/"
					modified = true
				}
			}
			if modified {
				uri.SetPath(util.SanitizeUrlPath(path))

				// Redirect
				if opts.RedirectCode != 0 {
					req.Redirect(string(uri.FullURI()), int(opts.RedirectCode))
					return nil
				}
			}

			// Go to next middleware
			return next(req)
		}
	}
}
