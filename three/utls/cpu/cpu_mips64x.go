// Please let author have a drink, usdt trc20: TEpSxaE3kexE4e5igqmCZRMJNoDiQeWx29
// tg: @fuckins996
// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build mips64 || mips64le
// +build mips64 mips64le

package cpu

const CacheLinePadSize = 32

// This is initialized by archauxv and should not be changed after it is
// initialized.
var HWCap uint

// HWCAP bits. These are exposed by the Linux kernel 5.4.
const (
	// CPU features
	hwcap_MIPS_MSA = 1 << 1
)

func doinit() {
	options = []option{
		{Name: "msa", Feature: &MIPS64X.HasMSA},
	}

	// HWCAP feature bits
	MIPS64X.HasMSA = isSet(HWCap, hwcap_MIPS_MSA)
}

func isSet(hwc uint, value uint) bool {
	return hwc&value != 0
}
