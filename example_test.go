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

package yaff_test

import (
	"bytes"
	"fmt"

	"github.com/nehemming/yaff"
	"github.com/nehemming/yaff/textformatter"
)

type dataRow struct {
	StringVal   string
	IntVal      int
	FloatVal    float64
	BoolVal     bool
	OnlyTrueVal bool `tabular:",trueonly"`
	hidden      int  //nolint:structcheck
	Embedded    embeddedData
}

type embeddedData struct {
	innerStr string
	IntTwo   int
}

func main() {
	// Get the text formatter
	formatter, err := yaff.Formatters().GetFormatter(textformatter.Text)
	if err != nil {
		fmt.Println(err)
		return
	}

	options := textformatter.NewOptions()

	// Output in grid style
	options.Style = textformatter.Grid

	// Generate some output, using io writer interface of Buffer
	var buf bytes.Buffer

	err = formatter.Format(&buf, options, []dataRow{
		{
			StringVal: "Hello", IntVal: 10, FloatVal: 9.8, BoolVal: true, OnlyTrueVal: true, hidden: 99,
			Embedded: embeddedData{innerStr: "hidden", IntTwo: 2},
		},
		{
			StringVal: "There", IntVal: 20, FloatVal: 3.14, BoolVal: false, OnlyTrueVal: false, hidden: 99,
			Embedded: embeddedData{innerStr: "hidden2", IntTwo: 2},
		},
	})

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(buf.String())
	}
}

func Example() {
	main()

	// Output:
	// +-----------+--------+----------+---------+-------------+--------+
	// | StringVal | IntVal | FloatVal | BoolVal | OnlyTrueVal | IntTwo |
	// +-----------+--------+----------+---------+-------------+--------+
	// | Hello     |     10 |      9.8 |    true |        true |      2 |
	// +-----------+--------+----------+---------+-------------+--------+
	// | There     |     20 |     3.14 |   false |             |      2 |
	// +-----------+--------+----------+---------+-------------+--------+
}
