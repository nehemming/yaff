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

package csvformatter

import (
	"bytes"
	"testing"

	"github.com/nehemming/testsupport"
	"github.com/nehemming/yaff"
)

func TestCSV(t *testing.T) {
	if CSV != yaff.Format("csv") {
		t.Errorf("Bad Format name %v", CSV)
	}
}

type testData struct {
	S string
	I int
	F float64
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

	err = fmt.Format(&buf, nil, []testData{
		{S: "Hello", I: 10, F: 3.14},
	})

	if err != nil {
		t.Errorf("Formatter Error %v", err)
	}

	expected := `S,I,F
Hello,10,3.14
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

	err = fmt.Format(&buf, nil,
		[]testData{
			{S: "Hello", I: 10, F: 3.14},
			{S: "Train", I: 11, F: 3.99},
		},
		[]testData{
			{S: "Hello", I: 99, F: 2.7},
		})

	if err != nil {
		t.Errorf("Formatter Error %v", err)
	}

	expected := `S,I,F
Hello,10,3.14
Train,11,3.99
S,I,F
Hello,99,2.7
`

	got := buf.String()

	testsupport.CompareStrings(t, expected, got)
}
