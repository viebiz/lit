package redis

type Type interface {
	string | int64 | uint64 | float64
}

type setMode string

const (
	setModeNone setMode = ""

	setModeNX setMode = "NX"

	setModeXX setMode = "XX"
)

func (s setMode) String() string {
	return string(s)
}

func (s setMode) IsValid() bool {
	return s == setModeNX || s == setModeXX
}
