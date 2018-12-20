package clccam

/*
 * Options handling: taken from and thanks to github.com/zpatrick/rclient
 */

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/rehttp"
	"github.com/grrtrr/clccam/logger"
)

// A ClientOption configures the @client using the functional option pattern.
type ClientOption func(client *Client)

// HostURL sets the base of the client if @host non-empty, otherwise uses the default base URL.
func HostURL(host string) ClientOption {
	return func(r *Client) {
		if host != "" {
			u, err := url.Parse(strings.TrimRight(host, "/"))
			if err != nil {
				logger.Fatalf("invalid URL %q: %s", host, err)
			}
			u.Scheme = "https"
			r.baseURL = u.String()
		}
	}
}

// Retryer configures the retry mechanism of the client
// @maxRetries: maximum number of retries per request
// @stepDelay:  base value for exponential backoff + jitter delay
// @maxTimeout: maximum overall client timeout
func Retryer(maxRetries int, stepDelay, maxTimeout time.Duration) ClientOption {
	return func(r *Client) {
		r.client.Transport = rehttp.NewTransport(
			r.client.Transport, // Wrap existing transport
			rehttp.RetryFn(func(at rehttp.Attempt) bool {
				if at.Index < maxRetries {
					if at.Response == nil {
						logger.Warnf("%s %s failed (%s) - retry #%d",
							at.Request.Method, at.Request.URL.Path, at.Error, at.Index+1)
						return true
					}
					switch at.Response.StatusCode {
					case http.StatusRequestTimeout:
						fallthrough
					case http.StatusInternalServerError:
						fallthrough
					case http.StatusBadGateway:
						fallthrough
					case http.StatusServiceUnavailable:
						fallthrough
					case http.StatusGatewayTimeout:

						logger.Warnf("%s %s returned %q - retry #%d",
							at.Request.Method, at.Request.URL.Path, at.Response.Status, at.Index+1)
						return true
					}
				}
				return false
			}),
			// Reuse @maxTimeout as upper bound for the exponential backoff.
			rehttp.ExpJitterDelay(stepDelay, maxTimeout),
		)
		// Set the overall client timeout in lock-step with that of the retryer.
		r.client.Timeout = maxTimeout
	}
}

// Context adds @ctx to the client
func Context(ctx context.Context) ClientOption {
	return func(r *Client) {
		r.ctx = ctx
	}
}

// Debug enables printing request/response to stderr.
func Debug(enabled bool) ClientOption {
	return func(r *Client) {
		r.requestDebug = enabled
	}
}

// JsonResponse enables printing JSON responses to stdout.
func JsonResponse(enabled bool) ClientOption {
	return func(r *Client) {
		r.jsonResponse = enabled
	}
}

// InsecureTLS disables SSL certificate validation. Use with caution.
func InsecureTLS(enable bool) ClientOption {
	return func(r *Client) {
		if tr, ok := r.client.Transport.(*http.Transport); !ok {
			logger.Fatalf("unable to access http client transport attributes - not using http.Transport?")
		} else if tr.TLSClientConfig != nil {
			tr.TLSClientConfig.InsecureSkipVerify = enable
		} else {
			tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: enable}
		}
	}
}

// RequestOptions sets the RequestOptions field of @r.
func RequestOptions(options ...RequestOption) ClientOption {
	return func(r *Client) {
		r.requestOptions = append(r.requestOptions, options...)
	}
}

// RequestOption modifies the HTTP @req in place
type RequestOption func(req *http.Request)

// Headers adds the specified names and values as headers to a request
func Headers(headers map[string]string) RequestOption {
	return func(req *http.Request) {
		for name, val := range headers {
			req.Header.Set(name, val)
		}
	}
}

// Query adds the specified query to a request.
func Query(query url.Values) RequestOption {
	return func(req *http.Request) {
		req.URL.RawQuery = query.Encode()
	}
}
