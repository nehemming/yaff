/*
Copyright (c) 2020-2021 The yaff Authors (Neil Hemming)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package textformatter

import (
	"github.com/nehemming/lpax"
	"github.com/nehemming/yaff/langpack"
)

// TableStyle is the style of the output table.
type TableStyle int

const (
	// Plain text.
	Plain TableStyle = iota

	// Aligned in columns.
	Aligned

	// Grid surrounding.
	Grid
)

// GetTextStyleFromString get the text style for a string.
func GetTextStyleFromString(style string) (TableStyle, error) {
	switch style {
	case "plain":
		return Plain, nil
	case "":
		fallthrough
	case "aligned":
		return Aligned, nil
	case "grid":
		return Grid, nil
	default:
		return Plain, lpax.Errorf(langpack.ErrorUnknownStyle, style)
	}
}
