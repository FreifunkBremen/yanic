package webserver

import (
	"errors"
	"fmt"
	"strings"
)

type PrometheusMetric struct {
	Name   string
	Value  interface{}
	Labels map[string]interface{}
}

func (m *PrometheusMetric) String() (string, error) {
	if m.Value == nil {
		return "", errors.New("no value of metric found")
	}
	output := m.Name
	if len(m.Labels) > 0 {
		output += "{"
		for label, v := range m.Labels {
			switch value := v.(type) {
			case string:
				output = fmt.Sprintf("%s%s=\"%s\",", output, label, strings.ReplaceAll(value, "\"", "'"))
			case float32:
				output = fmt.Sprintf("%s%s=\"%.4f\",", output, label, value)
			case float64:
				output = fmt.Sprintf("%s%s=\"%.4f\",", output, label, value)
			default:
				output = fmt.Sprintf("%s%s=\"%v\",", output, label, value)
			}
		}
		lastChar := len(output) - 1
		output = output[:lastChar] + "}"
	}
	return fmt.Sprintf("%s %v", output, m.Value), nil
}
