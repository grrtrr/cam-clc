package clccam

import (
	"bytes"
	"strconv"
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
		v, err = time.Parse("2006-01-02T15:04:05.000000Z", ts)
		if err != nil {
			//		v, err = time.Parse("2006-01-02T15:04:05.00000Z", ts)
			if err != nil {
				//				v, err = time.Parse("2006-01-02T15:04:05.0000Z", ts)
				//if err != nil {
				//	v, err = time.Parse("2006-01-02T15:04:05.000Z", ts)
			}
		}
		if err != nil {
			return errors.Errorf("invalid timestamp %s", string(b))
		}
	}
	t.Time = v
	return nil
}

func (t Timestamp) String() string {
	return t.Time.String()
}
