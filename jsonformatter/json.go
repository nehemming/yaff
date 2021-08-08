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

// Package jsonformatter is the yaff JSON formatter.
package jsonformatter

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/nehemming/lpax"
	"github.com/nehemming/yaff"
	"github.com/nehemming/yaff/langpack"
)

// JSON format.
const JSON = yaff.Format("json")

// NewFormatter return a new formatter.
func NewFormatter() (yaff.Formatter, error) {
	return &formatter{}, nil
}

type formatter struct{}

// Options for the JSON formatter.
type Options struct {
	Indent       int
	IndentString string
}

// NewOptions return new options.
func NewOptions() Options {
	return Options{
		Indent:       2,
		IndentString: " ",
	}
}

func (f *formatter) Format(writer io.Writer, options yaff.FormatOptions, data ...interface{}) error {
	if options == nil {
		options = NewOptions()
	}

	// convert options type
	jsonOptions, ok := options.(Options)
	if !ok {
		return lpax.Errorf(langpack.ErrorInvalidOptionType, options, JSON)
	}

	// use JSON serialization, concat data into a single doc
	var d interface{}

	if len(data) == 1 {
		d = data[0]
	} else {
		d = data
	}

	var buf []byte

	var err error

	ident := jsonOptions.IndentString
	if ident == "" {
		ident = " "
	}

	if jsonOptions.Indent > 0 {
		buf, err = json.MarshalIndent(d, "", strings.Repeat(ident, jsonOptions.Indent))
	} else {
		buf, err = json.Marshal(d)
	}

	if err != nil {
		return err
	}

	_, err = writer.Write(buf)
	//nolint:errcheck
	writer.Write([]byte("\n"))

	return err
}

func init() {
	// Register this formatter
	yaff.Formatters().Register(JSON, NewFormatter)
}
