package snowflake

import (
	"fmt"

	pkgerrors "github.com/pkg/errors"
)

// Generate will produce a new snowflake unique id
// The lifetime is (174 years). Ref - https://github.com/sony/sonyflake/blob/848d664ceea4c980874f2135c85c42409c530b1f/sonyflake_test.go#L179
func (g *Generator) Generate() (int64, error) {
	id, err := g.flake.NextID()
	if err != nil {
		return 0, pkgerrors.WithMessage(err, "snowflake ID generation failed")
	}

	if id <= 0 {
		return 0, pkgerrors.WithStack(fmt.Errorf("snowflake ID is invalid: %d", id))
	}

	return int64(id), nil
}
