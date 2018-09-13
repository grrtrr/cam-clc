package clccam

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"reflect"
	"regexp"
	"strings"

	"github.com/grrtrr/clccam/logger"
	"github.com/pkg/errors"
)

// Client is a reusable REST client for CAM API calls.
type Client struct {
	// client performs the actual requests.
	client *http.Client

	// Base URL to use.
	baseURL string

	// Per-request options.
	requestOptions []RequestOption

	// Cancellation context (used by @cancel). Can be overridden via WithContext()
	ctx context.Context

	// Print request / response to stderr.
	requestDebug bool

	// Print JSON response to stdout.
	jsonResponse bool
}

// NewClient returns a new standalone client.
func NewClient(options ...ClientOption) *Client {
	var c = &Client{client: &http.Client{}}

	for _, setOption := range options {
		setOption(c)
	}
	return c
}

// NewClient turns @t into a CAM client.
func (t Token) NewClient(options ...ClientOption) *Client {
	var c = NewClient(HostURL("https://cam.ctl.io"),
		RequestOptions(Headers(map[string]string{
			"Authorization": "Bearer " + string(t),
			"Content-Type":  "application/json; charset=utf-8",
			"Accept":        "application/json",
		})),
	)

	// Apply the provided options last, to override any defaults
	for _, setOption := range options {
		setOption(c)
	}
	return c
}

// WithDebug enables debugging on @c.
func (c *Client) WithDebug(enabled bool) *Client {
	c.requestDebug = enabled
	return c
}

// WithContext sets the context to @ctx.
func (c *Client) WithContext(ctx context.Context) *Client {
	c.ctx = ctx
	return c
}

// WithJsonResponse enables printing the JSON response to stdout.
func (c *Client) WithJsonResponse(enabled bool) *Client {
	c.jsonResponse = enabled
	return c
}

// Get performs a GET /path, with output into @resModel
func (c *Client) Get(path string, resModel interface{}) error {
	return c.getResponse(path, "GET", nil, resModel)
}

// getResponse performs a generic request
// @urlPath:  request path relative to %BaseURL
// @verb:     request verb
// @reqModel: request model to serialize, or nil.
// @resModel: result model to deserialize, must be a pointer to the expected result, or nil.
// Evaluates the StatusCode of the BaseResponse (embedded) in @inModel and sets @err accordingly.
// If @err == nil, fills in @resModel, else returns error.
func (c *Client) getResponse(urlPath, verb string, reqModel, resModel interface{}) error {
	var (
		url     = fmt.Sprintf("%s/%s", c.baseURL, strings.TrimLeft(urlPath, "/"))
		reqBody io.Reader
	)

	if reqModel != nil {
		jsonReq, err := json.Marshal(reqModel)
		if err != nil {
			return errors.Wrapf(err, "failed to encode request model %T %+v", reqModel, reqModel)
		}
		reqBody = bytes.NewBuffer(jsonReq)
	}

	// resModel must be a pointer type (call-by-value)
	if resModel != nil {
		if resType := reflect.TypeOf(resModel); resType.Kind() != reflect.Ptr {
			return errors.Errorf("expecting pointer to result model %T", resModel)
		}
	}

	req, err := http.NewRequest(verb, url, reqBody)
	if err != nil {
		return err
	}

	if c.ctx != nil {
		req = req.WithContext(c.ctx)
	}

	for _, setOption := range c.requestOptions {
		setOption(req)
	}

	if c.requestDebug {
		reqDump, _ := httputil.DumpRequest(req, true)
		logger.Debugf("%s", reqDump)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if c.requestDebug {
		resDump, _ := httputil.DumpResponse(res, true)
		logger.Debugf("%s", resDump)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil && res.ContentLength > 0 {
		res.Body.Close()
		return errors.Wrapf(err, "failed to read error response %d body", res.StatusCode)
	} else if err := res.Body.Close(); err != nil {
		return errors.Wrapf(err, "failed to close after reading response body")
	}

	switch res.StatusCode {
	case 200, 201, 202, 204: // OK | CREATED | ACCEPTED | NO CONTENT
		if c.requestDebug && len(body) > 0 && !strings.Contains(http.DetectContentType(body), "html") {
			logger.Debugf("%s", string(body))
		}

		if c.jsonResponse {
			var b bytes.Buffer

			if err := json.Indent(&b, body, "", "\t"); err != nil {
				return errors.Wrapf(err, "failed to decode JSON response %q", string(body))
			}
			fmt.Println(b.String())
		}

		if resModel != nil {
			if res.ContentLength == 0 {
				return errors.Errorf("unable do populate %T result model, due to empty %q response",
					resModel, res.Status)
			}
			return json.Unmarshal(body, resModel)
		} else if res.ContentLength > 0 {
			return errors.Errorf("unable to decode non-empty %q response (%d bytes) to nil response model",
				res.Status, res.ContentLength)
		}
		return nil
	default: // Errors and temporary failures
		if len(body) > 0 && !strings.Contains(http.DetectContentType(body), "html") {
			// Decode possible CAM error response:
			// 1) text/html:  HTML page - skip as per above check
			// 2) text/plain: use body after stripping whitespace
			// 3) bare JSON string
			// 4) struct { message: "string" }
			var payload map[string]interface{}
			var errMsg = string(bytes.TrimSpace(body))

			if err := json.Unmarshal(body, &payload); err != nil {
				// Failed to decode as struct, try string (2,3)
				if err = json.Unmarshal(body, &errMsg); err != nil {
					var nl = regexp.MustCompile(`(\r?\n)+`)

					errMsg = nl.ReplaceAllString(string(bytes.TrimSpace(body)), "; ")
				}
			} else if errors, ok := payload["message"]; ok {
				if msg, ok := errors.(string); ok {
					errMsg = strings.TrimRight(msg, " .") // sometimes they end error messages in '.'
				}
			} else if error, ok := payload["error"]; ok {
				if msg, ok := error.(string); ok {
					errMsg = fmt.Sprintf("Error - %s", msg)
				}
			}
			return errors.Errorf("%s (status: %d)", errMsg, res.StatusCode)
		}
		// FIXME: implement temporary / retryable errors (300)
		return errors.New(res.Status)
	}
}
