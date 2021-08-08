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

// Package yamlformatter is the yaff YAML formatter.
package yamlformatter

import (
	"io"

	"github.com/nehemming/yaff"
	"gopkg.in/yaml.v3"
)

// YAML format.
const YAML = yaff.Format("yaml")

// NewFormatter return a new formatter.
func NewFormatter() (yaff.Formatter, error) {
	return &formatter{}, nil
}

type formatter struct{}

func (f *formatter) Format(writer io.Writer, options yaff.FormatOptions, data ...interface{}) error {
	var err error

	n := len(data)

	if n > 0 { // will need a final new line
		//nolint:errcheck
		defer writer.Write([]byte("\n"))
	}

	for _, d := range data {
		buf, err := yaml.Marshal(d)
		if err != nil {
			return err
		}

		_, err = writer.Write(buf)
		if err != nil {
			return err
		}

		n--

		if n > 0 {
			_, err = writer.Write([]byte("\n---\n"))
			if err != nil {
				return err
			}
		}
	}

	return err
}

func init() {
	// Register this formatter
	yaff.Formatters().Register(YAML, NewFormatter)
}
