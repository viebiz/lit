package snowflake

import (
	"github.com/sony/sonyflake"
)

// Generator generates the snowflake ID
type Generator struct {
	flake *sonyflake.Sonyflake
}
