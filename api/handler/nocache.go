package handler

import (
	"time"

	"github.com/gin-gonic/gin"
)

// NoCache is a simple piece of middleware that sets a number of HTTP headers to prevent
// a router (or subrouter) from being cached by an upstream proxy and/or client.
//
// As per http://wiki.nginx.org/HttpProxyModule - NoCache sets:
//
//	Expires: Thu, 01 Jan 1970 00:00:00 UTC
//	Cache-Control: no-cache, no-store, no-transform, must-revalidate, private, max-age=0
//	Pragma: no-cache (for HTTP/1.0 proxies/clients)
//	X-Accel-Expires: 0
func NoCache() gin.HandlerFunc {
	// Unix epoch time.
	epoch := time.Unix(0, 0).Format(time.RFC1123)

	// Taken from https://github.com/mytrile/nocache
	noCacheHeaders := map[string]string{
		"Expires":                epoch,
		"Cache-Control":          "no-cache, no-store, must-revalidate;",
		"Pragma":                 "no-cache",
		"X-Content-Type-Options": "nosniff",
	}

	// ETag headers array.
	etagHeaders := [6]string{
		"ETag",
		"If-Modified-Since",
		"If-Match",
		"If-None-Match",
		"If-Range",
		"If-Unmodified-Since",
	}

	return func(c *gin.Context) {
		// Delete any ETag headers that may have been set
		for _, v := range etagHeaders {
			if c.Request.Header.Get(v) != "" {
				c.Request.Header.Del(v)
			}
		}

		// Set our NoCache headers
		for k, v := range noCacheHeaders {
			c.Writer.Header().Set(k, v)
		}

		c.Next()
	}
}
