package clccam

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// State of an instance or a machine.
type InstanceState uint32

const (
	InstanceState_processing InstanceState = iota
	InstanceState_done
	InstanceState_unavailable
)

// Implements encoding.TextMarshaler
func (i InstanceState) MarshalText() ([]byte, error) {
	switch i {
	case InstanceState_processing:
		return []byte("processing"), nil
	case InstanceState_done:
		return []byte("done"), nil
	case InstanceState_unavailable:
		return []byte("unavailable"), nil
	}
	return nil, fmt.Errorf("invalid InstanceState %d", i)
}

// Implements fmt.Stringer
func (i InstanceState) String() string {
	if b, err := i.MarshalText(); err != nil {
		return err.Error()
	} else {
		return string(b)
	}
}

// Implements encoding.TextUnmarshaler
func (i *InstanceState) UnmarshalText(data []byte) error {
	switch string(data) {
	case "processing":
		*i = InstanceState_processing
	case "done":
		*i = InstanceState_done
	case "unavailable":
		*i = InstanceState_unavailable
	default:
		return fmt.Errorf("invalid InstanceState %q", string(data))
	}
	return nil
}

// Implements flag.Value
func (i *InstanceState) Set(s string) error {
	return i.UnmarshalText([]byte(s))
}

// Implements pflag.Value (superset of flag.Value)
func (i InstanceState) Type() string {
	return "InstanceState"
}

// InstanceStateFromString attempts to parse @s as stringified InstanceState.
func InstanceStateFromString(s string) (val InstanceState, err error) {
	return val, val.Set(s)
}

// InstanceStateStrings returns the list of InstanceState string literals, or maps @vals if non-empty.
func InstanceStateStrings(vals ...InstanceState) (ret []string) {
	if len(vals) > 0 {
		for _, val := range vals {
			ret = append(ret, val.String())
		}
		return ret
	}
	return []string{"processing", "done", "unavailable"}
}

// Implements database/sql/driver.Valuer
func (i InstanceState) Value() (driver.Value, error) {
	return i.String(), nil
}

// Implements database/sql.Scanner
func (i *InstanceState) Scan(src interface{}) error {
	switch src := src.(type) {
	case int64:
		*i = InstanceState(src)
		return nil
	case []byte:
		return i.UnmarshalText(src)
	case string:
		return i.UnmarshalText([]byte(src))
	}
	return fmt.Errorf("unable to convert %T to InstanceState", src)
}

// Implements json.Marshaler
func (i InstanceState) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// Implements json.Unmarshaler
func (i *InstanceState) UnmarshalJSON(data []byte) error {
	var output string
	if err := json.Unmarshal(data, &output); err != nil {
		return fmt.Errorf("failed to parse '%s' as InstanceState: %s", string(data), err)
	}
	return i.UnmarshalText([]byte(output))
}

// Implements yaml.Marshaler
func (i InstanceState) MarshalYAML() (interface{}, error) {
	return i.String(), nil
}

// Implements yaml.Unmarshaler
func (i *InstanceState) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var output string
	if err := unmarshal(&output); err != nil {
		return fmt.Errorf("failed to unmarshal InstanceState from YAML: %s", err)
	}
	return i.UnmarshalText([]byte(output))
}

