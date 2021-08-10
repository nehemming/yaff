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
	"io"
	"strconv"
	"strings"

	"github.com/nehemming/lpax"
	"github.com/nehemming/yaff/langpack"

	"github.com/eidolon/wordwrap"
)

// rowID row ID.
type rowID int

// colID column id.
type colID int

type column struct {
	name       string
	width      int
	rightAlign bool
	minWidth   int
}

// Tabular data.
type tabular struct {
	columns    []*column
	rows       [][]string
	columnSet  map[string]bool
	excludeSet map[string]bool
}

// newTabular create a new tabular output.
func newTabular(columnSet map[string]bool, excludeSet map[string]bool) *tabular {
	if columnSet == nil {
		columnSet = make(map[string]bool)
	}

	if excludeSet == nil {
		excludeSet = make(map[string]bool)
	}

	return &tabular{
		columns:    make([]*column, 0, 2),
		rows:       make([][]string, 0, 2),
		columnSet:  columnSet,
		excludeSet: excludeSet,
	}
}

// shouldOutputColumn returns false if the column should not be output.
func (tablet *tabular) shouldOutputColumn(columnName string) bool {
	if columnName == "" || columnName == "-" {
		return false
	}

	columnName = strings.ToLower(columnName)

	if len(tablet.columnSet) > 0 && !tablet.columnSet[columnName] {
		return false
	}

	if tablet.excludeSet[columnName] {
		return false
	}

	return true
}

// AddColumn add a column to the output.
func (tablet *tabular) addColumn(name string, rightAlign bool, wrapFmt string) (colID, error) {
	// Add a column to end
	c := len(tablet.columns)

	nameLen := len(name)
	minWidth := nameLen

	if wrapFmt != "" {
		w, err := strconv.ParseInt(wrapFmt, 10, 32)
		if err != nil {
			return colID(c), err
		}
		if int(w) > minWidth {
			minWidth = int(w)
		}
	}

	tablet.columns = append(tablet.columns, &column{
		name:       name,
		rightAlign: rightAlign,
		width:      nameLen,
		minWidth:   minWidth, // minWidth is the smallest width of the column
	})

	return colID(c), nil
}

// newRow add a new row and return its row id.
func (tablet *tabular) newRow() rowID {
	c := len(tablet.columns)

	r := len(tablet.rows)

	tablet.rows = append(tablet.rows, make([]string, c))

	return rowID(r)
}

// setField set a field value.
func (tablet *tabular) setField(row rowID, column colID, format string, value ...interface{}) error {
	if column < 0 || int(column) >= len(tablet.columns) {
		return lpax.Errorf(langpack.ErrorColumnInvalidID, column)
	}
	if row < 0 || int(row) >= len(tablet.rows) {
		return lpax.Errorf(langpack.ErrorRowInvalidID, row)
	}

	s := fmt.Sprintf(format, value...)

	tablet.rows[row][column] = s

	w := tablet.columns[column].width

	if w < len(s) {
		tablet.columns[column].width = len(s)
	}

	return nil
}

// write output the table to an io writer using the supplied table style.
func (tablet *tabular) write(out io.Writer, style TableStyle, excludeHeader bool, columnSeparator string, terminalWidth int) error {
	switch style {
	case Plain:
		return tablet.writePlain(out, excludeHeader, columnSeparator)

	case Aligned:
		return tablet.writeAligned(out, excludeHeader, false, 0, terminalWidth)

	case Grid:
		return tablet.writeAligned(out, excludeHeader, true, 1, terminalWidth)

	case Markdown:
		return tablet.writeMarkdown(out)

	default:
		return lpax.Errorf(langpack.ErrorUnknownStyle, style)
	}
}

type wrappedTable struct {
	totalSpacing     int
	spacing          []int
	rowsLinesColumns [][][]string
}

func (tablet *tabular) calcWidths(minSpacing, terminalWidth int) {
	c := len(tablet.columns)

	//	What s needed
	needed := make([]int, c)
	var totalNeeded int
	for i, col := range tablet.columns {
		needs := col.width - col.minWidth
		if needs > 0 {
			needed[i] = needs
			totalNeeded += needs
		}
	}

	for available := terminalWidth - minSpacing; available > 0 && totalNeeded > 0; {
		for i, col := range tablet.columns {
			if n := needed[i]; n > 0 {
				// column widens
				needed[i]--
				col.minWidth++
				col.width = col.minWidth

				// available shrinks
				available--
				totalNeeded--
			}
		}
	}
}

func (tablet *tabular) calculateSpacing(hasGrid bool, pad int) (spacing []int, totalSpacing int, minSpacing int) {
	c := len(tablet.columns)

	// spacing is the width of each column
	spacing = make([]int, c)

	// extra specing is the left and right padding
	extraSpacing := (pad * 2)

	// if the output is a grid add in extra charact for final grid close char
	if hasGrid {
		minSpacing, totalSpacing = 1, 1
	}

	for i, col := range tablet.columns {
		spacing[i] = col.width + extraSpacing
		totalSpacing += col.width + extraSpacing + 1
		minSpacing += col.minWidth + extraSpacing + 1
	}

	return
}

func wrapText(text string, width int) []string {
	wrapped := wordwrap.Wrapper(width, true)(text)
	return strings.Split(wrapped, "\n")
}

func (tablet *tabular) buildWrappedTable(hasGrid bool, pad int, terminalWidth int) *wrappedTable {
	rows := make([][][]string, len(tablet.rows))

	spacing, totalSpacing, minSpacing := tablet.calculateSpacing(hasGrid, pad)

	if terminalWidth > 0 && totalSpacing > terminalWidth {
		tablet.calcWidths(minSpacing, terminalWidth)
		// recalculate spacings
		spacing, totalSpacing, _ = tablet.calculateSpacing(hasGrid, pad)

		// Build output in in row column line order and then pivot
		c := len(tablet.columns)
		for rowID, r := range tablet.rows {
			linesPerRow := 1
			columnLines := make([][]string, c)
			for fID, field := range r {
				width := tablet.columns[fID].minWidth
				lines := wrapText(field, width)
				if len(lines) > linesPerRow {
					linesPerRow = len(lines)
				}
				columnLines[fID] = lines
			}

			// Pivot
			lineColumns := make([][]string, linesPerRow)
			for j := 0; j < linesPerRow; j++ {
				cols := make([]string, c)
				for i, linesPerCol := range columnLines {
					if j >= len(linesPerCol) {
						cols[i] = ""
					} else {
						cols[i] = linesPerCol[j]
					}
				}
				lineColumns[j] = cols
			}
			rows[rowID] = lineColumns
		}
	} else {
		// sinple
		for rowID, r := range tablet.rows {
			l := []([]string){r}
			rows[rowID] = l
		}
	}

	wrapTable := &wrappedTable{
		totalSpacing:     totalSpacing,
		spacing:          spacing,
		rowsLinesColumns: rows,
	}

	return wrapTable
}

func (tablet *tabular) writeAligned(out io.Writer, excludeHeader bool, hasGrid bool, pad int, terminalWidth int) error {
	// Column width aligned output
	if len(tablet.rows) == 0 {
		return nil
	}

	wrappedTable := tablet.buildWrappedTable(hasGrid, pad, terminalWidth)

	// Write opening grid line
	if hasGrid {
		if err := writeGridLine(out, wrappedTable.totalSpacing,
			wrappedTable.spacing); err != nil {
			return err
		}
	}

	// Write header
	if !excludeHeader {
		if err := tablet.writeAlignedHeader(out, hasGrid, pad,
			wrappedTable.totalSpacing); err != nil {
			return err
		}

		// add in header grid line
		if hasGrid {
			if err := writeGridLine(out, wrappedTable.totalSpacing,
				wrappedTable.spacing); err != nil {
				return err
			}
		}
	}

	// Walk through rows
	for _, row := range wrappedTable.rowsLinesColumns {
		// output row
		if err := tablet.writeAlignedRow(out, row, hasGrid, pad,
			wrappedTable.totalSpacing); err != nil {
			return err
		}

		// add in grid line
		if hasGrid {
			if err := writeGridLine(out, wrappedTable.totalSpacing, wrappedTable.spacing); err != nil {
				return err
			}
		}
	}

	return nil
}

func (tablet *tabular) writePlain(out io.Writer, excludeHeader bool, columnSeparator string) error {
	// basic raw output, using a separator between columns

	if !excludeHeader {
		if err := tablet.writePlainHeader(out, columnSeparator); err != nil {
			return err
		}
	}

	for _, row := range tablet.rows {
		if err := tablet.writePlainRow(out, row, columnSeparator); err != nil {
			return err
		}
	}

	return nil
}

func (tablet *tabular) writeMarkdown(out io.Writer) error {
	// basic raw output, using a separator between columns

	if err := tablet.writeMarkdownHeader(out); err != nil {
		return err
	}

	for _, row := range tablet.rows {
		if err := tablet.writeMarkdownRow(out, row); err != nil {
			return err
		}
	}

	return nil
}

func (tablet *tabular) writeMarkdownHeader(out io.Writer) error {
	if len(tablet.columns) == 0 {
		return nil
	}

	// add fin line return
	//nolint:errcheck
	defer out.Write([]byte("\n"))

	for _, col := range tablet.columns {
		_, err := out.Write([]byte("|" + col.name))
		if err != nil {
			return err
		}
	}
	_, err := out.Write([]byte("|\n"))
	if err != nil {
		return err
	}
	_, err = out.Write([]byte(strings.Repeat("|-", len(tablet.columns))))
	if err != nil {
		return err
	}
	_, err = out.Write([]byte("|"))
	if err != nil {
		return err
	}
	return nil
}

func (tablet *tabular) writeMarkdownRow(out io.Writer, row []string) error {
	if len(row) == 0 {
		return nil
	}

	// add in line return on completion / error
	//nolint:errcheck
	defer out.Write([]byte("\n"))

	for _, field := range row {
		_, err := out.Write([]byte("|" + field))
		if err != nil {
			return err
		}
	}

	_, err := out.Write([]byte("|"))
	if err != nil {
		return err
	}
	return nil
}

func (tablet *tabular) writePlainHeader(out io.Writer, columnSeparator string) error {
	if len(tablet.columns) == 0 {
		return nil
	}

	// add fin line return
	//nolint:errcheck
	defer out.Write([]byte("\n"))

	var sep string
	for i, col := range tablet.columns {
		_, err := out.Write([]byte(sep + col.name))
		if err != nil {
			return err
		}
		if i == 0 {
			sep = columnSeparator
		}
	}

	return nil
}

func (tablet *tabular) writePlainRow(out io.Writer, row []string, columnSeparator string) error {
	if len(row) == 0 {
		return nil
	}

	// add in line return on completion / error
	//nolint:errcheck
	defer out.Write([]byte("\n"))

	var sep string
	for i, field := range row {
		_, err := out.Write([]byte(sep + field))
		if err != nil {
			return err
		}
		if i == 0 {
			sep = columnSeparator
		}
	}

	return nil
}

func (tablet *tabular) writeAlignedRow(out io.Writer, row [][]string, hasGrid bool, padding, total int) error {
	if len(tablet.columns) == 0 {
		return nil
	}

	// output the row
	var b strings.Builder
	b.Grow(total * len(row))

	for _, line := range row {
		for i, field := range line {
			col := tablet.columns[i]
			w := col.width
			fill := w - len(field)

			if hasGrid {
				b.WriteString("|")
			} else if i > 0 {
				b.WriteString(" ")
			}

			// add padding to left or right
			if col.rightAlign {
				b.WriteString(strings.Repeat(" ", fill+padding))
				b.WriteString(field)
				b.WriteString(strings.Repeat(" ", padding))
			} else {
				b.WriteString(strings.Repeat(" ", padding))
				b.WriteString(field)
				b.WriteString(strings.Repeat(" ", fill+padding))
			}
		}
		// complete grid
		if hasGrid {
			b.WriteString("|")
		}

		// add in close
		b.WriteString("\n")
	}

	_, err := out.Write([]byte(b.String()))

	return err
}

func (tablet *tabular) writeAlignedHeader(out io.Writer, hasGrid bool, padding, total int) error {
	// write out the aligned header
	if len(tablet.columns) == 0 {
		return nil
	}

	var b strings.Builder
	b.Grow(total)

	for i, col := range tablet.columns {
		w := col.width
		fill := w - len(col.name)

		if hasGrid {
			b.WriteString("|")
		} else if i > 0 {
			b.WriteString(" ")
		}

		if col.rightAlign {
			b.WriteString(strings.Repeat(" ", fill+padding))
			b.WriteString(col.name)
			b.WriteString(strings.Repeat(" ", padding))
		} else {
			b.WriteString(strings.Repeat(" ", padding))
			b.WriteString(col.name)
			b.WriteString(strings.Repeat(" ", fill+padding))
		}
	}

	if hasGrid {
		b.WriteString("|")
	}

	b.WriteString("\n")

	_, err := out.Write([]byte(b.String()))

	return err
}

func writeGridLine(out io.Writer, total int, spacing []int) error {
	// output a simple grid line
	var b strings.Builder
	b.Grow(total)

	b.WriteString("+")

	for _, c := range spacing {
		b.WriteString(strings.Repeat("-", c) + "+")
	}

	b.WriteString("\n")

	_, err := out.Write([]byte(b.String()))

	return err
}
