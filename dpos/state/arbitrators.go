package state

import (
	"bytes"
	"encoding/hex"
	"sort"

	"github.com/elastos/Elastos.ELA/blockchain/interfaces"
	"github.com/elastos/Elastos.ELA/common"
	"github.com/elastos/Elastos.ELA/common/config"
	"github.com/elastos/Elastos.ELA/common/log"
	"github.com/elastos/Elastos.ELA/core/contract"
	"github.com/elastos/Elastos.ELA/core/types"
	"github.com/elastos/Elastos.ELA/core/types/payload"
	"github.com/elastos/Elastos.ELA/events"
)

const (
	// Numerator of dpos majority ratio
	DPOSMajorityRatioNumerator = float64(2)

	// Denominator of dpos majority ratio
	DPOSMajorityRatioDenominator = float64(3)
)

type ArbitratorsConfig struct {
	ArbitratorsCount uint32
	CandidatesCount  uint32
	CRCArbitrators   []config.CRCArbitratorParams
	Versions         interfaces.HeightVersions

	GetCurrentHeader func() (*types.Header, error)
	GetBestHeight    func() uint32
}

type Arbitrators struct {
	cfg   ArbitratorsConfig
	State *State

	dutyChangedCount uint32

	currentArbitrators [][]byte
	currentCandidates  [][]byte

	currentArbitratorsProgramHashes []*common.Uint168
	currentCandidatesProgramHashes  []*common.Uint168

	nextArbitrators [][]byte
	nextCandidates  [][]byte

	crcArbitratorsProgramHashes map[common.Uint168]interface{}

	isShadow bool
	fork     *Arbitrators
}

func (a *Arbitrators) ForwardFork() interfaces.Arbitrators {
	if a.isShadow {
		log.Error("Can not fork form shadow")
		return nil
	}

	if a.fork != nil && (a.dutyChangedCount+1)%a.cfg.ArbitratorsCount ==
		a.fork.dutyChangedCount {
		return a.fork
	}

	a.fork = &Arbitrators{
		cfg:                             a.cfg,
		dutyChangedCount:                a.dutyChangedCount,
		isShadow:                        true,
		fork:                            nil,
		crcArbitratorsProgramHashes:     a.crcArbitratorsProgramHashes,
		currentArbitrators:              make([][]byte, 0),
		currentCandidates:               make([][]byte, 0),
		currentArbitratorsProgramHashes: make([]*common.Uint168, 0),
		currentCandidatesProgramHashes:  make([]*common.Uint168, 0),
		nextArbitrators:                 make([][]byte, 0),
		nextCandidates:                  make([][]byte, 0),
	}
	a.fork.currentArbitrators = append(a.fork.currentArbitrators,
		a.currentArbitrators...)
	a.fork.currentCandidates = append(a.fork.currentCandidates,
		a.currentCandidates...)
	a.fork.currentArbitratorsProgramHashes = append(a.fork.
		currentArbitratorsProgramHashes, a.currentArbitratorsProgramHashes...)
	a.fork.currentCandidatesProgramHashes = append(a.fork.
		currentCandidatesProgramHashes, a.currentCandidatesProgramHashes...)
	a.fork.nextArbitrators = append(a.fork.nextArbitrators,
		a.nextArbitrators...)
	a.fork.nextCandidates = append(a.fork.nextCandidates, a.nextCandidates...)

	a.fork.IncreaseChainHeight(nil)

	return a.fork
}

func (a *Arbitrators) ForceChange() error {
	block, err := a.cfg.GetCurrentHeader()
	if err != nil {
		return err
	}

	if err = a.updateNextArbitrators(block); err != nil {
		return err
	}

	if err = a.changeCurrentArbitrators(); err != nil {
		return err
	}

	events.Notify(events.ETNewArbiterElection, a.nextArbitrators)

	return nil
}

func (a *Arbitrators) IncreaseChainHeight(block *types.Block) {

	if a.isNewElection(block) {
		if err := a.changeCurrentArbitrators(); err != nil {
			log.Error("Change current arbitrators error: ", err)
			return
		}

		if a.isShadow {
			if err := a.updateNextArbitrators(&block.Header); err != nil {
				log.Error("Update arbitrators error: ", err)
				return
			}

			events.Notify(events.ETNewArbiterElection, a.nextArbitrators)
		}

	} else {
		a.dutyChangedCount++
	}
}

func (a *Arbitrators) DecreaseChainHeight(block *types.Block) {

	if a.dutyChangedCount == 0 {
		//todo complete me
	} else {
		a.dutyChangedCount--
	}
}

func (a *Arbitrators) IsArbitrator(pk []byte) bool {
	arbitrators := a.GetArbitrators()

	for _, v := range arbitrators {
		if bytes.Equal(pk, v) {
			return true
		}
	}
	return false
}

func (a *Arbitrators) GetArbitrators() [][]byte {
	a.State.mtx.RLock()
	result := a.currentArbitrators
	a.State.mtx.RUnlock()

	return result
}

func (a *Arbitrators) GetCandidates() [][]byte {
	a.State.mtx.RLock()
	result := a.currentCandidates
	a.State.mtx.RUnlock()

	return result
}

func (a *Arbitrators) GetNextArbitrators() [][]byte {
	a.State.mtx.RLock()
	result := a.nextArbitrators
	a.State.mtx.RUnlock()

	return result
}

func (a *Arbitrators) GetNextCandidates() [][]byte {
	a.State.mtx.RLock()
	result := a.nextCandidates
	a.State.mtx.RUnlock()

	return result
}

func (a *Arbitrators) IsCRCArbitratorProgramHash(hash *common.Uint168) bool {
	_, ok := a.crcArbitratorsProgramHashes[*hash]
	return ok
}

func (a *Arbitrators) IsCRCArbitrator(pk []byte) bool {
	for _, v := range a.cfg.CRCArbitrators {
		if bytes.Equal(v.PublicKey, pk) {
			return true
		}
	}
	return false
}

func (a *Arbitrators) GetCRCArbitrators() []config.CRCArbitratorParams {
	return a.cfg.CRCArbitrators
}

func (a *Arbitrators) GetArbitratorsProgramHashes() []*common.Uint168 {
	a.State.mtx.RLock()
	result := a.currentArbitratorsProgramHashes
	a.State.mtx.RUnlock()

	return result
}

func (a *Arbitrators) GetCandidatesProgramHashes() []*common.Uint168 {
	a.State.mtx.RLock()
	result := a.currentCandidatesProgramHashes
	a.State.mtx.RUnlock()

	return result
}

func (a *Arbitrators) GetOnDutyArbitrator() []byte {
	return a.cfg.Versions.GetNextOnDutyArbitrator(a.cfg.GetBestHeight(),
		a.dutyChangedCount, 0)
}

func (a *Arbitrators) GetNextOnDutyArbitrator(offset uint32) []byte {
	return a.cfg.Versions.GetNextOnDutyArbitrator(a.cfg.GetBestHeight()+1,
		a.dutyChangedCount, offset)
}

func (a *Arbitrators) GetArbitersCount() uint32 {
	a.State.mtx.RLock()
	result := a.getArbitersCount()
	a.State.mtx.RUnlock()

	return result
}

func (a *Arbitrators) GetArbitersMajorityCount() uint32 {
	a.State.mtx.RLock()
	minSignCount := float64(a.getArbitersCount()) *
		DPOSMajorityRatioNumerator / DPOSMajorityRatioDenominator
	a.State.mtx.RUnlock()

	return uint32(minSignCount)
}

func (a *Arbitrators) HasArbitersMajorityCount(num uint32) bool {
	return num > a.GetArbitersMajorityCount()
}

func (a *Arbitrators) HasArbitersMinorityCount(num uint32) bool {
	return num >= a.cfg.ArbitratorsCount-a.GetArbitersMajorityCount()
}

// getInactiveArbitrators returns inactive arbiters from a confirm data.
func (a *Arbitrators) GetInactiveArbitrators(confirm *payload.Confirm,
	onDutyArbitrator []byte) (result []string) {

	// check inactive arbitrators after producers has participated in
	if len(a.currentArbitrators) == len(a.cfg.CRCArbitrators) {
		return
	}

	if !bytes.Equal(onDutyArbitrator, confirm.Proposal.Sponsor) {
		arSequence := a.currentArbitrators
		arSequence = append(arSequence, arSequence...)

		start := onDutyArbitrator
		stop := confirm.Proposal.Sponsor
		reachedStart := false

		for i := 0; i < len(arSequence)-1; i++ {
			if bytes.Equal(start, arSequence[i]) {
				reachedStart = true
			}

			if reachedStart {
				if bytes.Equal(stop, arSequence[i]) {
					break
				}

				result = append(result, hex.EncodeToString(arSequence[i]))
			}
		}
	}

	return
}

func (a *Arbitrators) isNewElection(block *types.Block) bool {
	if a.cfg.Versions.GetDefaultBlockVersion(block.Height) >= 1 {

		// when change to "H1" height should fire new election immediately
		if block.Height == a.State.chainParams.HeightVersions[2] {
			return true
		}
		return a.dutyChangedCount == a.cfg.ArbitratorsCount-1
	}
	return false
}

func (a *Arbitrators) changeCurrentArbitrators() error {
	a.currentArbitrators = a.nextArbitrators
	a.currentCandidates = a.nextCandidates

	if err := a.sortArbitrators(); err != nil {
		return err
	}

	if err := a.updateArbitratorsProgramHashes(); err != nil {
		return err
	}

	a.dutyChangedCount = 0

	return nil
}

func (a *Arbitrators) updateNextArbitrators(header *types.Header) error {

	crcCount := uint32(0)
	a.nextArbitrators = make([][]byte, 0)
	for _, v := range a.cfg.CRCArbitrators {
		if !a.State.IsInactiveProducer(v.PublicKey) {
			a.nextArbitrators = append(a.nextArbitrators, v.PublicKey)
			crcCount++
		}
	}
	count := config.Parameters.ArbiterConfiguration.
		NormalArbitratorsCount + crcCount
	producers, err := a.cfg.Versions.GetNormalArbitratorsDesc(header.Height, count)
	if err != nil {
		return err
	}
	for _, v := range producers {
		a.nextArbitrators = append(a.nextArbitrators, v)
	}

	candidates, err := a.cfg.Versions.GetCandidatesDesc(header.Height, count)
	if err != nil {
		return err
	}
	a.nextCandidates = candidates

	return nil
}

func (a *Arbitrators) sortArbitrators() error {

	strArbitrators := make([]string, len(a.currentArbitrators))
	for i := 0; i < len(strArbitrators); i++ {
		strArbitrators[i] = common.BytesToHexString(a.currentArbitrators[i])
	}
	sort.Strings(strArbitrators)

	a.currentArbitrators = make([][]byte, len(strArbitrators))
	for i := 0; i < len(strArbitrators); i++ {
		value, err := common.HexStringToBytes(strArbitrators[i])
		if err != nil {
			return err
		}
		a.currentArbitrators[i] = value
	}

	return nil
}

func (a *Arbitrators) getArbitersCount() uint32 {
	return uint32(len(a.currentArbitrators))
}

func (a *Arbitrators) updateArbitratorsProgramHashes() error {
	a.currentArbitratorsProgramHashes = make([]*common.Uint168, len(a.currentArbitrators))
	for index, v := range a.currentArbitrators {
		hash, err := contract.PublicKeyToStandardProgramHash(v)
		if err != nil {
			return err
		}
		a.currentArbitratorsProgramHashes[index] = hash
	}

	a.currentCandidatesProgramHashes = make([]*common.Uint168, len(a.currentCandidates))
	for index, v := range a.currentCandidates {
		hash, err := contract.PublicKeyToStandardProgramHash(v)
		if err != nil {
			return err
		}
		a.currentCandidatesProgramHashes[index] = hash
	}

	return nil
}

func NewArbitrators(cfg *ArbitratorsConfig) (*Arbitrators, error) {

	a := &Arbitrators{
		cfg:             *cfg,
		isShadow:        false,
		nextArbitrators: make([][]byte, 0),
		nextCandidates:  make([][]byte, 0),
	}

	a.crcArbitratorsProgramHashes = make(map[common.Uint168]interface{})
	for _, v := range a.cfg.CRCArbitrators {
		a.nextArbitrators = append(a.nextArbitrators, v.PublicKey)

		hash, err := contract.PublicKeyToStandardProgramHash(v.PublicKey)
		if err != nil {
			return nil, err
		}
		a.crcArbitratorsProgramHashes[*hash] = nil
	}

	return a, nil
}