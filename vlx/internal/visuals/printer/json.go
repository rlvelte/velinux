package printer

import (
	"encoding/json"
	"os"
)

type jsonMessage struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type jsonTable struct {
	Headers []string   `json:"headers"`
	Rows    [][]string `json:"rows"`
}

// JsonPrinter primarily for consumption by other tools.
type JsonPrinter struct {
	encoder *json.Encoder
}

// Success returns a message.
func (j *JsonPrinter) Success(msg string) {
	_ = j.encoder.Encode(jsonMessage{Status: "success", Message: msg})
}

// Error returns an error message.
func (j *JsonPrinter) Error(err error) {
	_ = json.NewEncoder(os.Stderr).Encode(jsonMessage{Status: "error", Message: err.Error()})
}

// Table returns a table.
func (j *JsonPrinter) Table(headers []string, rows [][]string) {
	_ = j.encoder.Encode(jsonTable{Headers: headers, Rows: rows})
}
