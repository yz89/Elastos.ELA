package outputpayload

import (
	"io"
)

type DefaultOutput struct {
}

func (do *DefaultOutput) Data() []byte {
	return nil
}

func (do *DefaultOutput) Serialize(w io.Writer) error {
	return nil
}

func (do *DefaultOutput) Deserialize(r io.Reader) error {
	return nil
}

func (do *DefaultOutput) GetVersion() byte {
	return 0
}
