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

// Package yaff yet another flexible formatter.
package yaff

import (
	"io"
)

// Format identifier.
type Format string

// Formatter supports writing formated data to an writer.
type Formatter interface {

	// Format formats the supplied input data using the format options with output to the writer.
	Format(writer io.Writer, options FormatOptions, data ...interface{}) error
}

// FormatOptions options controlling formatting, every formatter can implement its own options.
type FormatOptions interface{}

// NewFormatter function type to create a new formatter.
// Each formatter type registered will provide an implementation to creeate an instance of an associated formatter.
type NewFormatter func() (Formatter, error)
