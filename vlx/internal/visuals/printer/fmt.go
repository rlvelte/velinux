package printer

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// FmtPrinter is a fallback backend that uses plain fmt and os operations.
type FmtPrinter struct{}

// Print prints a message.
func (f *FmtPrinter) Print(msg string) {
	fmt.Println(msg)
}

// Warn prints a warning message.
func (f *FmtPrinter) Warn(msg string) {
	fmt.Println("[WARN]", msg)
}

// Error prints an error message.
func (f *FmtPrinter) Error(msg string) {
	_, _ = fmt.Fprintln(os.Stderr, "[ERR]", msg)
}

const tabWidth = 8

// Table prints a table
func (f *FmtPrinter) Table(headers []string, rows [][]string) {
	if len(headers) == 0 && len(rows) == 0 {
		return
	}

	// Determine number of columns
	numCols := len(headers)
	for _, row := range rows {
		if len(row) > numCols {
			numCols = len(row)
		}
	}
	if numCols == 0 {
		return
	}

	// Calculate max width per column across headers and all rows
	maxWidths := make([]int, numCols)
	if len(headers) > 0 {
		for i, h := range headers {
			maxWidths[i] = len(h)
		}
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < numCols && len(cell) > maxWidths[i] {
				maxWidths[i] = len(cell)
			}
		}
	}

	printRow := func(cells []string) {
		for i, cell := range cells {
			if i > 0 {
				// Calculate how many tabs are needed to reach the target
				// tab stop for the previous column's maximum width.
				targetTab := maxWidths[i-1]/tabWidth + 1
				currentTab := len(cells[i-1]) / tabWidth
				tabs := targetTab - currentTab
				if tabs < 1 {
					tabs = 1
				}
				for j := 0; j < tabs; j++ {
					fmt.Print("\t")
				}
			}
			fmt.Print(cell)
		}
		fmt.Println()
	}

	if len(headers) > 0 {
		printRow(headers)
	}
	for _, row := range rows {
		printRow(row)
	}
}

// Confirm shows a simple confirmation dialog.
func (f *FmtPrinter) Confirm(msg string, defaultYes bool) bool {
	if defaultYes {
		fmt.Printf("%s [Y/n]: ", msg)
	} else {
		fmt.Printf("%s [y/N]: ", msg)
	}

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return false
	}

	response := strings.TrimSpace(scanner.Text())
	return response == "y" || response == "Y" || (defaultYes && response == "")
}
