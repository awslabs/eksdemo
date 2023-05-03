package printer

import (
	"encoding/json"
	"fmt"
	"io"

	"sigs.k8s.io/yaml"
)

// Switched to sigs.k8s.io/yaml to fix panic when outputting YAML using AWS SDK v2
// Issue: https://github.com/go-yaml/yaml/issues/463

type Output string

const (
	JSON  Output = "json"
	Table Output = "table"
	YAML  Output = "yaml"
)

type ResourcePrinter interface {
	PrintJSON(io.Writer) error
	PrintTable(io.Writer) error
	PrintYAML(io.Writer) error
}

func EncodeJSON(w io.Writer, obj interface{}) error {
	out, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %s", err)
	}

	_, err = w.Write(out)
	if err != nil {
		return fmt.Errorf("failed to output JSON: %s", err)
	}
	return nil
}

func EncodeYAML(w io.Writer, obj interface{}) error {
	out, err := yaml.Marshal(obj)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %s", err)
	}

	_, err = w.Write(out)
	if err != nil {
		return fmt.Errorf("failed to output YAML: %s", err)
	}
	return nil
}

func NewOutput(outputType string) (Output, error) {
	switch outputType {
	case string(Table):
		return Table, nil
	case string(JSON):
		return JSON, nil
	case string(YAML):
		return YAML, nil
	default:
		return "", fmt.Errorf("unsupported output type: %q", outputType)
	}
}

func (o Output) Print(out io.Writer, printer ResourcePrinter) error {
	switch o {
	default:
	case Table:
		return printer.PrintTable(out)
	case JSON:
		return printer.PrintJSON(out)
	case YAML:
		return printer.PrintYAML(out)
	}
	return nil
}
