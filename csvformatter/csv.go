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

// Package csvformatter is the yaff JSON formatter.
package csvformatter

import (
	"encoding/csv"
	"io"

	"github.com/gocarina/gocsv"
	"github.com/nehemming/lpax"
	"github.com/nehemming/yaff"
	"github.com/nehemming/yaff/langpack"
)

// CSV format.
const CSV = yaff.Format("csv")

// NewFormatter return a new formatter.
func NewFormatter() (yaff.Formatter, error) {
	return &formatter{}, nil
}

type formatter struct{}

// Options for the JSON formatter.
type Options struct {
	IncludeHeader   bool
	ColumnSeparator string
}

// NewOptions return new options.
func NewOptions() Options {
	return Options{
		IncludeHeader:   true,
		ColumnSeparator: ",",
	}
}

func (f *formatter) Format(writer io.Writer, options yaff.FormatOptions, data ...interface{}) error {
	if options == nil {
		options = NewOptions()
	}

	// convert options type
	csvOptions, ok := options.(Options)
	if !ok {
		return lpax.Errorf(langpack.ErrorInvalidOptionType, options, CSV)
	}

	// Set up writer with options
	out := gocsv.NewSafeCSVWriter(csv.NewWriter(writer))

	if csvOptions.ColumnSeparator != "" {
		out.Comma = ([]rune(csvOptions.ColumnSeparator))[0]
	} else {
		out.Comma = ','
	}

	// Set header mode
	var marshaller func(interface{}, gocsv.CSVWriter) error
	if csvOptions.IncludeHeader {
		marshaller = gocsv.MarshalCSV
	} else {
		marshaller = gocsv.MarshalCSVWithoutHeaders
	}

	// marahal each output
	for _, d := range data {
		err := marshaller(d, out)
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	// Register this formatter
	yaff.Formatters().Register(CSV, NewFormatter)
}
