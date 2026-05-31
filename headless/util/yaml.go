package util

import (
	"bufio"
	"fmt"
	"strings"
)

func Marshal(v interface{}) ([]byte, error) {
	switch m := v.(type) {
	case map[string]interface{}:
		return []byte(marshalMap(m, 0)), nil
	default:
		return []byte(fmt.Sprintf("%v", v)), nil
	}
}

func marshalMap(m map[string]interface{}, indent int) string {
	var result strings.Builder
	prefix := strings.Repeat("  ", indent)
	for k, v := range m {
		switch val := v.(type) {
		case map[string]interface{}:
			result.WriteString(fmt.Sprintf("%s%s:\n", prefix, k))
			result.WriteString(marshalMap(val, indent+1))
		case []interface{}:
			result.WriteString(fmt.Sprintf("%s%s:\n", prefix, k))
			for _, item := range val {
				if s, ok := item.(string); ok {
					result.WriteString(fmt.Sprintf("%s  - %s\n", prefix, s))
				} else {
					result.WriteString(fmt.Sprintf("%s  - %v\n", prefix, item))
				}
			}
		case string:
			result.WriteString(fmt.Sprintf("%s%s: %s\n", prefix, k, val))
		case bool:
			result.WriteString(fmt.Sprintf("%s%s: %t\n", prefix, k, val))
		case int:
			result.WriteString(fmt.Sprintf("%s%s: %d\n", prefix, k, val))
		case float64:
			result.WriteString(fmt.Sprintf("%s%s: %v\n", prefix, k, val))
		case nil:
			result.WriteString(fmt.Sprintf("%s%s: null\n", prefix, k))
		default:
			result.WriteString(fmt.Sprintf("%s%s: %v\n", prefix, k, val))
		}
	}
	return result.String()
}

func Unmarshal(data []byte, v interface{}) error {
	m, ok := v.(*map[string]interface{})
	if !ok {
		return fmt.Errorf("unsupported type")
	}
	*m = parseYAML(string(data))
	return nil
}

func parseYAML(s string) map[string]interface{} {
	result := make(map[string]interface{})
	scanner := bufio.NewScanner(strings.NewReader(s))
	var currentKey string

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		if strings.HasPrefix(trimmed, "- ") {
			item := strings.TrimPrefix(trimmed, "- ")
			if existing, ok := result[currentKey]; ok {
				if list, ok2 := existing.([]string); ok2 {
					result[currentKey] = append(list, item)
				}
			} else {
				result[currentKey] = []string{item}
			}
			continue
		}

		colonIdx := strings.Index(trimmed, ":")
		if colonIdx == -1 {
			continue
		}
		currentKey = strings.TrimSpace(trimmed[:colonIdx])
		rest := strings.TrimSpace(trimmed[colonIdx+1:])

		if rest == "" {
			continue
		}
		result[currentKey] = parseYAMLValue(rest)
	}
	return result
}

func parseYAMLValue(s string) interface{} {
	s = strings.TrimSpace(s)
	if s == "true" || s == "yes" {
		return true
	}
	if s == "false" || s == "no" {
		return false
	}
	if s == "null" || s == "~" {
		return nil
	}
	return s
}
