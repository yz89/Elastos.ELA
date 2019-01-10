package outputpayload

import (
	"bytes"
	"errors"
	"io"

	"github.com/elastos/Elastos.ELA/common"
	"github.com/elastos/Elastos.ELA/crypto"
)

type CancelProducer struct {
	PublicKey []byte
	Signature []byte
}

func (cp *CancelProducer) Data() []byte {
	buf := new(bytes.Buffer)
	if err := cp.Serialize(buf); err != nil {
		return []byte{0}
	}
	return buf.Bytes()
}

func (cp *CancelProducer) Serialize(w io.Writer) error {
	err := cp.SerializeUnsigned(w)
	if err != nil {
		return err
	}

	err = common.WriteVarBytes(w, cp.Signature)
	if err != nil {
		return errors.New("[CancelProducer], Signature serialize failed")
	}

	return nil
}

func (cp *CancelProducer) SerializeUnsigned(w io.Writer) error {
	err := common.WriteVarBytes(w, cp.PublicKey)
	if err != nil {
		return errors.New("[CancelProducer], Serialize failed")
	}
	return nil
}

func (cp *CancelProducer) Deserialize(r io.Reader) error {
	err := cp.DeserializeUnsigned(r)
	if err != nil {
		return err
	}
	sig, err := common.ReadVarBytes(r, crypto.SignatureLength, "signature")
	if err != nil {
		return errors.New("[CancelProducer], Signature deserialize failed")
	}

	cp.Signature = sig

	return nil
}

func (cp *CancelProducer) DeserializeUnsigned(r io.Reader) error {
	pk, err := common.ReadVarBytes(r, crypto.NegativeBigLength, "public key")
	if err != nil {
		return errors.New("[CancelProducer], Deserialize failed")
	}
	cp.PublicKey = pk
	return err
}

func (cp *CancelProducer) GetVersion() byte {
	return 0
}
