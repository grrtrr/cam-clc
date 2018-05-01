package clccam

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// Visibility of a CAM Entity
type Visibility uint32

const (
	// Visible to Cloud Application Manager users across all organizations.
	Visibility_Public Visibility = iota

	// Visible to all users in the organization where the box was created.
	Visibility_Organization

	// By default, the box is visible only to members of the workspace where it was created.
	Visibility_Workspace

	// This is used e.g. by variables
	Visibility_Internal

	// This is used e.g. by variables
	Visibility_Private
)

// Implements encoding.TextMarshaler
func (v Visibility) MarshalText() ([]byte, error) {
	switch v {
	case Visibility_Public:
		return []byte("public"), nil
	case Visibility_Organization:
		return []byte("organization"), nil
	case Visibility_Workspace:
		return []byte("workspace"), nil
	case Visibility_Internal:
		return []byte("internal"), nil
	case Visibility_Private:
		return []byte("private"), nil
	}
	return nil, fmt.Errorf("invalid Visibility %d", v)
}

// Implements encoding.TextUnmarshaler
func (v *Visibility) UnmarshalText(data []byte) error {
	switch string(data) {
	case "public":
		*v = Visibility_Public
	case "organization":
		*v = Visibility_Organization
	case "workspace":
		*v = Visibility_Workspace
	case "internal":
		*v = Visibility_Internal
	case "private":
		*v = Visibility_Private
	default:
		return fmt.Errorf("invalid Visibility %q", string(data))
	}
	return nil
}

// Implements fmt.Stringer
func (v Visibility) String() string {
	if b, err := v.MarshalText(); err != nil {
		return err.Error()
	} else {
		return string(b)
	}
}

// VisibilityFromString attempts to parse @s as stringified Visibility.
func VisibilityFromString(s string) (val Visibility, err error) {
	err = val.UnmarshalText([]byte(s))
	return val, err
}

// VisibilityStrings returns the list of Visibility string literals, or maps @vals if non-empty.
func VisibilityStrings(vals ...Visibility) (ret []string) {
	if len(vals) > 0 {
		for _, val := range vals {
			ret = append(ret, val.String())
		}
		return ret
	}
	return []string{"public", "organization", "workspace", "internal", "private"}
}

// Implements database/sql/driver.Valuer
func (v Visibility) Value() (driver.Value, error) {
	return v.String(), nil
}

// Implements database/sql.Scanner
func (v *Visibility) Scan(src interface{}) error {
	switch src := src.(type) {
	case int64:
		*v = Visibility(src)
		return nil
	case []byte:
		return v.UnmarshalText(src)
	case string:
		return v.UnmarshalText([]byte(src))
	}
	return fmt.Errorf("unable to convert %T to Visibility", src)
}

// Implements json.Marshaler
func (v Visibility) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

// Implements json.Unmarshaler
func (v *Visibility) UnmarshalJSON(data []byte) error {
	var output string
	if err := json.Unmarshal(data, &output); err != nil {
		return fmt.Errorf("failed to parse '%s' as Visibility: %s", string(data), err)
	}
	return v.UnmarshalText([]byte(output))
}

// Implements yaml.Marshaler
func (v Visibility) MarshalYAML() (interface{}, error) {
	return v.String(), nil
}

// Implements yaml.Unmarshaler
func (v *Visibility) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var output string
	if err := unmarshal(&output); err != nil {
		return fmt.Errorf("failed to unmarshal Visibility from YAML: %s", err)
	}
	return v.UnmarshalText([]byte(output))
}

