package clccam

import (
	"bytes"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

// URI embeds net/url in order to allow JSON marshaling.
// CAM uses URIs quite extensively.
type URI struct {
	*url.URL
}

func (u URI) String() string {
	if u.URL == nil {
		return ""
	}
	return u.URL.String()
}

// IsZero returns true if @u does not contain a meaningful value.
func (u *URI) IsZero() bool {
	return u == nil || u.URL == nil || u.URL.String() == ""
}

// Implements json.Marshaler
func (u URI) MarshalJSON() (text []byte, err error) {
	if u.URL == nil {
		return []byte(`""`), nil
	}
	return []byte(strconv.Quote(u.URL.String())), nil
}

// Implements json.Unmarshaler
func (u *URI) UnmarshalJSON(b []byte) error {
	u.URL = new(url.URL)
	return u.URL.UnmarshalBinary(bytes.Trim(b, `"`))
}

// UriFromString attempts to parse @s as a Clc CAM URI.
func UriFromString(s string) (*URI, error) {
	var u URI

	if err := json.Unmarshal([]byte(strconv.Quote(s)), &u); err != nil {
		return nil, errors.Errorf("invalid URI %q: %s", s, err)
	}
	return &u, nil
}
