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
	"testing"
)

func TestGetTextStyleFromString(t *testing.T) {
	ts, err := GetTextStyleFromString("")
	if err != nil {
		t.Errorf("Error %v (%v)", err, "")
	}

	if ts != Aligned {
		t.Errorf("Value %v (%v)", ts, Aligned)
	}

	ts, err = GetTextStyleFromString("aligned")

	if err != nil {
		t.Errorf("Error %v (%v)", err, "aligned")
	}

	if ts != Aligned {
		t.Errorf("Value %v (%v)", ts, "aligned")
	}

	ts, err = GetTextStyleFromString("grid")

	if err != nil {
		t.Errorf("Error %v (%v)", err, "grid")
	}

	if ts != Grid {
		t.Errorf("Value %v (%v)", ts, "grid")
	}

	ts, err = GetTextStyleFromString("plain")

	if err != nil {
		t.Errorf("Error %v (%v)", err, "plain")
	}

	if ts != Plain {
		t.Errorf("Value %v (%v)", ts, "plain")
	}

	_, err = GetTextStyleFromString("other")

	if err == nil {
		t.Error("No error")
	}
}
