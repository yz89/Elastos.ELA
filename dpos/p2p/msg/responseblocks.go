package msg

import (
	"io"

	"github.com/elastos/Elastos.ELA/core/types"

	"github.com/elastos/Elastos.ELA.Utility/common"
)

//todo move to config
const DefaultResponseBlocksMessageDataSize = 8000000 * 10

type ResponseBlocks struct {
	Command       string
	BlockConfirms []*types.BlockConfirm
}

func (m *ResponseBlocks) CMD() string {
	return CmdResponseBlocks
}

func (m *ResponseBlocks) MaxLength() uint32 {
	return DefaultResponseBlocksMessageDataSize
}

func (m *ResponseBlocks) Serialize(w io.Writer) error {
	if err := common.WriteVarUint(w, uint64(len(m.BlockConfirms))); err != nil {
		return err
	}

	for _, v := range m.BlockConfirms {
		if err := v.Serialize(w); err != nil {
			return err
		}
	}

	return nil
}

func (m *ResponseBlocks) Deserialize(r io.Reader) error {
	blockConfirmCount, err := common.ReadVarUint(r, 0)
	if err != nil {
		return err
	}

	m.BlockConfirms = make([]*types.BlockConfirm, 0)
	for i := uint64(0); i < blockConfirmCount; i++ {
		blockConfirm := &types.BlockConfirm{}
		if err = blockConfirm.Deserialize(r); err != nil {
			return err
		}
		m.BlockConfirms = append(m.BlockConfirms, blockConfirm)
	}

	return nil
}