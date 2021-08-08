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

// Package langpack for yaff cli flags.
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

	// FlagsReportingFormat cli arg for format.
	FlagsReportingFormat
	// FlagsReportingStyle cli arg for style (text format).
	FlagsReportingStyle
	// FlagsReportingTemplate cli arg for template contents (template format).
	FlagsReportingTemplate
	// FlagsReportingTemplateFile cli arg for template file (template format).
	FlagsReportingTemplateFile
	// FlagsReportingIndent cli arg for indent file (json format).
	FlagsReportingIndent
	// FlagsReportingInclude cli arg for columns to include.
	FlagsReportingInclude
	// FlagsReportingExclude cli arg for columns to exclude.
	FlagsReportingExclude
	// FlagsColumnSeparator column separator.
	FlagsColumnSeparator
	// FlagsTtyWidth terminal width.
	FlagsTtyWidth

	// ErrorIndentLessThanZero bad indent.
	ErrorIndentLessThanZero

	// ErrorTemplateAndTemplateFileSet bad template args.
	ErrorTemplateAndTemplateFileSet
)

var languagePack = lpax.TextMap{

	FlagsReportingFormat:            "output format (csv|json|yaml|text). Default is text.",
	FlagsReportingStyle:             "output style (plain|grid|aligned) Default for text is aligned",
	FlagsReportingTemplate:          "template string to use for text output. Uses Go templating syntax",
	FlagsReportingTemplateFile:      "template file path. File must ne in GO templating syntax ",
	FlagsReportingIndent:            "indenting to use with JSON formating, 0 for single line output",
	FlagsReportingInclude:           "columns to output",
	FlagsReportingExclude:           "columns to exclude from output",
	FlagsTtyWidth:                   "override to width of terminal for tty output",
	FlagsColumnSeparator:            "field separator for csv files",
	ErrorIndentLessThanZero:         "Indent (%d) is less than zero",
	ErrorTemplateAndTemplateFileSet: "%s and %s are both set, only one can be applied at once",
}

func init() {
	// Register formats and english language pack.
	lpax.Default().Register(LangPack, func(packID lpax.PackID, langTag language.Tag) lpax.TextMap {
		return languagePack
	}, lpax.Package, language.English)
}
