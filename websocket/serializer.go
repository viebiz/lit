package websocket

import (
	"encoding/json"
	"io"

	pkgerrors "github.com/pkg/errors"
)

type jsonSerializer struct{}

func NewJSONSerializer() Serializer {
	return &jsonSerializer{}
}

func (j jsonSerializer) Serialize(o any, w io.Writer) error {
	if err := json.NewEncoder(w).Encode(o); err != nil {
		return pkgerrors.WithStack(err)
	}
	return nil
}

func (j jsonSerializer) Deserialize(r io.Reader, o any) error {
	if err := json.NewDecoder(r).Decode(o); err != nil {
		return pkgerrors.WithStack(err)
	}

	return nil
}
