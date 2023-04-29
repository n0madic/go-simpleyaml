# go-simpleyaml

[![Go Reference](https://pkg.go.dev/badge/github.com/n0madic/go-simpleyaml.svg)](https://pkg.go.dev/github.com/n0madic/go-simpleyaml)


Package *simpleyaml* provides a simple function to parse YAML files.

Advantages:
- Simple and fast
- One file (`yaml.go`)
- Easy integration into the project
- No dependencies
- No reflection
- Supported value types: `string`, `int64`, `float64`, `bool`, `list`, `map`
- Access to the value by path (e.g. `yamlMap.Path("root[0].key")`)

Disadvantages:
- No support for YAML versions
- No support for anchors
- No support for tags
- No errors handling

## Why?

I needed a simple and fast way to parse YAML files in my project. I did not find a suitable solution, so I wrote my own.

## Installation

```bash
go get -u github.com/n0madic/go-simpleyaml
```

or move file `yaml.go` to your project.

## Usage

```go
package main

import (
	"fmt"

	simpleyaml "github.com/n0madic/go-simpleyaml"
)

func main() {
	yamlStr := `---
listen:
    address: 127.0.0.1
    ports:
    - 1555
    - 2222
paths:
- path: /home
  name: Home
- path: /tmp
  name: Temp
`
	yamlMap := simpleyaml.ParseYAML(yamlStr)

	// Get the value by path
	fmt.Printf("%s:%d\n", yamlMap.Path("listen.address"), yamlMap.Path("listen.ports[0]"))

	// Get the value from map by key and type assertion
	paths := yamlMap["paths"].([]interface{})
	for _, path := range paths {
		pathMap := path.(simpleyaml.YAMLNode)
		fmt.Printf("%s: %s\n", pathMap["name"], pathMap["path"])
	}
}
```

## Benchmark

Benchmark versus [gopkg.in/yaml.v3](https://gopkg.in/yaml.v3)

```bash
goos: darwin
goarch: arm64
pkg: github.com/n0madic/go-simpleyaml
BenchmarkParseYAML
BenchmarkParseYAML-10          	  184214	      6531 ns/op	    6649 B/op	     156 allocs/op
BenchmarkYAMLV3Unmarshal
BenchmarkYAMLV3Unmarshal-10    	   38187	     31400 ns/op	   21105 B/op	     352 allocs/op
```
