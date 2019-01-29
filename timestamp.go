package clccam

import (
	"bytes"
	"regexp"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const (
	// CAM-specific timestamp format (e.g. "2018-01-26 19:50:49.131726")
	ts_format = `2006-01-02 15:04:05.999999`

	// Alternative timestamp format that is also in use.
	// See below for rationale why fractional seconds are ignored in this case.
	ts_format_alt = `2006-01-02T15:04:05Z`
)

// Timestamp is a date/time string using the format "2018-01-26 19:50:49.131726"
type Timestamp struct {
	time.Time
}

// Implements json.Marshaler
func (t *Timestamp) MarshalJSON() (text []byte, err error) {
	return []byte(strconv.Quote(t.Time.Format(ts_format))), nil
}

// Implements json.Unmarshaler
func (t *Timestamp) UnmarshalJSON(b []byte) error {
	var ts = string(bytes.Trim(b, `"`))

	v, err := time.Parse(ts_format, ts)
	if err != nil {
		// Note: sometimes these alternative formats are used:
		// a) "2019-01-12T00:15:09.026751Z"
		// b) "2019-01-12T00:39:30.80067Z"
		// c) "2019-01-12T00:59:37.7309Z"
		// To avoid writing n formats to match these variants, simply ignore fractional seconds here.
		var ts_int = regexp.MustCompile(`\.\d+Z\s*$`).ReplaceAllString(ts, "Z")

		v, err = time.Parse(ts_format_alt, ts_int)
		if err != nil {
			return errors.Errorf("invalid timestamp format %q", string(b))
		}
	}
	t.Time = v
	return nil
}

func (t Timestamp) String() string {
	return t.Time.String()
}
