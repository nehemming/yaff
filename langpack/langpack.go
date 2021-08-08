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

// Package langpack for yaff.
package langpack

import (
	"github.com/nehemming/lpax"
	"golang.org/x/text/language"
)

type (
	// PackID is the local package id.
	PackID int

	// TextID is the id type for all messages defined by this package.
	TextID int
)

// Single returns the id of for a single version of a message.
// Typical implementations use the lpax helper functions to implement the id system.
// Positive id's are used for single messages while -ve id's areused for plural versions.
func (id TextID) Single() lpax.TextID {
	return TextID(lpax.IntTypeSingle(int(id)))
}

// Plural returns the plural version of the id.
func (id TextID) Plural() lpax.TextID {
	return TextID(lpax.IntTypePlural(int(id)))
}

// String implements stringer function.
// ReflectCoderString reflects the package name used
// by the parent package along with the +ve integer id of the key.
func (id TextID) String() string {
	return lpax.ReflectCoderString(id.Single(), 1)
}

const (
	// LangPack language pack for this package.
	LangPack = PackID(1)

	// None is a default empty message (id 0).
	None = TextID(iota)

	// ErrorUnknownFormatter Unknown formatter.
	ErrorUnknownFormatter

	// ErrorInvalidOptionType unknown option type.
	ErrorInvalidOptionType

	// ErrorNoTemplateDefinition no template specified error.
	ErrorNoTemplateDefinition

	// ErrorUnknownStyle unknown text style.
	ErrorUnknownStyle

	// ErrorColumnInvalidID Column has an invalid ID.
	ErrorColumnInvalidID

	// ErrorRowInvalidID Row has an invalid ID.
	ErrorRowInvalidID
)

var languagePack = lpax.TextMap{

	ErrorUnknownFormatter:     "Unknown formatter %s",
	ErrorInvalidOptionType:    "Option type %[1]T is not valid for formatter %[2]v",
	ErrorNoTemplateDefinition: "No template definition specified",
	ErrorUnknownStyle:         "Unknown style %v",

	ErrorColumnInvalidID: "Column %d is an invalid id",
	ErrorRowInvalidID:    "Row %d is an invalid id",
}

func init() {
	// Register formats and english language pack.
	lpax.Default().Register(LangPack, func(packID lpax.PackID, langTag language.Tag) lpax.TextMap {
		return languagePack
	}, lpax.Package, language.English)
}
