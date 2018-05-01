package clccam

import (
	"bytes"
	"net/url"
	"strconv"
)

// URI embeds net/url in order to allow JSON marshaling.
// CAM uses URIs quite extensively.
type URI struct {
	*url.URL
}

func (u URI) String() string {
	if u.URL == nil {
		return "<nil>"
	}
	return u.URL.String()
}

// Implements json.Marshaler
func (u URI) MarshalJSON() (text []byte, err error) {
	if u.URL == nil {
		return []byte(nil), nil
	}
	return []byte(strconv.Quote(u.URL.String())), nil
}

// Implements json.Unmarshaler
func (u *URI) UnmarshalJSON(b []byte) error {
	u.URL = new(url.URL)
	return u.URL.UnmarshalBinary(bytes.Trim(b, `"`))
}
