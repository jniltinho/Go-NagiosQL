// Package nagimport parses Nagios .cfg files and imports objects into the database.
package nagimport

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ParsedObject holds one parsed Nagios object definition.
type ParsedObject struct {
	Type   string
	Fields map[string]string
}

// ParseFile reads a .cfg file and returns all define blocks it contains.
// Unknown directive lines (outside define blocks) are skipped silently.
func ParseFile(path string) ([]ParsedObject, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening %s: %w", path, err)
	}
	defer f.Close()

	var objects []ParsedObject
	var current *ParsedObject

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and blank lines.
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		if strings.HasPrefix(line, "define ") && strings.HasSuffix(line, "{") {
			objType := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, "define "), "{"))
			current = &ParsedObject{Type: strings.TrimSpace(objType), Fields: make(map[string]string)}
			continue
		}

		if line == "}" && current != nil {
			objects = append(objects, *current)
			current = nil
			continue
		}

		if current != nil {
			// Split on first whitespace.
			parts := strings.SplitN(line, " ", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				val := strings.TrimSpace(parts[1])
				// Strip inline comments.
				if idx := strings.Index(val, ";"); idx != -1 {
					val = strings.TrimSpace(val[:idx])
				}
				current.Fields[key] = val
			}
		}
	}
	return objects, scanner.Err()
}
