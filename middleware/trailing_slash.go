package middleware

import (
	"bytes"

	webserver "github.com/randlabs/go-webserver/v2"
	"github.com/randlabs/go-webserver/v2/util"
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
func NewTrailingSlash(opts TrailingSlashOptions) webserver.HandlerFunc {
	// Setup middleware function
	return func(req *webserver.RequestContext) error {
		uri := req.URI()

		// Check if the path contains (or not) the trailing slash
		modified := false
		path := uri.Path()
		if opts.Remove {
			if bytes.HasSuffix(path, []byte{47}) {
				path = bytes.TrimRight(path, "/\\")
				modified = true
			}
		} else {
			if !bytes.HasSuffix(path, []byte{47}) {
				path = append(path, '/')
				modified = true
			}
		}
		if modified {
			strPath, err := util.SanitizeUrlPath(util.UnsafeByteSlice2String(path))
			if err != nil {
				req.BadRequest(err.Error())
				return nil
			}
			uri.SetPath(strPath)

			// Redirect
			if opts.RedirectCode != 0 {
				req.Redirect(util.UnsafeByteSlice2String(uri.FullURI()), int(opts.RedirectCode))
				return nil
			}
		}

		// Go to next middleware
		return req.Next()
	}
}
