package outputpayload

import (
	"bytes"
	"errors"
	"io"

	"github.com/elastos/Elastos.ELA/common"
	"github.com/elastos/Elastos.ELA/crypto"
)

type RegisterProducer struct {
	PublicKey []byte
	NickName  string
	Url       string
	Location  uint64
	Address   string
	Signature []byte
}

func (rp *RegisterProducer) Data() []byte {
	buf := new(bytes.Buffer)
	if err := rp.Serialize(buf); err != nil {
		return []byte{0}
	}
	return buf.Bytes()
}

func (rp *RegisterProducer) Serialize(w io.Writer) error {
	err := rp.SerializeUnsigned(w)
	if err != nil {
		return err
	}

	err = common.WriteVarBytes(w, rp.Signature)
	if err != nil {
		return errors.New("[RegisterProducer], Signature serialize failed")
	}

	return nil
}

func (rp *RegisterProducer) SerializeUnsigned(w io.Writer) error {
	err := common.WriteVarBytes(w, rp.PublicKey)
	if err != nil {
		return errors.New("[RegisterProducer], PublicKey serialize failed")
	}

	err = common.WriteVarString(w, rp.NickName)
	if err != nil {
		return errors.New("[RegisterProducer], NickName serialize failed")
	}

	err = common.WriteVarString(w, rp.Url)
	if err != nil {
		return errors.New("[RegisterProducer], Url serialize failed")
	}

	err = common.WriteUint64(w, rp.Location)
	if err != nil {
		return errors.New("[RegisterProducer], Location serialize failed")
	}

	err = common.WriteVarString(w, rp.Address)
	if err != nil {
		return errors.New("[RegisterProducer], Address serialize failed")
	}
	return nil
}

func (rp *RegisterProducer) Deserialize(r io.Reader) error {
	err := rp.DeserializeUnsigned(r)
	if err != nil {
		return err
	}
	sig, err := common.ReadVarBytes(r, crypto.SignatureLength, "signature")
	if err != nil {
		return errors.New("[RegisterProducer], Signature deserialize failed")
	}

	rp.Signature = sig

	return nil
}

func (rp *RegisterProducer) DeserializeUnsigned(r io.Reader) error {
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

	rp.PublicKey = publicKey
	rp.NickName = nickName
	rp.Url = url
	rp.Location = location
	rp.Address = addr

	return nil
}

func (rp *RegisterProducer) GetVersion() byte {
	return 0
}
