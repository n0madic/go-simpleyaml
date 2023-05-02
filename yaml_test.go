package simpleyaml

import (
	"reflect"
	"testing"
)

func TestYAMLNode_Path(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		n    YAMLNode
		args args
		want interface{}
	}{
		{
			name: "empty",
			n:    YAMLNode{},
			args: args{path: ""},
			want: nil,
		},
		{
			name: "simple",
			n:    YAMLNode{"a": "v1"},
			args: args{path: "a"},
			want: "v1",
		},
		{
			name: "wrong path",
			n:    YAMLNode{"a": "v1"},
			args: args{path: "b"},
			want: nil,
		},
		{
			name: "nested",
			n:    YAMLNode{"a": YAMLNode{"b": "v1"}},
			args: args{path: "a.b"},
			want: "v1",
		},
		{
			name: "nested list",
			n:    YAMLNode{"a": []interface{}{YAMLNode{"b": "v1"}}},
			args: args{path: "a[0].b"},
			want: "v1",
		},
		{
			name: "out of range",
			n:    YAMLNode{"a": []interface{}{YAMLNode{"b": "v1"}}},
			args: args{path: "a[1].b"},
			want: nil,
		},
		{
			name: "wrong index",
			n:    YAMLNode{"a": []interface{}{YAMLNode{"b": "v1"}}},
			args: args{path: "a[b].b"},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.Path(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("YAMLNode.Path() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseYAML(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want YAMLNode
	}{
		{
			name: "empty",
			args: args{input: ""},
			want: YAMLNode{},
		},
		{
			name: "simple",
			args: args{input: "a: 1\nb: str\nc:"},
			want: YAMLNode{"a": int64(1), "b": "str", "c": nil},
		},
		{
			name: "bool values",
			args: args{input: "a: true\nb: false"},
			want: YAMLNode{"a": true, "b": false},
		},
		{
			name: "list",
			args: args{input: "a:\n  - 1\n  - 2"},
			want: YAMLNode{"a": []interface{}{int64(1), int64(2)}},
		},
		{
			name: "inline list and map",
			args: args{input: "a: [v1, v2, v3]\nb: {c: v4, d: v5}"},
			want: YAMLNode{"a": []interface{}{"v1", "v2", "v3"}, "b": YAMLNode{"c": "v4", "d": "v5"}},
		},
		{
			name: "quoted string",
			args: args{input: "a: \"v1\"\nb: 'v2'"},
			want: YAMLNode{"a": "v1", "b": "v2"},
		},
		{
			name: "mixed tabs and spaces",
			args: args{input: "a:\n\tb: 1\n\tc:\n        - 2\n        - 3"},
			want: YAMLNode{"a": YAMLNode{"b": int64(1), "c": []interface{}{int64(2), int64(3)}}},
		},
		{
			name: "nested",
			args: args{input: "a:\n  b: 1.1"},
			want: YAMLNode{"a": YAMLNode{"b": 1.1}},
		},
		{
			name: "nested list with map",
			args: args{input: "a:\n  - b: 1\n  - c: 2"},
			want: YAMLNode{"a": []interface{}{YAMLNode{"b": int64(1)}, YAMLNode{"c": int64(2)}}},
		},
		{
			name: "nested list with map and list",
			args: args{input: "a:\n  - b: v1\n  - c:\n    - v2\n    - v3"},
			want: YAMLNode{"a": []interface{}{YAMLNode{"b": "v1"}, YAMLNode{"c": []interface{}{"v2", "v3"}}}},
		},
		{
			name: "mix map and list",
			args: args{input: "a:\n  b: v1\n  c:\n    - v2\n    - v3\n  d:\n    e: v4\n    f: v5"},
			want: YAMLNode{"a": YAMLNode{"b": "v1", "c": []interface{}{"v2", "v3"}, "d": YAMLNode{"e": "v4", "f": "v5"}}},
		},
		{
			name: "list of maps",
			args: args{input: "a:\n  - b: v1\n    c: v2\n  - d: v3\n    e: v4"},
			want: YAMLNode{"a": []interface{}{YAMLNode{"b": "v1", "c": "v2"}, YAMLNode{"d": "v3", "e": "v4"}}},
		},
		{
			name: "list of maps with list",
			args: args{input: "a:\n  - b: v1\n    c:\n      - v2\n      - v3\n  - d: v4\n    e: v5"},
			want: YAMLNode{"a": []interface{}{YAMLNode{"b": "v1", "c": []interface{}{"v2", "v3"}}, YAMLNode{"d": "v4", "e": "v5"}}},
		},
		{
			name: "list of maps with list and map",
			args: args{input: "a:\n  - b: v1\n    c:\n      - v2\n      - v3\n    d:\n      e: v4\n      f: v5\n  - g: v6\n    h: v7"},
			want: YAMLNode{"a": []interface{}{YAMLNode{"b": "v1", "c": []interface{}{"v2", "v3"}, "d": YAMLNode{"e": "v4", "f": "v5"}}, YAMLNode{"g": "v6", "h": "v7"}}},
		},
		{
			name: "list of maps with list and map and list",
			args: args{input: "a:\n  - b: v1\n    c:\n      - v2\n      - v3\n    d:\n      e: v4\n      f:\n        - v5\n        - v6\ng:\n  - v7\n  - v8"},
			want: YAMLNode{"a": []interface{}{YAMLNode{"b": "v1", "c": []interface{}{"v2", "v3"}, "d": YAMLNode{"e": "v4", "f": []interface{}{"v5", "v6"}}}}, "g": []interface{}{"v7", "v8"}},
		},
		{
			name: "maps with maps and list and list",
			args: args{input: "a:\n  b: v1\n  c:\n    d: v2\n    e: v3\n  f:\n    - v4\n    - v5\ng:\n  - v6\n  - v7"},
			want: YAMLNode{"a": YAMLNode{"b": "v1", "c": YAMLNode{"d": "v2", "e": "v3"}, "f": []interface{}{"v4", "v5"}}, "g": []interface{}{"v6", "v7"}},
		},
		{
			name: "multiline block",
			args: args{input: "a: |\n  line1\n  line2\nb: >\n  line3\n  line4"},
			want: YAMLNode{"a": "line1\nline2", "b": "line3 line4"},
		},
		{
			name: "multiline block with empty line",
			args: args{input: "a: |\n  line1\n  \n  line2\nb: >\n  line3\n  \n  line4"},
			want: YAMLNode{"a": "line1\n\nline2", "b": "line3\nline4"},
		},
		{
			name: "divider",
			args: args{input: "---\na: v1\n---\nb: v2"},
			want: YAMLNode{"a": "v1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseYAML(tt.args.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("input:\n%s\nParseYAML() = %v, want %v", tt.args.input, got, tt.want)
			}
		})
	}
}
