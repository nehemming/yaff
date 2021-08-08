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

import "testing"

func TestParseTagsEmpty(t *testing.T) {
	d := parseTags("")
	if d != nil {
		t.Error("unexpected", *d)
	}
}

func TestParseTagsSimple(t *testing.T) {
	d := parseTags("winter")
	if d == nil {
		t.Error("unexpected nil d")
		return
	}

	if d.Name != "winter" {
		t.Error("unexpected name", d.Name)
	}
}

func TestParseTagsOptions(t *testing.T) {
	d := parseTags("winter, summer, spring")
	if d == nil {
		t.Error("unexpected nil d")
		return
	}

	if d.Name != "winter" {
		t.Error("unexpected name", d.Name)
	}

	if !d.Options["summer"] && !d.Options["spring"] {
		t.Error("unexpected name", d.Options)
	}
}

func TestParseTagsOptionsAndParams(t *testing.T) {
	d := parseTags("winter, summer, spring=12.7")
	if d == nil {
		t.Error("unexpected nil d")
		return
	}

	if d.Name != "winter" {
		t.Error("unexpected name", d.Name)
	}

	if !d.Options["summer"] || d.Options["spring"] {
		t.Error("unexpected name", d.Options)
	}

	if d.Params["spring"] != "12.7" {
		t.Error("unexpected spring", d.Params["spring"])
	}
}
