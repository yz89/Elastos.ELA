package outputpayload

import (
	"bytes"
	"errors"
	"io"

	"github.com/elastos/Elastos.ELA/common"
	"github.com/elastos/Elastos.ELA/crypto"
)

type UpdateProducer struct {
	PublicKey []byte
	NickName  string
	Url       string
	Location  uint64
	Address   string
	Signature []byte
}

func (up *UpdateProducer) Data() []byte {
	buf := new(bytes.Buffer)
	if err := up.Serialize(buf); err != nil {
		return []byte{0}
	}
	return buf.Bytes()
}

func (up *UpdateProducer) Serialize(w io.Writer) error {
	err := up.SerializeUnsigned(w)
	if err != nil {
		return err
	}

	err = common.WriteVarBytes(w, up.Signature)
	if err != nil {
		return errors.New("[UpdateProducer], Signature serialize failed")
	}

	return nil
}

func (up *UpdateProducer) SerializeUnsigned(w io.Writer) error {
	err := common.WriteVarBytes(w, up.PublicKey)
	if err != nil {
		return errors.New("[UpdateProducer], PublicKey serialize failed")
	}

	err = common.WriteVarString(w, up.NickName)
	if err != nil {
		return errors.New("[UpdateProducer], NickName serialize failed")
	}

	err = common.WriteVarString(w, up.Url)
	if err != nil {
		return errors.New("[UpdateProducer], Url serialize failed")
	}

	err = common.WriteUint64(w, up.Location)
	if err != nil {
		return errors.New("[UpdateProducer], Location serialize failed")
	}

	err = common.WriteVarString(w, up.Address)
	if err != nil {
		return errors.New("[UpdateProducer], Address serialize failed")
	}
	return nil
}

func (up *UpdateProducer) Deserialize(r io.Reader) error {
	err := up.DeserializeUnsigned(r)
	if err != nil {
		return err
	}
	sig, err := common.ReadVarBytes(r, crypto.SignatureLength, "signature")
	if err != nil {
		return errors.New("[UpdateProducer], Signature deserialize failed")
	}

	up.Signature = sig

	return nil
}

func (up *UpdateProducer) DeserializeUnsigned(r io.Reader) error {
	publicKey, err := common.ReadVarBytes(r, crypto.NegativeBigLength, "public key")
	if err != nil {
		return errors.New("[UpdateProducer], PublicKey deserialize failed")
	}

	nickName, err := common.ReadVarString(r)
	if err != nil {
		return errors.New("[UpdateProducer], NickName deserialize failed")
	}

	url, err := common.ReadVarString(r)
	if err != nil {
		return errors.New("[UpdateProducer], Url deserialize failed")
	}

	location, err := common.ReadUint64(r)
	if err != nil {
		return errors.New("[UpdateProducer], Location deserialize failed")
	}

	addr, err := common.ReadVarString(r)
	if err != nil {
		return errors.New("[UpdateProducer], Address deserialize failed")
	}

	up.PublicKey = publicKey
	up.NickName = nickName
	up.Url = url
	up.Location = location
	up.Address = addr

	return nil
}

func (up *UpdateProducer) GetVersion() byte {
	return 0
}

func ConvertToRegisterProducerPayload(update *UpdateProducer) *RegisterProducer {
	return &RegisterProducer{
		PublicKey: update.PublicKey,
		NickName:  update.NickName,
		Url:       update.Url,
		Location:  update.Location,
		Address:   update.Address,
	}
}
