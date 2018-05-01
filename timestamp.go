package clccam

import (
	"bytes"
	"time"

	"github.com/pkg/errors"
)

// CAM-specific timestamp format (e.g. "2018-01-26 19:50:49.131726")
const ts_format = `2006-01-02 15:04:05.999999`

// Timestamp is a date/time string using the format "2018-01-26 19:50:49.131726"
type Timestamp struct {
	time.Time
}

// Implements json.Marshaler
func (t *Timestamp) MarshalJSON() (text []byte, err error) {
	return []byte(t.Time.Format(ts_format)), nil
}

// Implements json.Unmarshaler
func (t *Timestamp) UnmarshalJSON(b []byte) error {
	v, err := time.Parse(ts_format, string(bytes.Trim(b, `"`)))
	if err != nil {
		return errors.Errorf("invalid timestamp %q", string(b))
	}
	t.Time = v
	return nil
}

func (t Timestamp) String() string {
	return t.Time.String()
}
