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
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

// wrapParamName is the name of the width param used to support wrapping text.
const widthParamName = "width"

func reflectInterface(table *tabular, value reflect.Value) error {
	// If interface is to a pointer de reference it
	if value.Kind() == reflect.Ptr {
		value = reflect.Indirect(value)
	}

	// Determin output type (arrays are tabular) a struct will display as detail
	switch value.Kind() {
	case reflect.Array, reflect.Slice:
		return reflectArray(table, value)
	case reflect.Struct:
		return reflectStructDetail(table, value)
	case reflect.Func:
		return nil
	default:
		col, err := table.addColumn("Output", false, "")
		if err != nil {
			return err
		}
		return table.setField(table.newRow(), col, "%v", value.Interface())
	}
}

func reflectStructDetail(table *tabular, value reflect.Value) error {
	// Validate this is a struct
	if value.Kind() != reflect.Struct {
		panic("unexpected type")
	}

	// Find num of fields in struct to iterate over
	n := value.NumField()
	if n == 0 {
		return nil
	}

	// Add 2 detail columns for field name and value
	_, _ = table.addColumn("Name", true, "")
	_, _ = table.addColumn("Output", false, "")

	// Iterate over the structure
	for i := 0; i < n; i++ {
		if err := reflectFieldNameValue(table, value.Type().Field(i), value.Field(i)); err != nil {
			return err
		}
	}

	return nil
}

func reflectFieldNameValue(table *tabular, t reflect.StructField, value reflect.Value) error {
	// Get tags
	name := getFieldName(t.Name, getTags(t.Tag, tabularTagName))
	if !table.shouldOutputColumn(name) {
		return nil
	}

	// Kind is the declared kind so values of interfaces and pointers are excluded
	kind := value.Kind()

	switch kind {
	case reflect.Array, reflect.Slice, reflect.Func, reflect.Map, reflect.Interface, reflect.Ptr:
		return nil

	case reflect.Struct:
		n := value.NumField()
		if n == 0 {
			return nil
		}

		for i := 0; i < n; i++ {
			if err := reflectFieldNameValue(table, value.Type().Field(i), value.Field(i)); err != nil {
				return err
			}
		}

	default:
		row := table.newRow()

		if err := table.setField(row, colID(0), "%s", name); err != nil {
			return err
		}
		if err := table.setField(row, colID(1), "%v", value.Interface()); err != nil {
			return err
		}
	}
	return nil
}

func reflectArray(table *tabular, value reflect.Value) error {
	// Check array or slice
	k := value.Kind()
	if k != reflect.Array && k != reflect.Slice {
		return fmt.Errorf("Item is not an array or slice - type is %v", value.Type().Name())
	}

	// How big is the array
	n := value.Len()

	if n == 0 {
		return nil
	}

	var col colID
	var err error

	for i := 0; i < n; i++ {
		// Get the array item
		item := value.Index(i)

		// Dereference pointers
		if item.Kind() == reflect.Ptr {
			item = reflect.Indirect(item)
		}

		// Support struct ort or simple value types
		switch item.Kind() {
		case reflect.Struct:
			// row of data
			if i == 0 {
				if err = reflectStructHeader(table, item); err != nil {
					return err
				}
			}

			row := table.newRow()

			if err = reflectStructRow(table, row, item); err != nil {
				return err
			}

		case reflect.Array, reflect.Slice, reflect.Map, reflect.Func, reflect.Ptr, reflect.Interface:
			continue

		default:
			// Output column data
			if i == 0 {
				col, err = table.addColumn("Output", false, "")
				if err != nil {
					return err
				}
			}

			row := table.newRow()

			if err = table.setField(row, col, "%v", item.Interface()); err != nil {
				return err
			}
		}
	}

	return nil
}

func reflectStructHeader(table *tabular, value reflect.Value) error {
	// Check is a struct
	if value.Kind() != reflect.Struct {
		panic("unexpected type")
	}

	// Iterate ver a struct to get the files in the type
	n := value.NumField()
	for i := 0; i < n; i++ {
		if err := reflectFieldHead(table, value.Type().Field(i)); err != nil {
			return err
		}
	}

	return nil
}

type tagData struct {
	Name    string
	Options map[string]bool
	Params  map[string]string
}

func (tagData *tagData) getWidthParam() string {
	if tagData == nil {
		return ""
	}
	return tagData.Params[widthParamName]
}

func getTags(tag reflect.StructTag, tagType string) *tagData {
	// Get the passed in tag and return it with its options or nil if not found
	return parseTags(tag.Get(tagType))
}

func parseTags(tagEntry string) *tagData {
	if tagEntry == "" {
		return nil
	}

	d := tagData{Options: make(map[string]bool), Params: make(map[string]string)}

	for i, s := range strings.Split(tagEntry, ",") {
		s := strings.Trim(s, " ")

		if s == "" {
			continue
		}

		if i == 0 {
			d.Name = s
		} else {
			parts := strings.SplitN(s, "=", 2)
			if len(parts) == 1 {
				d.Options[s] = true
			} else {
				d.Params[parts[0]] = parts[1]
			}
		}
	}

	return &d
}

const (
	tabularTagName = "tabular"
)

func getFieldName(name string, tagInfo *tagData) string {
	// cannot reflect non exported fields
	if name == "" || !unicode.IsUpper([]rune(name)[0]) {
		return ""
	}

	if tagInfo != nil && tagInfo.Name != "" {
		return tagInfo.Name
	}

	return name
}

func reflectFieldHead(table *tabular, t reflect.StructField) error {
	tagData := getTags(t.Tag, tabularTagName)
	name := getFieldName(t.Name, tagData)
	if !table.shouldOutputColumn(name) {
		return nil
	}

	switch t.Type.Kind() {
	case reflect.Array, reflect.Slice, reflect.Func, reflect.Map, reflect.Ptr, reflect.Interface:
		return nil

	case reflect.Struct:

		// Embedded structure, flatten struct
		n := t.Type.NumField()
		for i := 0; i < n; i++ {
			if err := reflectFieldHead(table, t.Type.Field(i)); err != nil {
				return err
			}
		}

	case reflect.String:
		_, err := table.addColumn(name, false, tagData.getWidthParam())
		if err != nil {
			return err
		}

	default:
		_, err := table.addColumn(name, true, tagData.getWidthParam())
		if err != nil {
			return err
		}
	}

	return nil
}

func reflectStructRow(table *tabular, row rowID, value reflect.Value) error {
	// Check we have a struct
	if value.Kind() != reflect.Struct {
		panic("unexpected type")
	}

	n := value.NumField()
	var err error
	col := colID(0)
	for i := 0; i < n; i++ {
		// t is the structure declared type not the type of the value itself
		t := value.Type().Field(i)
		col, err = reflectFieldValue(table, row, col, value.Field(i), t)

		if err != nil {
			return err
		}
	}

	return nil
}

func reflectFieldValue(table *tabular, row rowID, col colID, value reflect.Value, t reflect.StructField) (colID, error) {
	tagInfo := getTags(t.Tag, tabularTagName)
	name := getFieldName(t.Name, tagInfo)
	if !table.shouldOutputColumn(name) {
		return col, nil
	}

	switch t.Type.Kind() {
	case reflect.Array, reflect.Slice, reflect.Func, reflect.Map, reflect.Interface, reflect.Ptr:
		return col, nil

	case reflect.Struct:

		// Nested struct, follow
		n := value.NumField()
		var err error

		for i := 0; i < n; i++ {
			t := value.Type().Field(i)
			col, err = reflectFieldValue(table, row, col, value.Field(i), t)
			if err != nil {
				return col, err
			}
		}

	case reflect.Bool:
		// Allow alternative handling of bool false
		if tagInfo != nil && tagInfo.Options["trueonly"] {
			if !value.Interface().(bool) {
				return col + 1, table.setField(row, col, "")
			}
		}
		fallthrough
	default:
		return col + 1, table.setField(row, col, "%v", value.Interface())
	}

	return col, nil
}
