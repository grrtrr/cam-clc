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
	"time"

	"github.com/grrtrr/clccam/logger"
	"github.com/pkg/errors"
)

// Client is a reusable REST client for CAM API calls.
type Client struct {
	// client performs the actual requests
	client *http.Client

	// Base URL to use
	baseURL string

	// Per-request options
	requestOptions []RequestOption

	// Cancellation context (used by @cancel). Can be overridden via WithContext()
	ctx context.Context

	// enable verbose per-request debugging
	requestDebug bool
}

// NewClient returns a new standalone client.
func NewClient(options ...ClientOption) *Client {
	var c = &Client{client: &http.Client{}}

	for _, setOption := range options {
		setOption(c)
	}
	return c
}

// NewClient turns @t into a CAM client with usable defaults
func (t Token) NewClient(options ...ClientOption) *Client {
	var c = NewClient(HostURL("https://cam.ctl.io"),
		RequestOptions(Headers(map[string]string{
			"Authorization": "Bearer " + string(t),
			"Content-Type":  "application/json; charset=utf-8",
			"Accept":        "application/json",
		})),
		Retryer(3, 10*time.Second, 180*time.Second),
		Debug(true),
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

// WithContext sets the context to @ctx
func (c *Client) WithContext(ctx context.Context) *Client {
	c.ctx = ctx
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
			return errors.Errorf("failed to encode request model %T %+v: %s", reqModel, reqModel, err)
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
	defer res.Body.Close()

	if c.requestDebug {
		resDump, _ := httputil.DumpResponse(res, true)
		logger.Debugf("%s", resDump)
	}

	switch res.StatusCode {
	case 200, 201, 202, 204: // OK | CREATED | ACCEPTED | NO CONTENT
		if resModel != nil {
			if res.ContentLength == 0 {
				return errors.Errorf("unable do populate %T result model, due to empty %q response",
					resModel, res.Status)
			}
			return json.NewDecoder(res.Body).Decode(resModel)
		} else if res.ContentLength > 0 {
			return errors.Errorf("unable to decode non-empty %q response (%d bytes) to nil response model",
				res.Status, res.ContentLength)
		}
		return nil
	}

	// Remaining cases
	if body, err := ioutil.ReadAll(res.Body); err != nil && res.ContentLength > 0 {
		return errors.Errorf("failed to read error response %d body: %s", res.StatusCode, err)
	} else if len(body) > 0 && !strings.Contains(http.DetectContentType(body), "html") {
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
	return errors.New(res.Status)
}
