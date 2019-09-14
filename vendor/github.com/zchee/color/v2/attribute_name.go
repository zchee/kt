// Copyright 2019 The color Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Code generated by "stringer -type=Attribute"; DO NOT EDIT.

package color

import "strconv"

const (
	_Attribute_name_0 = "ResetBoldFaintItalicUnderlineBlinkSlowBlinkRapidReverseVideoConcealedCrossedOut"
	_Attribute_name_1 = "FgBlackFgRedFgGreenFgYellowFgBlueFgMagentaFgCyanFgWhite"
	_Attribute_name_2 = "BgBlackBgRedBgGreenBgYellowBgBlueBgMagentaBgCyanBgWhite"
	_Attribute_name_3 = "FgHiBlackFgHiRedFgHiGreenFgHiYellowFgHiBlueFgHiMagentaFgHiCyanFgHiWhite"
	_Attribute_name_4 = "BgHiBlackBgHiRedBgHiGreenBgHiYellowBgHiBlueBgHiMagentaBgHiCyanBgHiWhite"
)

var (
	_Attribute_index_0 = [...]uint8{0, 5, 9, 14, 20, 29, 38, 48, 60, 69, 79}
	_Attribute_index_1 = [...]uint8{0, 7, 12, 19, 27, 33, 42, 48, 55}
	_Attribute_index_2 = [...]uint8{0, 7, 12, 19, 27, 33, 42, 48, 55}
	_Attribute_index_3 = [...]uint8{0, 9, 16, 25, 35, 43, 54, 62, 71}
	_Attribute_index_4 = [...]uint8{0, 9, 16, 25, 35, 43, 54, 62, 71}
)

// Name name of Attribute.
func (i Attribute) Name() string {
	switch {
	case 0 <= i && i <= 9:
		return _Attribute_name_0[_Attribute_index_0[i]:_Attribute_index_0[i+1]]
	case 30 <= i && i <= 37:
		i -= 30
		return _Attribute_name_1[_Attribute_index_1[i]:_Attribute_index_1[i+1]]
	case 40 <= i && i <= 47:
		i -= 40
		return _Attribute_name_2[_Attribute_index_2[i]:_Attribute_index_2[i+1]]
	case 90 <= i && i <= 97:
		i -= 90
		return _Attribute_name_3[_Attribute_index_3[i]:_Attribute_index_3[i+1]]
	case 100 <= i && i <= 107:
		i -= 100
		return _Attribute_name_4[_Attribute_index_4[i]:_Attribute_index_4[i+1]]
	default:
		return "Attribute(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}