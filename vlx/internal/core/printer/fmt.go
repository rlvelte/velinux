package printer

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// FmtPrinter is a fallback backend that uses plain fmt and os operations.
type FmtPrinter struct{}

// Available reports whether this backend can be used.
func (b *FmtPrinter) Available() bool {
	return true
}

// Print prints a message.
func (b *FmtPrinter) Print(msg string) {
	fmt.Println(msg)
}

// Warn prints a warning message.
func (b *FmtPrinter) Warn(msg string) {
	fmt.Println("[WARN]", msg)
}

// Error prints an error message.
func (b *FmtPrinter) Error(msg string) {
	fmt.Fprintln(os.Stderr, "[ERR]", msg)
}

// Table prints data in a tabular format.
func (b *FmtPrinter) Table(headers []string, rows [][]string) {
	if len(headers) > 0 {
		for i, h := range headers {
			if i > 0 {
				fmt.Print("\t")
			}
			fmt.Print(h)
		}
		fmt.Println()
	}
	for _, row := range rows {
		for i, cell := range row {
			if i > 0 {
				fmt.Print("\t")
			}
			fmt.Print(cell)
		}
		fmt.Println()
	}
}

// Confirm shows a simple confirmation dialog.
func (b *FmtPrinter) Confirm(msg string, defaultYes bool) bool {
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

// Spinner shows a progress indicator.
func (b *FmtPrinter) Spinner(label string, fn func() error) error {
	fmt.Printf("%s...", label)
	err := fn()

	if err != nil {
		fmt.Println(" failed")
	} else {
		fmt.Println(" done")
	}

	return err
}
