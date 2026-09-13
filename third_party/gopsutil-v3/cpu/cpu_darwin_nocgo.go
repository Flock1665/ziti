// PATCHED (local fork): also selected on iOS, where the cgo variant above is
// excluded; see ../../PATCHES.md
//go:build darwin && (!cgo || ios)
// +build darwin,!cgo ios

package cpu

import "github.com/shirou/gopsutil/v3/internal/common"

func perCPUTimes() ([]TimesStat, error) {
	return []TimesStat{}, common.ErrNotImplementedError
}

func allCPUTimes() ([]TimesStat, error) {
	return []TimesStat{}, common.ErrNotImplementedError
}
