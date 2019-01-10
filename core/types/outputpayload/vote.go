package outputpayload

import (
	"fmt"
	"io"

	"github.com/elastos/Elastos.ELA/common"
	"github.com/elastos/Elastos.ELA/crypto"
)

const (
	Delegate VoteType = 0x00
	CRC      VoteType = 0x01
)

type VoteType byte

var VoteTypes = []VoteType{
	Delegate,
	CRC,
}

type VoteContent struct {
	VoteType   VoteType
	Candidates [][]byte
}

func (vc *VoteContent) Serialize(w io.Writer, version byte) error {
	if _, err := w.Write([]byte{byte(vc.VoteType)}); err != nil {
		return err
	}
	if err := common.WriteVarUint(w, uint64(len(vc.Candidates))); err != nil {
		return err
	}
	for _, candidate := range vc.Candidates {
		if err := common.WriteVarBytes(w, candidate); err != nil {
			return err
		}
	}

	return nil
}

func (vc *VoteContent) Deserialize(r io.Reader, version byte) error {
	voteType, err := common.ReadBytes(r, 1)
	if err != nil {
		return err
	}
	vc.VoteType = VoteType(voteType[0])

	candidatesCount, err := common.ReadVarUint(r, 0)
	if err != nil {
		return err
	}

	for i := uint64(0); i < candidatesCount; i++ {
		candidate, err := common.ReadVarBytes(r, crypto.COMPRESSEDLEN, "producer")
		if err != nil {
			return err
		}
		vc.Candidates = append(vc.Candidates, candidate)
	}

	return nil
}

func (vc VoteContent) String() string {
	return fmt.Sprint("Content: {",
		"VoteType: ", vc.VoteType, " ",
		"Candidates: ", vc.Candidates, "}")
}

type VoteOutput struct {
	Version  byte
	Contents []VoteContent
}

func (vo *VoteOutput) Data() []byte {
	return nil
}

func (vo *VoteOutput) Serialize(w io.Writer) error {
	if _, err := w.Write([]byte{byte(vo.Version)}); err != nil {
		return err
	}
	if err := common.WriteVarUint(w, uint64(len(vo.Contents))); err != nil {
		return err
	}
	for _, content := range vo.Contents {
		if err := content.Serialize(w, vo.Version); err != nil {
			return err
		}
	}
	return nil
}

func (vo *VoteOutput) Deserialize(r io.Reader) error {
	version, err := common.ReadBytes(r, 1)
	if err != nil {
		return err
	}
	vo.Version = version[0]

	contentsCount, err := common.ReadVarUint(r, 0)
	if err != nil {
		return err
	}

	for i := uint64(0); i < contentsCount; i++ {
		var content VoteContent
		if err := content.Deserialize(r, vo.Version); err != nil {
			return err
		}
		vo.Contents = append(vo.Contents, content)
	}

	return nil
}

func (vo *VoteOutput) GetVersion() byte {
	return vo.Version
}

func (vo VoteOutput) String() string {
	return fmt.Sprint("Vote: {",
		"Version: ", vo.Version, " ",
		"Contents: ", vo.Contents, "\n\t}")
}
