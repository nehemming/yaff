# yaff

Yet another flexible formatter 

Reflect, render and output arbitrary data structures using plugin formatters.  Supported formats include CSV, JSON, YAML, Text and Go templates.

 * [Installation](#install) 
 * [Features](#features)
 * [Getting Started](#start)
 * [Contributing](#contrib)
 * [License](#license)

## Status

Yaff is feature rich and a working MVP   The documentation is still a bit light and increased test coverage is required.

![Status](https://img.shields.io/badge/Status-ALPHA-red?style=for-the-badge)
[![Build Status](https://img.shields.io/circleci/build/gh/nehemming/yaff/master?style=for-the-badge)](https://github.com/nehemming/yaff) 
[![Release](https://img.shields.io/github/v/release/nehemming/yaff.svg?style=for-the-badge)](https://github.com/nehemming/yaff/releases/latest)
[![Coveralls](https://img.shields.io/coveralls/github/nehemming/yaff?style=for-the-badge)](https://coveralls.io/github/nehemming/yaff)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=for-the-badge)](#license)
[![GoReportCard](https://goreportcard.com/badge/github.com/nehemming/yaff?test=0&style=for-the-badge)](https://goreportcard.com/report/github.com/nehemming/yaff)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](http://godoc.org/github.com/nehemming/cirocket)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=for-the-badge)](https://conventionalcommits.org)
[![Uses: cirocket](https://img.shields.io/badge/Uses-cirocket-orange?style=for-the-badge)](https://github.com/nehemming/cirocket)
[![Uses: GoReleaser](https://img.shields.io/badge/uses-goreleaser-green.svg?style=for-the-badge)](https://github.com/goreleaser)

## <a name="install"></a>Installation

```bash
go get -u https://github.com/nehemming/yaff
```

or clone this repo to your local machine using

```bash
git clone https://github.com/nehemming/yaff
```

This project requires Go 1.15 or newer and supports modules.

## <a name="features"></a>Key features

 *  Reflects arbitrary data structures to output formatted text
 *  Plug in formatter model, with built in support for csv, json, yaml, text and go templates
 * Integrates with [Viper](https://github.com/spf13/viper) and [Cobra](https://github.com/spf13/cobra) cli application commands to bind formatting flags and configuration parameters.
 * Text formatter supports auto sizing word wrapping grid.

## <a name="start"></a>Getting started

Here is an example outputting an array of structures in a grid

```go
package main

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
	hidden      int
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
		{StringVal: "Hello", IntVal: 10, FloatVal: 9.8, BoolVal: true, OnlyTrueVal: true, hidden: 99,
			Embedded: embeddedData{innerStr: "hidden", IntTwo: 2}},
		{StringVal: "There", IntVal: 20, FloatVal: 3.14, BoolVal: false, OnlyTrueVal: false, hidden: 99,
			Embedded: embeddedData{innerStr: "hidden2", IntTwo: 2}},
	})

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(buf.String())
	}
}

// Output:
// +-----------+--------+----------+---------+-------------+--------+
// | StringVal | IntVal | FloatVal | BoolVal | OnlyTrueVal | IntTwo |
// +-----------+--------+----------+---------+-------------+--------+
// | Hello     |     10 |      9.8 |    true |        true |      2 |
// +-----------+--------+----------+---------+-------------+--------+
// | There     |     20 |     3.14 |   false |             |      2 |
// +-----------+--------+----------+---------+-------------+--------+
```

## <a name="contrib"></a>Contributing
We would welcome contributions to this project.  Please read our [CONTRIBUTION](https://github.com/nehemming/yaff/blob/master/CONTRIBUTING.md) file for further details on how you can participate or report any issues.

## <a name="license"></a>License

[![FOSSA Status](https://app.fossa.com/api/projects/custom%2B26823%2Fgit%40github.com%3Anehemming%2Fyaff.git.svg?type=small)](https://app.fossa.com/projects/custom%2B26823%2Fgit%40github.com%3Anehemming%2Fyaff.git?ref=badge_small)

This software is licensed under the [Apache License](https://github.com/nehemming/yaff/blob/master/LICENSE). 