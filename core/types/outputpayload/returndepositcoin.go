package outputpayload

import "io"

type ReturnDepositCoin struct {
}

func (a *ReturnDepositCoin) Data() []byte {
	return nil
}

func (a *ReturnDepositCoin) Serialize(w io.Writer) error {
	return nil
}

func (a *ReturnDepositCoin) Deserialize(r io.Reader) error {
	return nil
}
