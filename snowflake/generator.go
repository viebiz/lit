package snowflake

// Generator generates the snowflake ID
type Generator struct {
	flake idProvider
}

// idProvider defines an entity that can provide the next unique ID.
type idProvider interface {
	NextID() (uint64, error)
}
