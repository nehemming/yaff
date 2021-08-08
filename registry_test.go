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

package yaff

import (
	"errors"
	"io"
	"testing"
)

// exampleCount is the number of example formatters registered in the shared formatter.
const exampleCount = 1

func TestNewRegistry(t *testing.T) {
	reg := NewRegistry()

	formats := reg.Formats()

	if len(formats) > 0 {
		t.Errorf("Formats too long %v", len(formats))
	}

	f, err := reg.GetFormatter(Format("json"))

	if f != nil {
		t.Error("Magic format exists")
	}
	if err == nil {
		t.Error("No error exists")
	}
}

func TestRegistrationWithNoFormatterCreated(t *testing.T) {
	reg := NewRegistry()

	testFormat := Format("test")

	reg.Register(testFormat, func() (Formatter, error) {
		return nil, errors.New("testerr")
	})

	formats := reg.Formats()

	if len(formats) != 1 {
		t.Errorf("Formats len not 1  %v", len(formats))
	}

	f, err := reg.GetFormatter(testFormat)

	if err == nil || err.Error() != "testerr" {
		t.Errorf("Error get wrong %v", err)
	}

	if f != nil {
		t.Error("Magic format exists")
	}
}

type testFormatter struct{}

func (f *testFormatter) Format(writer io.Writer, options FormatOptions, data ...interface{}) error {
	return nil
}

func TestRegistrationWithFormatterCreation(t *testing.T) {
	reg := NewRegistry()

	testFormat := Format("test")

	reg.Register(testFormat, func() (Formatter, error) {
		return &testFormatter{}, nil
	})

	formats := reg.Formats()

	if len(formats) != 1 {
		t.Errorf("Formats len not 1  %v", len(formats))
	}

	f, err := reg.GetFormatter(testFormat)
	if err != nil {
		t.Errorf("Error get %v", err)
	}

	if f == nil {
		t.Error("No test format exists")
	}
}

func TestRegistrationWithNilPassed(t *testing.T) {
	reg := NewRegistry()

	testFormat := Format("test")

	reg.Register(testFormat, nil)

	formats := reg.Formats()

	if len(formats) != 1 {
		t.Errorf("Formats len not 1  %v", len(formats))
	}

	f, err := reg.GetFormatter(testFormat)

	if err == nil {
		t.Error("Error get no error")
	}

	if f != nil {
		t.Error("Magic format exists")
	}
}

func TestSharedRegistry(t *testing.T) {
	reg := Formatters()

	formats := reg.Formats()

	if len(formats) > exampleCount {
		t.Errorf("Formats too long %v", len(formats))
	}
}
