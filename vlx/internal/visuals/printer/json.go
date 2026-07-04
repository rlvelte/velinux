package printer

import (
	"encoding/json"
	"os"
)

type jsonMessage struct {
	Level   string `json:"level,omitempty"`
	Message string `json:"message"`
}

type jsonTable struct {
	Headers []string   `json:"headers"`
	Rows    [][]string `json:"rows"`
}

// JSONPrinter primarily for consumption by other tools.
type JSONPrinter struct {
	encoder *json.Encoder
}

// Print returns a message.
func (b *JSONPrinter) Print(msg string) {
	_ = b.encoder.Encode(jsonMessage{Message: msg})
}

// Warn returns a warning message.
func (b *JSONPrinter) Warn(msg string) {
	_ = b.encoder.Encode(jsonMessage{Level: "warn", Message: msg})
}

// Error returns an error message.
func (b *JSONPrinter) Error(msg string) {
	_ = json.NewEncoder(os.Stderr).Encode(jsonMessage{Level: "error", Message: msg})
}

// Table returns a table.
func (b *JSONPrinter) Table(headers []string, rows [][]string) {
	_ = b.encoder.Encode(jsonTable{Headers: headers, Rows: rows})
}

// Confirm does not work in JSON mode.
func (b *JSONPrinter) Confirm(_ string, _ bool) bool {
	return false
}
