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

package cliflags

import (
	"testing"

	"github.com/nehemming/yaff/csvformatter"
	"github.com/nehemming/yaff/jsonformatter"
	"github.com/nehemming/yaff/templateformatter"
	"github.com/nehemming/yaff/textformatter"
	"github.com/nehemming/yaff/yamlformatter"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func TestMapFromListEmpty(t *testing.T) {
	m := mapFromList("")

	if len(m) != 0 {
		t.Error("unexpected len", len(m))
	}
}

func TestMapFromListSingle(t *testing.T) {
	m := mapFromList("one")

	if len(m) != 1 {
		t.Error("unexpected len", len(m))
	}

	if !m["one"] {
		t.Error("missing m one")
	}
}

func TestMapFromListMany(t *testing.T) {
	m := mapFromList("one, two,    three, two")

	if len(m) != 3 {
		t.Error("unexpected len", len(m))
	}

	if !m["one"] {
		t.Error("missing m one")
	}
	if !m["two"] {
		t.Error("missing m two")
	}
	if !m["three"] {
		t.Error("missing m three")
	}
}

func TestGetFormmatterFromFlagsDefault(t *testing.T) {
	flags := new(pflag.FlagSet)
	v := viper.New()
	_, fo, err := GetFormmatterFromFlags(flags, v, jsonformatter.JSON, "cfg")
	if err != nil {
		t.Error("err:", err)
	}

	if fo == nil {
		t.Error("fo:", fo)
	}
}

func TestGetFormmatterFromFlagsJSON(t *testing.T) {
	flags := new(pflag.FlagSet)
	v := viper.New()
	_, fo, err := GetFormmatterFromFlags(flags, v, jsonformatter.JSON, "cfg")
	if err != nil {
		t.Error("err:", err)
	}

	if fo == nil {
		t.Error("fo:", fo)
	}

	jOpt := fo.(jsonformatter.Options)

	if jOpt.Indent != 0 {
		t.Error("indent:", jOpt.Indent)
	}
}

func TestGetFormmatterFromFlagsJSONFmt(t *testing.T) {
	flags := new(pflag.FlagSet)
	v := viper.New()

	v.SetDefault("cfg."+ParamsReportingIndent, 4)
	_, fo, err := GetFormmatterFromFlags(flags, v, jsonformatter.JSON, "cfg")
	if err != nil {
		t.Error("err:", err)
	}

	if fo == nil {
		t.Error("fo:", fo)
	}

	jOpt := fo.(jsonformatter.Options)

	if jOpt.Indent != 4 {
		t.Error("indent:", jOpt.Indent)
	}
}

func TestGetFormmatterFromFlagsYAML(t *testing.T) {
	flags := new(pflag.FlagSet)
	v := viper.New()
	_, fo, err := GetFormmatterFromFlags(flags, v, yamlformatter.YAML, "cfg")
	if err != nil {
		t.Error("err:", err)
	}

	if fo != nil {
		t.Error("fo:", fo)
		return
	}
}

func TestGetFormmatterFromFlagsCSVFmt(t *testing.T) {
	flags := new(pflag.FlagSet)
	v := viper.New()

	v.SetDefault("cfg."+ParamColumnSeparator, "|")
	_, fo, err := GetFormmatterFromFlags(flags, v, csvformatter.CSV, "cfg")
	if err != nil {
		t.Error("err:", err)
	}

	if fo == nil {
		t.Error("fo:", fo)
	}

	csvOpt := fo.(csvformatter.Options)

	if csvOpt.ColumnSeparator != "|" {
		t.Error("ColumnSeparator:", csvOpt.ColumnSeparator)
	}
}

func TestGetFormmatterFromFlagsTextFmt(t *testing.T) {
	flags := new(pflag.FlagSet)
	v := viper.New()

	v.SetDefault("cfg."+ParamColumnSeparator, "|")
	_, fo, err := GetFormmatterFromFlags(flags, v, textformatter.Text, "cfg")
	if err != nil {
		t.Error("err:", err)
	}

	if fo == nil {
		t.Error("fo:", fo)
	}

	textOut := fo.(textformatter.Options)

	if textOut.ColumnSeparator != "|" {
		t.Error("ColumnSeparator:", textOut.ColumnSeparator)
	}
}

func TestGetFormmatterFromFlagsTemplateFmt(t *testing.T) {
	flags := new(pflag.FlagSet)
	v := viper.New()

	AddFormattingFlags(flags)

	_ = flags.Parse([]string{"--template", "{{hello}}"})
	_, fo, err := GetFormmatterFromFlags(flags, v, templateformatter.Template, "cfg")
	if err != nil {
		t.Error("err:", err)
	}

	if fo == nil {
		t.Error("fo:", fo)
	}

	textOut := fo.(templateformatter.Options)

	if textOut.Template != "{{hello}}" {
		t.Error("Template:", textOut.Template)
	}
}

func TestBindFormattingParamsToFlags(t *testing.T) {
	flags := new(pflag.FlagSet)
	v := viper.New()
	AddFormattingFlags(flags)
	err := BindFormattingParamsToFlags(flags, v, "cfg")
	if err != nil {
		t.Error("err:", err)
	}
	_ = flags.Parse([]string{"--format", "JSON"})
	fmt := v.GetString("cfg." + ParamsReportingFormat)

	if fmt != "JSON" {
		t.Error("format:", fmt)
	}
}
