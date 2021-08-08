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

// Package cliflags allows an application to bind parameters to control outputt formats
package cliflags

import (
	"strings"

	"github.com/nehemming/lpax"
	"github.com/nehemming/yaff"
	lp "github.com/nehemming/yaff/cliflags/langpack"
	"github.com/nehemming/yaff/csvformatter"
	"github.com/nehemming/yaff/jsonformatter"
	"github.com/nehemming/yaff/templateformatter"
	"github.com/nehemming/yaff/textformatter"
	"github.com/nehemming/yaff/yamlformatter"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// FlagsReportingFormat output format.
	FlagsReportingFormat = "format"
	// FlagsReportingTemplate inline reporting template.
	FlagsReportingTemplate = "template"
	// FlagsReportingTemplateFile template file.
	FlagsReportingTemplateFile = "templatefile"
	// FlagsReportingIndent indent level.
	FlagsReportingIndent = "jsonindent"
	// FlagsReportingStyle text style.
	FlagsReportingStyle = "style"
	// FlagsReportingInclude column set to include in results.
	FlagsReportingInclude = "colset"
	// FlagsReportingExclude Columns to exclude.
	FlagsReportingExclude = "excludecols"
	// FlagsColumnSeparator Columns to exclude.
	FlagsColumnSeparator = "separator"
	// FlagsTtyWidth width override for tty.
	FlagsTtyWidth = "ttywidth"
)

const (
	// ParamsReportingFormat reporting format.
	ParamsReportingFormat = "format"
	// ParamsReportingStyle reporting style.
	ParamsReportingStyle = "style"
	// ParamsReportingIndent indent.
	ParamsReportingIndent = "indent"
	// ParamColumnSeparator is the column separator.
	ParamColumnSeparator = "separator"
	// ParamTtyWidth is the tty width to use.
	ParamTtyWidth = "ttywidth"
)

// AddFormattingFlags adds formatting flags to a flag set.
func AddFormattingFlags(flags *pflag.FlagSet) {
	tf := lpax.Default()

	flags.String(FlagsReportingFormat, "", tf.Text(lp.FlagsReportingFormat))
	flags.String(FlagsReportingStyle, "", tf.Text(lp.FlagsReportingStyle))
	flags.String(FlagsReportingTemplate, "", tf.Text(lp.FlagsReportingTemplate))
	flags.String(FlagsReportingTemplateFile, "", tf.Text(lp.FlagsReportingTemplateFile))
	flags.Int(FlagsReportingIndent, 2, tf.Text(lp.FlagsReportingIndent))
	flags.String(FlagsReportingInclude, "", tf.Text(lp.FlagsReportingInclude))
	flags.String(FlagsReportingExclude, "", tf.Text(lp.FlagsReportingExclude))
	flags.String(FlagsColumnSeparator, ",", tf.Text(lp.FlagsColumnSeparator))
	flags.Int(FlagsTtyWidth, 0, tf.Text(lp.FlagsTtyWidth))
}

// BindFormattingParamsToFlags binds a flags set to viper params.
// configBase is the base configuration path to use for the passed config viper.
func BindFormattingParamsToFlags(flags *pflag.FlagSet, config *viper.Viper, configBase string) error {
	if len(configBase) > 0 && !strings.HasSuffix(configBase, ".") {
		configBase = configBase + "."
	}

	if err := config.BindPFlag(configBase+ParamsReportingFormat, flags.Lookup(FlagsReportingFormat)); err != nil {
		return err
	}
	if err := config.BindPFlag(configBase+ParamsReportingStyle, flags.Lookup(FlagsReportingStyle)); err != nil {
		return err
	}
	if err := config.BindPFlag(configBase+ParamsReportingIndent, flags.Lookup(FlagsReportingIndent)); err != nil {
		return err
	}
	if err := config.BindPFlag(configBase+ParamColumnSeparator, flags.Lookup(FlagsColumnSeparator)); err != nil {
		return err
	}
	if err := config.BindPFlag(configBase+ParamTtyWidth, flags.Lookup(FlagsTtyWidth)); err != nil {
		return err
	}

	config.SetDefault(configBase+ParamColumnSeparator, ",")

	return nil
}

// GetFormmatterFromFlags returns a formatter and its options based off command line a config settings.
func GetFormmatterFromFlags(flags *pflag.FlagSet, v *viper.Viper, defaultFormat yaff.Format, configBase string) (yaff.Formatter, yaff.FormatOptions, error) {
	if len(configBase) > 0 && !strings.HasSuffix(configBase, ".") {
		configBase = configBase + "."
	}

	// Validate format
	format := yaff.Format(v.GetString(configBase + ParamsReportingFormat))
	if format == "" {
		format = defaultFormat
	}

	formatter, err := yaff.Formatters().GetFormatter(format)
	if err != nil {
		return nil, nil, err
	}

	var formatOptions yaff.FormatOptions

	// Bind format options from args
	switch format {
	case templateformatter.Template:
		option := templateformatter.NewOptions()

		option.Template, _ = flags.GetString(FlagsReportingTemplate)
		option.TemplateFile, _ = flags.GetString(FlagsReportingTemplateFile)

		if option.Template != "" && option.TemplateFile != "" {
			return nil, nil, lpax.Errorf(lp.ErrorTemplateAndTemplateFileSet,
				FlagsReportingTemplate, FlagsReportingTemplateFile)
		}

		formatOptions = option

	case textformatter.Text:
		option := textformatter.NewOptions()

		option.Style, err = textformatter.GetTextStyleFromString(v.GetString(configBase + ParamsReportingStyle))
		if err != nil {
			return nil, nil, err
		}

		list, _ := flags.GetString(FlagsReportingInclude)
		option.ColumnSet = mapFromList(list)
		list, _ = flags.GetString(FlagsReportingExclude)
		option.ExcludeSet = mapFromList(list)
		option.ColumnSeparator = v.GetString(configBase + ParamColumnSeparator)
		option.TerminalWidth = v.GetInt(configBase + ParamTtyWidth)
		formatOptions = option

	case jsonformatter.JSON:
		option := jsonformatter.NewOptions()

		indent := v.GetInt(configBase + ParamsReportingIndent)
		if indent < 0 {
			return nil, nil, lpax.Errorf(lp.ErrorIndentLessThanZero, indent)
		}

		option.Indent = indent
		formatOptions = option

	case csvformatter.CSV:
		option := csvformatter.NewOptions()
		option.ColumnSeparator = v.GetString(configBase + ParamColumnSeparator)
		formatOptions = option

	case yamlformatter.YAML:
		fallthrough
	default:
	}

	return formatter, formatOptions, nil
}

func mapFromList(list string) map[string]bool {
	m := make(map[string]bool)
	items := strings.Split(list, ",")

	for _, s := range items {
		s = strings.ToLower(strings.Trim(s, " "))
		if s != "" {
			m[s] = true
		}
	}

	return m
}
