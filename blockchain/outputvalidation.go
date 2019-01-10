package blockchain

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/elastos/Elastos.ELA/common"
	"github.com/elastos/Elastos.ELA/core/contract"
	"github.com/elastos/Elastos.ELA/core/types"
	"github.com/elastos/Elastos.ELA/core/types/outputpayload"
	"github.com/elastos/Elastos.ELA/crypto"
)

const (
	// MinDepositAmount is the minimum deposit as a producer.
	MinDepositAmount = 5000
	// DepositLockupBlocks indicates how many blocks need to wait when cancel producer was triggered, and can submit return deposit coin request.
	DepositLockupBlocks            = 2160
	MaxStringLength                = 100
	MaxVoteProducersPerTransaction = 50
)

func CheckOutputPayloadSanity(o *types.Output) error {
	// common check

	// output payload check
	switch o.OutputType {
	case types.DefaultOutput:
	case types.VoteOutput:
		return CheckOutputVoteSanity(o)
	case types.RegisterProducerOutput:
		return CheckOutputRegisterProducerSanity(o)
	default:
		return errors.New("invalid output type")
	}

	return nil
}

func CheckOutputPayloadContext(o *types.Output) error {
	// common check

	// output payload check
	switch o.OutputType {
	case types.DefaultOutput:
	case types.VoteOutput:
		return CheckOutputVoteContext(o)
	case types.RegisterProducerOutput:
		return CheckOutputRegisterProducerContext(o)
	default:
		return errors.New("invalid output type")
	}
	return nil
}

// check vote output payload
func CheckOutputVoteSanity(o *types.Output) error {
	vo, ok := o.OutputPayload.(*outputpayload.VoteOutput)
	if !ok {
		return errors.New("invalid vote payload")
	}

	typeMap := make(map[outputpayload.VoteType]struct{})
	for _, content := range vo.Contents {
		if _, exists := typeMap[content.VoteType]; exists {
			return errors.New("duplicate vote type")
		}
		typeMap[content.VoteType] = struct{}{}

		if len(content.Candidates) == 0 || len(content.Candidates) > MaxVoteProducersPerTransaction {
			return errors.New("invalid public key count")
		}
		// only use Delegate as a vote type for now
		if content.VoteType != outputpayload.Delegate {
			return errors.New("invalid vote type")
		}

		candidateMap := make(map[string]struct{})
		for _, candidate := range content.Candidates {
			c := common.BytesToHexString(candidate)
			if _, exists := candidateMap[c]; exists {
				return errors.New("duplicate candidate")
			}
			candidateMap[c] = struct{}{}
		}
	}

	return nil
}

func CheckOutputVoteContext(o *types.Output) error {
	return nil
}

// check register producer payload
func CheckOutputRegisterProducerSanity(o *types.Output) error {
	rp, ok := o.OutputPayload.(*outputpayload.RegisterProducer)
	if !ok {
		return errors.New("invalid register producer payload")
	}

	// check public key
	hash, err := contract.PublicKeyToDepositProgramHash(rp.PublicKey)
	if err != nil {
		return errors.New("invalid public key")
	}

	// check url
	if err := checkStringField(rp.Url, "Url"); err != nil {
		return err
	}

	// check ip
	if err := checkStringField(rp.Address, "IP"); err != nil {
		return err
	}

	// check signature
	publicKey, err := crypto.DecodePoint(rp.PublicKey)
	if err != nil {
		return errors.New("invalid public key in payload")
	}
	signedBuf := new(bytes.Buffer)
	err = rp.SerializeUnsigned(signedBuf)
	if err != nil {
		return err
	}
	err = crypto.Verify(*publicKey, signedBuf.Bytes(), rp.Signature)
	if err != nil {
		return errors.New("invalid signature in payload")
	}

	// check the deposit coin
	if !o.ProgramHash.IsEqual(*hash) {
		return errors.New("deposit address does not match the public key in payload")
	}
	if o.Value < MinDepositAmount {
		return errors.New("producer deposit amount is insufficient")
	}

	return nil
}

func CheckOutputRegisterProducerContext(o *types.Output) error {
	rp, ok := o.OutputPayload.(*outputpayload.RegisterProducer)
	if !ok {
		return errors.New("invalid register producer payload")
	}

	height, err := DefaultLedger.Store.GetCancelProducerHeight(rp.PublicKey)
	if err == nil {
		return fmt.Errorf("invalid producer, canceled at height: %d", height)
	}

	// check duplicate public key and nick name
	producers := DefaultLedger.Store.GetRegisteredProducers()
	if err := checkStringField(rp.NickName, "NickName"); err != nil {
		return err
	}
	for _, p := range producers {
		if bytes.Equal(p.PublicKey, rp.PublicKey) {
			return errors.New("duplicated public key")
		}
		if p.NickName == rp.NickName {
			return errors.New("duplicated nick name")
		}
	}

	return nil
}

func checkStringField(rawStr string, field string) error {
	if len(rawStr) == 0 || len(rawStr) > MaxStringLength {
		return fmt.Errorf("field %s has invalid string length", field)
	}

	return nil
}
