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

// Package templateformatter is the yaff Template formatter.
package templateformatter

import (
	"io"

	tt "text/template"

	"github.com/nehemming/fsio"
	"github.com/nehemming/lpax"
	"github.com/nehemming/yaff"
	"github.com/nehemming/yaff/langpack"
)

// Template format.
const Template = yaff.Format("template")

// NewFormatter return a new formatter.
func NewFormatter() (yaff.Formatter, error) {
	return &formatter{}, nil
}

type formatter struct{}

// Options for the JSON formatter.
type Options struct {
	Template     string
	TemplateFile string
}

// NewOptions return new options.
func NewOptions() Options {
	return Options{}
}

func (f *formatter) Format(writer io.Writer, options yaff.FormatOptions, data ...interface{}) error {
	if options == nil {
		options = NewOptions()
	}

	// convert options type.
	templateOptions, ok := options.(Options)
	if !ok {
		return lpax.Errorf(langpack.ErrorInvalidOptionType, options, Template)
	}

	if templateOptions.Template != "" {
		return reportTextTemplate(writer, templateOptions.Template, data)
	} else if templateOptions.TemplateFile != "" {
		buf, err := fsio.ReadFileFromPath(templateOptions.TemplateFile)
		if err != nil {
			return err
		}
		return reportTextTemplate(writer, string(buf), data)
	}

	// No format specified
	return lpax.Errorf(langpack.ErrorNoTemplateDefinition, Template)
}

func reportTextTemplate(writer io.Writer, template string, data []interface{}) error {
	// prep template
	t, err := tt.New("main").Parse(template)
	if err != nil {
		return err
	}

	// run the template per input

	for _, d := range data {
		err = t.Execute(writer, d)
		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
	// Register this formatter
	yaff.Formatters().Register(Template, NewFormatter)
}
