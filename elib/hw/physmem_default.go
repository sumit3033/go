// Copyright 2016 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !uio_pci_dma,!uio_pci_generic,novfio

package hw

import (
	"github.com/platinasystems/go/elib"
)

func DmaAllocAligned(n, log2Align uint) (b []byte, id elib.Index, offset, cap uint) {
	return heap.GetAligned(n, log2Align)
}
func DmaPhysAddress(a uintptr) uintptr { return a }
