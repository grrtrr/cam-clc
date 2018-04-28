package clccam

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/go-resty/resty"
	"github.com/grrtrr/clccam/logger"
	"github.com/pkg/errors"
)

const (
	// Upper bound on client operations - default timeout value
	ClientTimeout = 180 * time.Second

	// CAM main API url
	BaseURL = "https://cam.ctl.io"

	// Maximum number of retries per request.
	MaxRetries = 3

	// Per-request retry delay for the retryer.
	StepDelay = time.Second * 10
)

// GLOBALS
var (
	// Enable per-request debugging
	RequestDebug bool
)

// Client wraps resty.Client, since some of the CAM defaults are not exactly REST
type Client struct {
	client *resty.Client

	// Cancellation context for this client and its requests
	ctx context.Context
}

// NewClient returns a new client initialied from @t
func (t Token) NewClient() *Client {
	var c = resty.New().
		SetRESTMode().
		SetDebug(RequestDebug && false).
		SetAuthToken(string(t)).
		SetHostURL("https://cam.ctl.io").
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		}).
		SetTimeout(ClientTimeout).
		AddRetryCondition(resty.RetryConditionFunc(func(res *resty.Response) (bool, error) {
			switch res.RawResponse.StatusCode {
			// Request timeout, server error, bad gateway, service unavailable, gateway timeout
			case 408, 500, 502, 503, 504:
				logger.Warnf("%s %s returned %q - retrying",
					res.Request.Method, res.RawResponse.Request.URL.Path, res.RawResponse.Status)
				return true, nil
			}
			return false, nil
		})).
		SetRetryCount(MaxRetries)
		//			SetRetryWaitTime(time),

	// Try to work around the use of log.Logger
	c.Log.SetPrefix("CAM")
	c.Log.SetFlags(0)
	c.Log.SetOutput(logger.Writer())
	return &Client{client: c, ctx: context.Background()}
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
// @url:      request path relative to %BaseURL
// @verb:     request verb
// @reqModel: request model to serialize, or nil.
// @resModel: result model to deserialize, must be a pointer to the expected result, or nil.
// Evaluates the StatusCode of the BaseResponse (embedded) in @inModel and sets @err accordingly.
// If @err == nil, fills in @resModel, else returns error.
func (c *Client) getResponse(url, verb string, reqModel, resModel interface{}) error {
	var (
		req = c.client.R()
		res *resty.Response
		err error
	)

	if reqModel != nil {
		req = req.SetBody(reqModel)
	}

	/* resModel must be a pointer type (call-by-value) */
	if resModel != nil {
		if resType := reflect.TypeOf(resModel); resType.Kind() != reflect.Ptr {
			return errors.Errorf("expecting pointer to result model %T", resModel)
		}
	}

	if c.ctx != nil {
		req = req.SetContext(c.ctx)
	}

	if RequestDebug {
		logger.Debugf("%s %s", verb, url)
	}

	switch verb {
	case "GET":
		res, err = req.Get(url)
	default:
		// func (*Request) SetContentLength
		/* XXX
		req, err := http.NewRequest(verb, url, reqBody)
		if err != nil {
			return err

		}
		*/
	}

	if RequestDebug {
		logger.Debugf("%s: %s", res.Status(), res.Body())
	}

	if err != nil {
		return err
	}

	switch res.StatusCode() {
	case 200, 201, 202, 204: /* OK / CREATED / ACCEPTED / NO CONTENT */
		if resModel != nil {
			if res.Size() == 0 {
				return errors.Errorf("Unable do populate %T result model, due to empty %q response",
					resModel, res.Status)
			}
			return json.Unmarshal(res.Body(), resModel)
		} else if res.Size() > 0 {
			return errors.Errorf("Unable to decode non-empty %q response (%d bytes) to nil response model",
				res.Status(), res.Size())
		}
		return nil
	}

	// Remaining error cases: res.ContentLength is not reliable - in the SBS case, it uses
	// Transfer-Encoding "chunked", without a Content-Length.
	if res.Size() > 0 {
		var payload map[string]interface{}
		var errMsg = string(res.Body())

		//
		// Decode error response:
		// 1) bare JSON string
		// 2) struct { message: "string" }
		//
		if err := json.Unmarshal(res.Body(), &payload); err != nil {
			/* Failed to decode as struct, try string (1) next. */
			if err = json.Unmarshal(res.Body(), &errMsg); err != nil {
				errMsg = string(res.Body())
			}
		} else if errors, ok := payload["message"]; ok {
			if msg, ok := errors.(string); ok {
				errMsg = msg
			}
		} else if error, ok := payload["error"]; ok {
			if msg, ok := error.(string); ok {
				errMsg = fmt.Sprintf("Error - %s", msg)
			}
		}
		return errors.Errorf("%s (status: %d)", errMsg, res.StatusCode())
	}
	return errors.New(res.Status())
}
