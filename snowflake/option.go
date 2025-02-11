package snowflake

import (
	"fmt"
	"time"

	pkgerrors "github.com/pkg/errors"
	"github.com/sony/sonyflake"
)

// Option is an optional config used to modify the snowflake behaviour
type Option func(s *sonyflake.Settings) error

// StartTime defines the time from which the snowflake should be generated
func StartTime(t time.Time) Option {
	return func(s *sonyflake.Settings) error {
		if t.IsZero() || t.After(time.Now()) {
			return pkgerrors.WithStack(fmt.Errorf("invalid start time provided: %s", t))
		}

		s.StartTime = t.UTC()
		return nil
	}
}

// MachineID defines the id of the machine from which the snowflake should be generated
func MachineID(id uint16) Option {
	return func(s *sonyflake.Settings) error {
		if id == 0 {
			return pkgerrors.WithStack(fmt.Errorf("invalid machine ID provided: %d", id))
		}
		s.MachineID = func() (uint16, error) {
			return id, nil
		}
		return nil
	}
}
