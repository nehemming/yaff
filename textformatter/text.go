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

// Package textformatter is the yaff Text formatter.
package textformatter

import (
	"io"
	"reflect"
	"strings"

	"github.com/nehemming/lpax"
	"github.com/nehemming/yaff"
	"github.com/nehemming/yaff/langpack"
)

// Text format.
const Text = yaff.Format("text")

// NewFormatter return a new formatter.
func NewFormatter() (yaff.Formatter, error) {
	return &formatter{}, nil
}

type formatter struct{}

// Options for the JSON formatter.
type Options struct {
	Style           TableStyle
	ExcludeHeader   bool
	ColumnSeparator string
	ColumnSet       map[string]bool
	ExcludeSet      map[string]bool
	// TerminalWidth is the width oth the terminal the aligned text is being output to
	// if this is 0 no wrapping will be used.  For values > column min width this value will
	// be used to wrap text.
	TerminalWidth int
}

// NewOptions return new options.
func NewOptions() Options {
	return Options{
		Style:           Aligned,
		ColumnSeparator: "\t",
		ColumnSet:       make(map[string]bool),
		ExcludeSet:      make(map[string]bool),
	}
}

func (f *formatter) Format(writer io.Writer, options yaff.FormatOptions, data ...interface{}) error {
	if options == nil {
		options = NewOptions()
	}

	// convert options type
	textOptions, ok := options.(Options)
	if !ok {
		return lpax.Errorf(langpack.ErrorInvalidOptionType, options, Text)
	}

	textOptions = normalizeOptions(textOptions)

	for _, d := range data {
		if err := renderStyledText(writer, d, textOptions); err != nil {
			return err
		}
	}

	return nil
}

func normalizeOptions(options Options) Options {
	// Use lower case for all exclusion and colset settings
	xSet := make(map[string]bool)
	cSet := make(map[string]bool)

	if options.ExcludeSet != nil {
		for k, v := range options.ExcludeSet {
			xSet[strings.ToLower(k)] = v
		}
	}

	if options.ColumnSet != nil {
		for k, v := range options.ColumnSet {
			cSet[strings.ToLower(k)] = v
		}
	}

	options.ExcludeSet = xSet
	options.ColumnSet = cSet
	return options
}

func renderStyledText(writer io.Writer, d interface{}, options Options) error {
	table := newTabular(options.ColumnSet, options.ExcludeSet)

	value := reflect.ValueOf(d)

	if err := reflectInterface(table, value); err != nil {
		return err
	}

	return table.write(writer, options.Style, options.ExcludeHeader, options.ColumnSeparator, options.TerminalWidth)
}

func init() {
	// Register this formatter.
	yaff.Formatters().Register(Text, NewFormatter)
}
