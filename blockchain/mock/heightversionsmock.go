package mock

import (
	"github.com/elastos/Elastos.ELA/common"
	"github.com/elastos/Elastos.ELA/core/types"
)

type HeightVersionsMock struct {
	Producers           [][]byte
	ShouldConfirm       bool
	CurrentArbitrator   []byte
	DefaultTxVersion    byte
	DefaultBlockVersion uint32
}

func NewBlockHeightMock() *HeightVersionsMock {
	const arbitratorStr = "8a6cb4b5ff1a4f8368c6513a536c663381e3fdeff738e9b437bd8fce3fb30b62"
	arbitrator, _ := common.HexStringToBytes(arbitratorStr)

	mockObj := &HeightVersionsMock{
		Producers:           make([][]byte, 0),
		ShouldConfirm:       true,
		CurrentArbitrator:   arbitrator,
		DefaultTxVersion:    1,
		DefaultBlockVersion: 1,
	}

	return mockObj
}

func (b *HeightVersionsMock) GetCandidatesDesc(block *types.Block, startIndex uint32) ([][]byte, error) {
	return nil, nil
}

func (b *HeightVersionsMock) GetNormalArbitratorsDesc(block *types.Block, arbitratorsCount uint32) ([][]byte, error) {
	return nil, nil
}

func (b *HeightVersionsMock) GetDefaultTxVersion(blockHeight uint32) byte {
	return b.DefaultTxVersion
}

func (b *HeightVersionsMock) GetDefaultBlockVersion(blockHeight uint32) uint32 {
	return b.DefaultBlockVersion
}

func (b *HeightVersionsMock) CheckOutputPayload(blockHeight uint32, tx *types.Transaction, output *types.Output) error {
	return nil
}

func (b *HeightVersionsMock) CheckOutputProgramHash(blockHeight uint32, tx *types.Transaction, programHash common.Uint168) error {
	return nil
}

func (b *HeightVersionsMock) CheckCoinbaseMinerReward(blockHeight uint32, tx *types.Transaction, totalReward common.Fixed64) error {
	return nil
}

func (b *HeightVersionsMock) CheckCoinbaseArbitratorsReward(blockHeight uint32, coinbase *types.Transaction, rewardInCoinbase common.Fixed64) error {
	return nil
}

func (b *HeightVersionsMock) CheckVoteProducerOutputs(blockHeight uint32, tx *types.Transaction, outputs []*types.Output,
	references map[*types.Input]*types.Output, producers [][]byte) error {
	return nil
}

func (b *HeightVersionsMock) CheckTxHasNoPrograms(blockHeight uint32, tx *types.Transaction) error {
	return nil
}

func (b *HeightVersionsMock) AddBlock(block *types.Block) (bool, bool, error) {
	return true, false, nil
}

func (b *HeightVersionsMock) AddDposBlock(block *types.DposBlock) (bool, bool, error) {
	return true, false, nil
}

func (b *HeightVersionsMock) AssignCoinbaseTxRewards(block *types.Block, totalReward common.Fixed64) error {
	return nil
}

func (b *HeightVersionsMock) CheckConfirmedBlockOnFork(block *types.Block) error {
	return nil
}

func (b *HeightVersionsMock) GetNextOnDutyArbitrator(blockHeight, dutyChangedCount, offset uint32) []byte {
	return b.CurrentArbitrator
}