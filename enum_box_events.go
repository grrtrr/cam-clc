package clccam

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// Events in the lifecycle of an Elastic Box
type BoxEvent uint32

const (
	// Event before configuring the box
	BoxEvent_PreConfigure BoxEvent = iota

	// Configuring the box
	BoxEvent_Configure

	// Before installing onto the box
	BoxEvent_PreInstall

	// Installing onto the box
	BoxEvent_Install

	// Before starting the box
	BoxEvent_PreStart

	// Starting the box
	BoxEvent_Start

	// Before stopping the box
	BoxEvent_PreStop

	// Stopping the box
	BoxEvent_Stop

	// Before disposing the box
	BoxEvent_PreDispose

	// Disposing the box
	BoxEvent_Dispose
)

// Implements encoding.TextMarshaler
func (b BoxEvent) MarshalText() ([]byte, error) {
	switch b {
	case BoxEvent_PreConfigure:
		return []byte("pre_configure"), nil
	case BoxEvent_Configure:
		return []byte("configure"), nil
	case BoxEvent_PreInstall:
		return []byte("pre_install"), nil
	case BoxEvent_Install:
		return []byte("install"), nil
	case BoxEvent_PreStart:
		return []byte("pre_start"), nil
	case BoxEvent_Start:
		return []byte("start"), nil
	case BoxEvent_PreStop:
		return []byte("pre_stop"), nil
	case BoxEvent_Stop:
		return []byte("stop"), nil
	case BoxEvent_PreDispose:
		return []byte("pre_dispose"), nil
	case BoxEvent_Dispose:
		return []byte("dispose"), nil
	}
	return nil, fmt.Errorf("invalid BoxEvent %d", b)
}

// Implements fmt.Stringer
func (b BoxEvent) String() string {
	if b, err := b.MarshalText(); err != nil {
		return err.Error()
	} else {
		return string(b)
	}
}

// Implements encoding.TextUnmarshaler
func (b *BoxEvent) UnmarshalText(data []byte) error {
	switch string(data) {
	case "pre_configure":
		*b = BoxEvent_PreConfigure
	case "configure":
		*b = BoxEvent_Configure
	case "pre_install":
		*b = BoxEvent_PreInstall
	case "install":
		*b = BoxEvent_Install
	case "pre_start":
		*b = BoxEvent_PreStart
	case "start":
		*b = BoxEvent_Start
	case "pre_stop":
		*b = BoxEvent_PreStop
	case "stop":
		*b = BoxEvent_Stop
	case "pre_dispose":
		*b = BoxEvent_PreDispose
	case "dispose":
		*b = BoxEvent_Dispose
	default:
		return fmt.Errorf("invalid BoxEvent %q", string(data))
	}
	return nil
}

// Implements flag.Value
func (b *BoxEvent) Set(s string) error {
	return b.UnmarshalText([]byte(s))
}

// Implements pflag.Value (superset of flag.Value)
func (b BoxEvent) Type() string {
	return "BoxEvent"
}

// BoxEventFromString attempts to parse @s as stringified BoxEvent.
func BoxEventFromString(s string) (val BoxEvent, err error) {
	return val, val.Set(s)
}

// BoxEventStrings returns the list of BoxEvent string literals, or maps @vals if non-empty.
func BoxEventStrings(vals ...BoxEvent) (ret []string) {
	if len(vals) > 0 {
		for _, val := range vals {
			ret = append(ret, val.String())
		}
		return ret
	}
	return []string{"pre_configure", "configure", "pre_install", "install", "pre_start", "start", "pre_stop", "stop", "pre_dispose", "dispose"}
}

// Implements database/sql/driver.Valuer
func (b BoxEvent) Value() (driver.Value, error) {
	return b.String(), nil
}

// Implements database/sql.Scanner
func (b *BoxEvent) Scan(src interface{}) error {
	switch src := src.(type) {
	case int64:
		*b = BoxEvent(src)
		return nil
	case []byte:
		return b.UnmarshalText(src)
	case string:
		return b.UnmarshalText([]byte(src))
	}
	return fmt.Errorf("unable to convert %T to BoxEvent", src)
}

// Implements json.Marshaler
func (b BoxEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}

// Implements json.Unmarshaler
func (b *BoxEvent) UnmarshalJSON(data []byte) error {
	var output string
	if err := json.Unmarshal(data, &output); err != nil {
		return fmt.Errorf("failed to parse '%s' as BoxEvent: %s", string(data), err)
	}
	return b.UnmarshalText([]byte(output))
}

// Implements yaml.Marshaler
func (b BoxEvent) MarshalYAML() (interface{}, error) {
	return b.String(), nil
}

// Implements yaml.Unmarshaler
func (b *BoxEvent) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var output string
	if err := unmarshal(&output); err != nil {
		return fmt.Errorf("failed to unmarshal BoxEvent from YAML: %s", err)
	}
	return b.UnmarshalText([]byte(output))
}

