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

package templateformatter

import (
	"bytes"
	"testing"

	"github.com/nehemming/testsupport"
	"github.com/nehemming/yaff"
)

func TestTemplate(t *testing.T) {
	if Template != yaff.Format("template") {
		t.Errorf("Bad Format name %v", Template)
	}
}

type innerData struct {
	Sin string
}

type testData struct {
	S string
	I int
	F float64
	N innerData
}

func TestNewFormatterNoOptions(t *testing.T) {
	fmt, err := NewFormatter()
	if err != nil {
		t.Errorf("Error %v", err)
	}

	if fmt == nil {
		t.Error("No formatter")
	}

	// Generate some output
	var buf bytes.Buffer

	err = fmt.Format(&buf, nil, &testData{
		S: "Hello", I: 10, F: 3.14, N: innerData{Sin: "Inside"},
	})

	if err == nil {
		t.Errorf("No error for missing option")
	}
}

func TestNewFormatter(t *testing.T) {
	fmt, err := NewFormatter()
	if err != nil {
		t.Errorf("Error %v", err)
	}

	if fmt == nil {
		t.Error("No formatter")
	}

	options := Options{
		Template: "{{ .S }} {{.N.Sin}}",
	}

	// Generate some output
	var buf bytes.Buffer

	err = fmt.Format(&buf, options, &testData{
		S: "Hello", I: 10, F: 3.14, N: innerData{Sin: "Inside"},
	})

	if err != nil {
		t.Errorf("Formatter Error %v", err)
	}

	expected := "Hello Inside"

	got := buf.String()

	testsupport.CompareStrings(t, expected, got)
}

func TestNewFormatterLoadFromFile(t *testing.T) {
	fmt, err := NewFormatter()
	if err != nil {
		t.Errorf("Error %v", err)
	}

	if fmt == nil {
		t.Error("No formatter")
	}

	options := Options{
		TemplateFile: "./testdata/template.tmpl",
	}

	// Generate some output
	var buf bytes.Buffer

	err = fmt.Format(&buf, options, &testData{
		S: "Hello", I: 10, F: 3.14, N: innerData{Sin: "Inside"},
	})

	if err != nil {
		t.Errorf("Formatter Error %v", err)
	}

	expected := "Hello Inside!"

	got := buf.String()

	testsupport.CompareStrings(t, expected, got)
}
