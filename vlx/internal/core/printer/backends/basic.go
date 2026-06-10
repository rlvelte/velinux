package backends

import (
	"fmt"
	"os"
)

// Basic is a fallback backend that uses plain fmt and os operations.
type Basic struct{}

// Available reports whether this backend can be used.
func (b *Basic) Available() bool {
	return true
}

// Info prints an info message.
func (b *Basic) Info(msg string) {
	fmt.Println("[INFO]", msg)
}

// Success prints a success message.
func (b *Basic) Success(msg string) {
	fmt.Println("[OK]", msg)
}

// Warn prints a warning message.
func (b *Basic) Warn(msg string) {
	fmt.Println("[WARN]", msg)
}

// Error prints an error message.
func (b *Basic) Error(msg string) {
	fmt.Fprintln(os.Stderr, "[ERR]", msg)
}

// Header signals a new section.
func (b *Basic) Header(msg string) {
	fmt.Println("==>", msg)
}

// Table prints data in a tabular format.
func (b *Basic) Table(headers []string, rows [][]string) {
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
func (b *Basic) Confirm(msg string) bool {
	fmt.Printf("%s [y/N]: ", msg)
	var response string

	fmt.Scanln(&response)
	return response == "y" || response == "Y"
}

// Spinner shows a progress indicator.
func (b *Basic) Spinner(label string, fn func() error) error {
	fmt.Printf("%s...", label)
	err := fn()

	if err != nil {
		fmt.Println(" failed")
	} else {
		fmt.Println(" done")
	}

	return err
}
