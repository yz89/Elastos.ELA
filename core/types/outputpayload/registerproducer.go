package outputpayload

import (
	"bytes"
	"errors"
	"io"

	"github.com/elastos/Elastos.ELA/common"
	"github.com/elastos/Elastos.ELA/crypto"
)

const PayloadRegisterProducerVersion byte = 0x00

type RegisterProducer struct {
	PublicKey []byte
	NickName  string
	Url       string
	Location  uint64
	Address   string
	Signature []byte
}

func (a *RegisterProducer) Data() []byte {
	buf := new(bytes.Buffer)
	if err := a.Serialize(buf); err != nil {
		return []byte{0}
	}
	return buf.Bytes()
}

func (a *RegisterProducer) Serialize(w io.Writer) error {
	err := a.SerializeUnsigned(w)
	if err != nil {
		return err
	}

	err = common.WriteVarBytes(w, a.Signature)
	if err != nil {
		return errors.New("[RegisterProducer], Signature serialize failed")
	}

	return nil
}

func (a *RegisterProducer) SerializeUnsigned(w io.Writer) error {
	err := common.WriteVarBytes(w, a.PublicKey)
	if err != nil {
		return errors.New("[RegisterProducer], PublicKey serialize failed")
	}

	err = common.WriteVarString(w, a.NickName)
	if err != nil {
		return errors.New("[RegisterProducer], NickName serialize failed")
	}

	err = common.WriteVarString(w, a.Url)
	if err != nil {
		return errors.New("[RegisterProducer], Url serialize failed")
	}

	err = common.WriteUint64(w, a.Location)
	if err != nil {
		return errors.New("[RegisterProducer], Location serialize failed")
	}

	err = common.WriteVarString(w, a.Address)
	if err != nil {
		return errors.New("[RegisterProducer], Address serialize failed")
	}
	return nil
}

func (a *RegisterProducer) Deserialize(r io.Reader) error {
	err := a.DeserializeUnsigned(r)
	if err != nil {
		return err
	}
	sig, err := common.ReadVarBytes(r, crypto.SignatureLength, "signature")
	if err != nil {
		return errors.New("[RegisterProducer], Signature deserialize failed")
	}

	a.Signature = sig

	return nil
}

func (a *RegisterProducer) DeserializeUnsigned(r io.Reader) error {
	publicKey, err := common.ReadVarBytes(r, crypto.NegativeBigLength, "public key")
	if err != nil {
		return errors.New("[RegisterProducer], PublicKey deserialize failed")
	}

	nickName, err := common.ReadVarString(r)
	if err != nil {
		return errors.New("[RegisterProducer], NickName deserialize failed")
	}

	url, err := common.ReadVarString(r)
	if err != nil {
		return errors.New("[RegisterProducer], Url deserialize failed")
	}

	location, err := common.ReadUint64(r)
	if err != nil {
		return errors.New("[RegisterProducer], Location deserialize failed")
	}

	addr, err := common.ReadVarString(r)
	if err != nil {
		return errors.New("[RegisterProducer], Address deserialize failed")
	}

	a.PublicKey = publicKey
	a.NickName = nickName
	a.Url = url
	a.Location = location
	a.Address = addr

	return nil
}

func (a *RegisterProducer) GetVersion() byte {
	return 0
}

func (a *RegisterProducer) Validate() error {
	return nil
}
