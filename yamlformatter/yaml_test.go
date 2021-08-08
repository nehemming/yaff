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

package yamlformatter

import (
	"bytes"
	"testing"

	"github.com/nehemming/testsupport"
	"github.com/nehemming/yaff"
)

func TestYAML(t *testing.T) {
	if YAML != yaff.Format("yaml") {
		t.Errorf("Bad Format name %v", YAML)
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

func TestNewFormatter(t *testing.T) {
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

	expected := `s: Hello
i: 10
f: 3.14
"n":
    sin: Inside

`

	got := buf.String()

	testsupport.CompareStrings(t, expected, got)
}

func TestNewFormatterTwoOutputs(t *testing.T) {
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
	},
		&testData{S: "Again", I: 99, F: 2.71})

	if err != nil {
		t.Errorf("Formatter Error %v", err)
	}

	expected := `s: Hello
i: 10
f: 3.14
"n":
    sin: Inside

---
s: Again
i: 99
f: 2.71
"n":
    sin: ""

`

	got := buf.String()

	testsupport.CompareStrings(t, expected, got)
}
