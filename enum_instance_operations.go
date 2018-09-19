package clccam

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// Operations performed on an instance.
type InstanceOp uint32

const (
	// Deploy
	InstanceOp_deploy InstanceOp = iota

	// Shut down
	InstanceOp_shutdown

	// Shut down service
	InstanceOp_shutdown_service

	// Power On
	InstanceOp_poweron

	// Re-Install
	InstanceOp_reinstall

	// Reconfigure
	InstanceOp_reconfigure

	// Terminate
	InstanceOp_terminate

	// Terminate Service
	InstanceOp_terminate_service

	// Snapshot
	InstanceOp_snapshot
)

// Implements encoding.TextMarshaler
func (i InstanceOp) MarshalText() ([]byte, error) {
	switch i {
	case InstanceOp_deploy:
		return []byte("deploy"), nil
	case InstanceOp_shutdown:
		return []byte("shutdown"), nil
	case InstanceOp_shutdown_service:
		return []byte("shutdown_service"), nil
	case InstanceOp_poweron:
		return []byte("poweron"), nil
	case InstanceOp_reinstall:
		return []byte("reinstall"), nil
	case InstanceOp_reconfigure:
		return []byte("reconfigure"), nil
	case InstanceOp_terminate:
		return []byte("terminate"), nil
	case InstanceOp_terminate_service:
		return []byte("terminate_service"), nil
	case InstanceOp_snapshot:
		return []byte("snapshot"), nil
	}
	return nil, fmt.Errorf("invalid InstanceOp %d", i)
}

// Implements fmt.Stringer
func (i InstanceOp) String() string {
	if b, err := i.MarshalText(); err != nil {
		return err.Error()
	} else {
		return string(b)
	}
}

// Implements encoding.TextUnmarshaler
func (i *InstanceOp) UnmarshalText(data []byte) error {
	switch string(data) {
	case "deploy":
		*i = InstanceOp_deploy
	case "shutdown":
		*i = InstanceOp_shutdown
	case "shutdown_service":
		*i = InstanceOp_shutdown_service
	case "poweron":
		*i = InstanceOp_poweron
	case "reinstall":
		*i = InstanceOp_reinstall
	case "reconfigure":
		*i = InstanceOp_reconfigure
	case "terminate":
		*i = InstanceOp_terminate
	case "terminate_service":
		*i = InstanceOp_terminate_service
	case "snapshot":
		*i = InstanceOp_snapshot
	default:
		return fmt.Errorf("invalid InstanceOp %q", string(data))
	}
	return nil
}

// Implements flag.Value
func (i *InstanceOp) Set(s string) error {
	return i.UnmarshalText([]byte(s))
}

// Implements pflag.Value (superset of flag.Value)
func (i InstanceOp) Type() string {
	return "InstanceOp"
}

// InstanceOpFromString attempts to parse @s as stringified InstanceOp.
func InstanceOpFromString(s string) (val InstanceOp, err error) {
	return val, val.Set(s)
}

// InstanceOpStrings returns the list of InstanceOp string literals, or maps @vals if non-empty.
func InstanceOpStrings(vals ...InstanceOp) (ret []string) {
	if len(vals) > 0 {
		for _, val := range vals {
			ret = append(ret, val.String())
		}
		return ret
	}
	return []string{"deploy", "shutdown", "shutdown_service", "poweron", "reinstall", "reconfigure", "terminate", "terminate_service", "snapshot"}
}

// Implements database/sql/driver.Valuer
func (i InstanceOp) Value() (driver.Value, error) {
	return i.String(), nil
}

// Implements database/sql.Scanner
func (i *InstanceOp) Scan(src interface{}) error {
	switch src := src.(type) {
	case int64:
		*i = InstanceOp(src)
		return nil
	case []byte:
		return i.UnmarshalText(src)
	case string:
		return i.UnmarshalText([]byte(src))
	}
	return fmt.Errorf("unable to convert %T to InstanceOp", src)
}

// Implements json.Marshaler
func (i InstanceOp) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// Implements json.Unmarshaler
func (i *InstanceOp) UnmarshalJSON(data []byte) error {
	var output string
	if err := json.Unmarshal(data, &output); err != nil {
		return fmt.Errorf("failed to parse '%s' as InstanceOp: %s", string(data), err)
	}
	return i.UnmarshalText([]byte(output))
}

// Implements yaml.Marshaler
func (i InstanceOp) MarshalYAML() (interface{}, error) {
	return i.String(), nil
}

// Implements yaml.Unmarshaler
func (i *InstanceOp) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var output string
	if err := unmarshal(&output); err != nil {
		return fmt.Errorf("failed to unmarshal InstanceOp from YAML: %s", err)
	}
	return i.UnmarshalText([]byte(output))
}

