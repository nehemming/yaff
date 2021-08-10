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
	"bytes"
	"testing"

	"github.com/nehemming/testsupport"
	"github.com/nehemming/yaff"
)

func TestText(t *testing.T) {
	if Text != yaff.Format("text") {
		t.Errorf("Bad Format name %v", Text)
	}
}

type innerData struct {
	Sin  string    `tabular:"Sun"`
	next *testData //nolint:structcheck,unused
	list [10]int   //nolint:structcheck,unused
}

type testData struct {
	S string
	I int
	F float64
	N innerData
	h int //nolint:structcheck,unused
}

type testData2 struct {
	S  string
	B1 bool
	B2 bool `tabular:",trueonly"`
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

	if err != nil {
		t.Errorf("Formatter Error %v", err)
	}

	expected :=
		`Name Output
   S Hello 
   I 10    
   F 3.14  
 Sun Inside
`
	got := buf.String()

	testsupport.CompareStrings(t, expected, got)
}

func TestNewFormatterGridOption(t *testing.T) {
	fmt, err := NewFormatter()
	if err != nil {
		t.Errorf("Error %v", err)
	}

	if fmt == nil {
		t.Error("No formatter")
	}

	// Generate some output
	var buf bytes.Buffer

	options := NewOptions()

	options.Style = Grid
	options.ExcludeSet["I"] = true
	options.ExcludeHeader = true

	err = fmt.Format(&buf, options, &testData{
		S: "Hello", I: 10, F: 3.14, N: innerData{Sin: "Inside"},
	})

	if err != nil {
		t.Errorf("Formatter Error %v", err)
	}

	expected := `+------+--------+
|    S | Hello  |
+------+--------+
|    F | 3.14   |
+------+--------+
|  Sun | Inside |
+------+--------+
`
	got := buf.String()

	testsupport.CompareStrings(t, expected, got)
}

func TestNewFormatterGridArrayOption(t *testing.T) {
	fmt, err := NewFormatter()
	if err != nil {
		t.Errorf("Error %v", err)
	}

	if fmt == nil {
		t.Error("No formatter")
	}

	// Generate some output
	var buf bytes.Buffer

	options := NewOptions()

	options.Style = Grid
	options.ExcludeSet["I"] = true

	err = fmt.Format(&buf, options, []testData{
		{S: "Hello", I: 10, F: 3.14, N: innerData{Sin: "Inside"}},
		{S: "Bye", I: 32, F: 2.77, N: innerData{Sin: "Outside"}},
	})

	if err != nil {
		t.Errorf("Formatter Error %v", err)
	}

	expected := `+-------+------+---------+
| S     |    F | Sun     |
+-------+------+---------+
| Hello | 3.14 | Inside  |
+-------+------+---------+
| Bye   | 2.77 | Outside |
+-------+------+---------+
`
	got := buf.String()

	testsupport.CompareStrings(t, expected, got)
}

func TestNewFormatterPlainArrayOption(t *testing.T) {
	fmt, err := NewFormatter()
	if err != nil {
		t.Errorf("Error %v", err)
	}

	if fmt == nil {
		t.Error("No formatter")
	}

	// Generate some output
	var buf bytes.Buffer

	options := NewOptions()

	options.Style = Plain
	options.ExcludeSet["I"] = true

	err = fmt.Format(&buf, options, []testData{
		{S: "Hello", I: 10, F: 3.14, N: innerData{Sin: "Inside"}},
		{S: "Bye", I: 32, F: 2.77, N: innerData{Sin: "Outside"}},
	})

	if err != nil {
		t.Errorf("Formatter Error %v", err)
	}

	expected := `S  F  Sun
Hello  3.14  Inside
Bye  2.77  Outside
`
	got := buf.String()

	testsupport.CompareStrings(t, expected, got)
}

func TestNewFormatterSimpleText(t *testing.T) {
	fmt, err := NewFormatter()
	if err != nil {
		t.Errorf("Error %v", err)
	}

	if fmt == nil {
		t.Error("No formatter")
	}

	// Generate some output
	var buf bytes.Buffer

	options := NewOptions()

	options.Style = Plain
	options.ExcludeSet["I"] = true

	err = fmt.Format(&buf, options, "hello there", "today")

	if err != nil {
		t.Errorf("Formatter Error %v", err)
	}

	expected := `Output
hello there
Output
today
`
	got := buf.String()

	testsupport.CompareStrings(t, expected, got)
}

func TestNewFormatterMarkdown(t *testing.T) {
	fmt, err := NewFormatter()
	if err != nil {
		t.Errorf("Error %v", err)
	}

	if fmt == nil {
		t.Error("No formatter")
	}

	// Generate some output
	var buf bytes.Buffer

	options := NewOptions()

	options.Style = Markdown
	options.ExcludeSet["I"] = true

	err = fmt.Format(&buf, options, []testData{
		{S: "Hello", I: 10, F: 3.14, N: innerData{Sin: "Inside"}},
		{S: "Bye", I: 32, F: 2.77, N: innerData{Sin: "Outside"}},
	})

	if err != nil {
		t.Errorf("Formatter Error %v", err)
	}

	expected := `|S|F|Sun|
|-|-|-|
|Hello|3.14|Inside|
|Bye|2.77|Outside|
`
	got := buf.String()

	testsupport.CompareStrings(t, expected, got)
}

func TestNewFormatterBadStyle(t *testing.T) {
	fmt, err := NewFormatter()
	if err != nil {
		t.Errorf("Error %v", err)
	}

	if fmt == nil {
		t.Error("No formatter")
	}

	// Generate some output
	var buf bytes.Buffer

	options := NewOptions()

	options.Style = TableStyle(890)
	options.ExcludeSet["I"] = true

	err = fmt.Format(&buf, options, "hello there", "today")

	if err == nil {
		t.Errorf("No style error")
	}
}

func TestNewFormatterGridBooleanOption(t *testing.T) {
	fmt, err := NewFormatter()
	if err != nil {
		t.Errorf("Error %v", err)
	}

	if fmt == nil {
		t.Error("No formatter")
	}

	// Generate some output
	var buf bytes.Buffer

	options := NewOptions()

	options.Style = Grid

	err = fmt.Format(&buf, options, []testData2{
		{S: "Hello", B1: true, B2: true},
		{S: "Hello", B1: false, B2: false},
	})

	if err != nil {
		t.Errorf("Formatter Error %v", err)
	}

	expected := `+-------+-------+------+
| S     |    B1 |   B2 |
+-------+-------+------+
| Hello |  true | true |
+-------+-------+------+
| Hello | false |      |
+-------+-------+------+
`
	got := buf.String()

	testsupport.CompareStrings(t, expected, got)
}

func TestNewFormatterGridArray(t *testing.T) {
	fmt, err := NewFormatter()
	if err != nil {
		t.Errorf("Error %v", err)
	}

	if fmt == nil {
		t.Error("No formatter")
	}

	// Generate some output
	var buf bytes.Buffer

	options := NewOptions()

	options.Style = Grid

	err = fmt.Format(&buf, options, []string{
		"one", "two",
	})

	if err != nil {
		t.Errorf("Formatter Error %v", err)
	}

	expected := `+--------+
| Output |
+--------+
| one    |
+--------+
| two    |
+--------+
`
	got := buf.String()

	testsupport.CompareStrings(t, expected, got)
}

func TestNewFormatterGridOptionWithWrapping(t *testing.T) {
	fmt, err := NewFormatter()
	if err != nil {
		t.Errorf("Error %v", err)
	}

	if fmt == nil {
		t.Error("No formatter")
	}

	// Generate some output
	var buf bytes.Buffer

	options := NewOptions()

	options.Style = Grid
	options.ExcludeSet["I"] = true
	options.ExcludeHeader = false
	options.TerminalWidth = 30

	err = fmt.Format(&buf, options, &testData{
		S: "Hello this is some long text that needs to wrap", I: 10, F: 3.14, N: innerData{Sin: "Inside"},
	})

	if err != nil {
		t.Errorf("Formatter Error %v", err)
	}

	expected := `+------+---------------------+
| Name | Output              |
+------+---------------------+
|    S | Hello this is some  |
|      | long text that      |
|      | needs to wrap       |
+------+---------------------+
|    F | 3.14                |
+------+---------------------+
|  Sun | Inside              |
+------+---------------------+
`
	got := buf.String()

	testsupport.CompareStrings(t, expected, got)
}
