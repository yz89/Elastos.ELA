// Copyright (c) 2017-2019 Elastos Foundation
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.
// 

package state

// ICRRecord defines necessary operations about CR checkpoint
type ICRRecord interface {
	GetHeightsDesc() ([]uint32, error)
	GetCheckpoint(height uint32) (*CheckPoint, error)
	SaveCheckpoint(point *CheckPoint) error
}