package output

import (
	"encoding/json"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

// Format represents the output format.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
	FormatYAML Format = "yaml"
)

// ParseFormat parses a string into a Format.
func ParseFormat(s string) Format {
	switch s {
	case "json":
		return FormatJSON
	case "yaml":
		return FormatYAML
	default:
		return FormatText
	}
}

// Write outputs data in the specified format.
// For text format, textFn is called to produce styled output.
// For json/yaml, data is marshaled from the structured value.
func Write(w io.Writer, format Format, data interface{}, textFn func(io.Writer)) {
	switch format {
	case FormatJSON:
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		enc.Encode(data)
	case FormatYAML:
		enc := yaml.NewEncoder(w)
		enc.SetIndent(2)
		enc.Encode(data)
		enc.Close()
	default:
		textFn(w)
	}
}

// Println writes a line to the writer.
func Println(w io.Writer, s string) {
	fmt.Fprintln(w, s)
}
