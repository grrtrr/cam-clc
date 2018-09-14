package clccam

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// Operations performed on an instance.
type InstanceEvent uint32

const (
	// Shut down
	InstanceEvent_shutdown InstanceEvent = iota

	// Power On
	InstanceEvent_poweron

	// Re-Install
	InstanceEvent_reinstall

	// Reconfigure
	InstanceEvent_reconfigure

	// Terminate
	InstanceEvent_terminate

	// Terminate Service
	InstanceEvent_terminate_service
)

// Implements encoding.TextMarshaler
func (i InstanceEvent) MarshalText() ([]byte, error) {
	switch i {
	case InstanceEvent_shutdown:
		return []byte("shutdown"), nil
	case InstanceEvent_poweron:
		return []byte("poweron"), nil
	case InstanceEvent_reinstall:
		return []byte("reinstall"), nil
	case InstanceEvent_reconfigure:
		return []byte("reconfigure"), nil
	case InstanceEvent_terminate:
		return []byte("terminate"), nil
	case InstanceEvent_terminate_service:
		return []byte("terminate_service"), nil
	}
	return nil, fmt.Errorf("invalid InstanceEvent %d", i)
}

// Implements encoding.TextUnmarshaler
func (i *InstanceEvent) UnmarshalText(data []byte) error {
	switch string(data) {
	case "shutdown":
		*i = InstanceEvent_shutdown
	case "poweron":
		*i = InstanceEvent_poweron
	case "reinstall":
		*i = InstanceEvent_reinstall
	case "reconfigure":
		*i = InstanceEvent_reconfigure
	case "terminate":
		*i = InstanceEvent_terminate
	case "terminate_service":
		*i = InstanceEvent_terminate_service
	default:
		return fmt.Errorf("invalid InstanceEvent %q", string(data))
	}
	return nil
}

// Implements fmt.Stringer
func (i InstanceEvent) String() string {
	if b, err := i.MarshalText(); err != nil {
		return err.Error()
	} else {
		return string(b)
	}
}

// InstanceEventFromString attempts to parse @s as stringified InstanceEvent.
func InstanceEventFromString(s string) (val InstanceEvent, err error) {
	err = val.UnmarshalText([]byte(s))
	return val, err
}

// InstanceEventStrings returns the list of InstanceEvent string literals, or maps @vals if non-empty.
func InstanceEventStrings(vals ...InstanceEvent) (ret []string) {
	if len(vals) > 0 {
		for _, val := range vals {
			ret = append(ret, val.String())
		}
		return ret
	}
	return []string{"shutdown", "poweron", "reinstall", "reconfigure", "terminate", "terminate_service"}
}

// Implements database/sql/driver.Valuer
func (i InstanceEvent) Value() (driver.Value, error) {
	return i.String(), nil
}

// Implements database/sql.Scanner
func (i *InstanceEvent) Scan(src interface{}) error {
	switch src := src.(type) {
	case int64:
		*i = InstanceEvent(src)
		return nil
	case []byte:
		return i.UnmarshalText(src)
	case string:
		return i.UnmarshalText([]byte(src))
	}
	return fmt.Errorf("unable to convert %T to InstanceEvent", src)
}

// Implements json.Marshaler
func (i InstanceEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// Implements json.Unmarshaler
func (i *InstanceEvent) UnmarshalJSON(data []byte) error {
	var output string
	if err := json.Unmarshal(data, &output); err != nil {
		return fmt.Errorf("failed to parse '%s' as InstanceEvent: %s", string(data), err)
	}
	return i.UnmarshalText([]byte(output))
}

// Implements yaml.Marshaler
func (i InstanceEvent) MarshalYAML() (interface{}, error) {
	return i.String(), nil
}

// Implements yaml.Unmarshaler
func (i *InstanceEvent) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var output string
	if err := unmarshal(&output); err != nil {
		return fmt.Errorf("failed to unmarshal InstanceEvent from YAML: %s", err)
	}
	return i.UnmarshalText([]byte(output))
}

