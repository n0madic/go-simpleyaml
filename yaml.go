package simpleyaml

import (
	"strconv"
	"strings"
)

// YAMLNode is a recursive structure for representing YAML data.
type YAMLNode map[string]interface{}

// Path returns a value from YAMLNode by path.
// If the path is not found, it returns nil.
func (n YAMLNode) Path(path string) interface{} {
	parts := strings.Split(path, ".")
	var value interface{} = n
	for _, part := range parts {
		if strings.Contains(part, "[") && strings.Contains(part, "]") {
			// handle indexed path
			key := part[:strings.Index(part, "[")]
			indexStr := part[strings.Index(part, "[")+1 : strings.Index(part, "]")]
			index, err := strconv.Atoi(indexStr)
			if err != nil {
				return nil
			}
			list, ok := n[key].([]interface{})
			if !ok {
				return nil
			}
			if index >= len(list) {
				return nil
			}
			value = list[index]
		} else {
			// handle regular path
			value = n[part]
		}
		if value == nil {
			return nil
		}
		n, _ = value.(YAMLNode)
	}
	return value
}

// ParseYAML parses a YAML string and returns a YAMLNode
// which is a map[string]interface{}.
// YAMLNode is a recursive structure.
func ParseYAML(input string) YAMLNode {
	return parseYAMLLines(strings.Split(input, "\n"))
}

// parseYAMLLines parses YAML lines and returns a YAMLNode
func parseYAMLLines(lines []string) YAMLNode {
	node := make(YAMLNode)
	var key string

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		if line == "" || line[0] == '#' {
			continue
		}

		if line == "---" {
			if i == 0 {
				continue
			} else {
				break
			}
		}

		// parse key-value pair
		kv := strings.SplitN(line, ":", 2)
		if len(kv) == 2 {
			key = strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])

			if value == "" && i+1 < len(lines) {
				if isList(lines[i+1]) {
					var list []interface{}
					i++
					for i < len(lines) {
						line = strings.TrimSpace(lines[i])
						if isList(line) {
							var item interface{}
							if strings.Contains(line, ":") {
								// nested map
								block := lines[i:]
								block[0] = strings.Replace(block[0], "-", " ", 1)
								var childLines int
								item, childLines = getIndentedBlock(block)
								i += childLines - 1
							} else {
								// simple list
								item = parseValue(strings.TrimSpace(strings.TrimPrefix(line, "-")))
							}
							list = append(list, item)
						} else if line == "" {
							break
						} else {
							i--
							break
						}
						i++
					}
					node[key] = list
				} else {
					// parse child node
					childNode, childLines := getIndentedBlock(lines[i+1:])
					i += childLines
					node[key] = childNode
				}
			} else if strings.HasPrefix(value, "|") || strings.HasPrefix(value, ">") {
				// parse multiline block
				var blockLines []string
				i++
				indent := getIndent(lines[i])
				for i < len(lines) {
					if getIndent(lines[i]) < indent {
						break
					}
					line := strings.TrimSpace(lines[i])
					if line == "" {
						if strings.HasPrefix(value, ">") {
							blockLines = append(blockLines, "\n")
						} else {
							blockLines = append(blockLines, "")
						}
					} else {
						blockLines = append(blockLines, line)
					}
					i++
				}
				i--
				if strings.HasPrefix(value, ">") {
					result := strings.Join(blockLines, " ")
					result = strings.Replace(result, " \n ", "\n", -1)
					node[key] = result
				} else {
					node[key] = strings.Join(blockLines, "\n")
				}
			} else {
				node[key] = parseValue(value)
			}
		}
	}

	return node
}

// getIndent returns the number of leading spaces or tabs.
func getIndent(s string) int {
	indent := 0
	for _, c := range s {
		if c == ' ' {
			indent++
		} else if c == '\t' {
			indent += 8
		} else {
			break
		}
	}
	return indent
}

// getIndentedBlock returns a child node and the number of lines
func getIndentedBlock(lines []string) (YAMLNode, int) {
	indent := getIndent(lines[0])
	var blockLines []string
	var line string

	count := 0
	for _, line = range lines {
		if getIndent(line) < indent {
			break
		}
		blockLines = append(blockLines, line)
		count++
	}

	childNode := parseYAMLLines(blockLines)
	return childNode, count
}

// parseValue parses a YAML value and returns as interface{}
func parseValue(value string) interface{} {
	if value == "" {
		return nil
	}
	if strings.ToLower(value) == "true" || strings.ToLower(value) == "false" {
		boolValue, _ := strconv.ParseBool(value)
		return boolValue
	}
	if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
		return intValue
	}
	if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
		return floatValue
	}
	if strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") {
		return value[1 : len(value)-1]
	}
	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
		return value[1 : len(value)-1]
	}
	if strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}") {
		// inline map
		yamlObject := make(YAMLNode)
		pairs := strings.Split(value[1:len(value)-1], ",")
		for _, pair := range pairs {
			parts := strings.SplitN(pair, ":", 2)
			key := strings.TrimSpace(parts[0])
			val := parseValue(strings.TrimSpace(parts[1]))
			yamlObject[key] = val
		}
		return yamlObject
	}
	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		// inline array
		array := make([]interface{}, 0)
		items := strings.Split(value[1:len(value)-1], ",")
		for _, item := range items {
			array = append(array, parseValue(strings.TrimSpace(item)))
		}
		return array
	}
	return value
}

// isList returns true if the line is a list item
func isList(line string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), "-")
}
