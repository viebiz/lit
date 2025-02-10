package snowflake

import (
	"time"

	"github.com/sony/sonyflake"
)

// Generator generates the snowflake ID
type Generator struct {
	flake *sonyflake.Sonyflake
}

// New returns a new instance of Generator
func New(opts ...Option) (*Generator, error) {
	s := sonyflake.Settings{
		StartTime: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	for _, opt := range opts {
		if err := opt(&s); err != nil {
			return nil, err
		}
	}

	return &Generator{
		flake: sonyflake.NewSonyflake(s),
	}, nil
}
