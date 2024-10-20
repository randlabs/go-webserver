// See the LICENSE file for license details.

package util

// -----------------------------------------------------------------------------

var (
	BytesHttp = []byte("http")

	BytesTrue = []byte("true")

	MethodOptions = []byte("OPTIONS")

	HeaderAccessControlAllowCredentials      = []byte("Access-Control-Allow-Credentials")
	HeaderAccessControlAllowHeaders          = []byte("Access-Control-Allow-Headers")
	HeaderAccessControlAllowMethods          = []byte("Access-Control-Allow-Methods")
	HeaderAccessControlAllowOrigin           = []byte("Access-Control-Allow-Origin")
	HeaderAccessControlAllowPrivateNetwork   = []byte("Access-Control-Allow-Private-Network")
	HeaderAccessControlExposeHeaders         = []byte("Access-Control-Expose-Headers")
	HeaderAccessControlMaxAge                = []byte("Access-Control-Max-Age")
	HeaderAccessControlRequestHeaders        = []byte("Access-Control-Request-Headers")
	HeaderAccessControlRequestMethod         = []byte("Access-Control-Request-Method")
	HeaderAccessControlRequestPrivateNetwork = []byte("Access-Control-Request-Private-Network")
	HeaderAuthorization                      = []byte("Authorization")
	HeaderCacheControl                       = []byte("Cache-Control")
	HeaderContentSecurityPolicy              = []byte("Content-Security-Policy")
	HeaderContentSecurityPolicyReportOnly    = []byte("Content-Security-Policy-Report-Only")
	HeaderContentType                        = []byte("Content-Type")
	HeaderContentLength                      = []byte("Content-Length")
	HeaderETag                               = []byte("Etag")
	HeaderHeaderRetryAfter                   = []byte("Retry-After")
	HeaderIfNoneMatch                        = []byte("If-None-Match")
	HeaderOrigin                             = []byte("Origin")
	HeaderReferrerPolicy                     = []byte("Referrer-Policy")
	HeaderStrictTransportSecurity            = []byte("Strict-Transport-Security")
	HeaderTrueClientIP                       = []byte("True-Client-IP")
	HeaderVary                               = []byte("Vary")
	HeaderXContentTypeOptions                = []byte("X-Content-Type-Options")
	HeaderXForwardedHost                     = []byte("X-Forwarded-Host")
	HeaderXForwardedFor                      = []byte("X-Forwarded-For")
	HeaderXForwardedProto                    = []byte("X-Forwarded-Proto")
	HeaderXForwardedProtocol                 = []byte("X-Forwarded-Protocol")
	HeaderXForwardedSsl                      = []byte("X-Forwarded-Ssl")
	HeaderXFrameOptions                      = []byte("X-Frame-Options")
	HeaderXRateLimitLimit                    = []byte("X-Rate-Limit-Limit")
	HeaderXRateLimitRemaining                = []byte("X-Rate-Limit-Remaining")
	HeaderXRateLimitReset                    = []byte("X-Rate-Limit-Reset")
	HeaderXUrlScheme                         = []byte("X-Url-Scheme")
	HeaderXXSSProtection                     = []byte("X-XSS-Protection")

	ContentTypeApplicationJSON = []byte("application/json; charset=utf-8")
	ContentTypeTextPlain       = []byte("text/plain; charset=utf-8")
)
