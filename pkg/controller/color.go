// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package controller

import (
	xxhash "github.com/cespare/xxhash/v2"
	color "github.com/zchee/color/v2"
)

var colorList = [][2]*color.Color{
	{color.New(color.FgHiCyan), color.New(color.FgCyan)},
	{color.New(color.FgHiGreen), color.New(color.FgGreen)},
	{color.New(color.FgHiMagenta), color.New(color.FgMagenta)},
	{color.New(color.FgHiYellow), color.New(color.FgYellow)},
	{color.New(color.FgHiBlue), color.New(color.FgBlue)},
	{color.New(color.FgHiRed), color.New(color.FgRed)},
	{color.New(color.FgHiCyan, color.Faint), color.New(color.FgCyan, color.Faint)},
	{color.New(color.FgHiGreen, color.Faint), color.New(color.FgGreen, color.Faint)},
	{color.New(color.FgHiMagenta, color.Faint), color.New(color.FgMagenta, color.Faint)},
	{color.New(color.FgHiYellow, color.Faint), color.New(color.FgYellow, color.Faint)},
	{color.New(color.FgHiBlue, color.Faint), color.New(color.FgBlue, color.Faint)},
	{color.New(color.FgHiRed, color.Faint), color.New(color.FgRed, color.Faint)},
}

func findColors(podName string) (podColor, containerColor *color.Color) {
	digest := xxhash.New()
	digest.WriteString(podName)
	idx := digest.Sum64() % uint64(len(colorList))

	colors := colorList[idx]
	return colors[0], colors[1]
}