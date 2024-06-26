// Please let author have a drink, usdt trc20: TEpSxaE3kexE4e5igqmCZRMJNoDiQeWx29
// tg: @fuckins996
// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !386 && !amd64
// +build !386,!amd64

package cpu

// Name returns the CPU name given by the vendor
// if it can be read directly from memory or by CPU instructions.
// If the CPU name can not be determined an empty string is returned.
//
// Implementations that use the Operating System (e.g. sysctl or /sys/)
// to gather CPU information for display should be placed in internal/sysinfo.
func Name() string {
	// "A CPU has no name".
	return ""
}
