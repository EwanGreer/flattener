package flattener_test

import (
	"testing"

	"github.com/EwanGreer/flattener"
	"github.com/stretchr/testify/assert"
)

func TestFlattener_JSON(t *testing.T) {
	tests := []struct {
		name      string
		delimiter string
		input     []byte
		want      []byte
		wantErr   bool
	}{
		{
			name:      "UnNested Key",
			delimiter: ".",
			input:     []byte(`{"name":"james"}`),
			want:      []byte(`{"name":"james"}`),
			wantErr:   false,
		},
		{
			name:      "Nil Input",
			delimiter: ".",
			input:     nil,
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "Nested Object",
			delimiter: ".",
			input:     []byte(`{"user":{"name":"john","age":30}}`),
			want:      []byte(`{"user.age":30,"user.name":"john"}`),
			wantErr:   false,
		},
		{
			name:      "Deeply Nested Object",
			delimiter: ".",
			input:     []byte(`{"user":{"address":{"city":"NYC","zip":"10001"}}}`),
			want:      []byte(`{"user.address.city":"NYC","user.address.zip":"10001"}`),
			wantErr:   false,
		},
		{
			name:      "Array",
			delimiter: ".",
			input:     []byte(`{"items":["a","b","c"]}`),
			want:      []byte(`{"items.0":"a","items.1":"b","items.2":"c"}`),
			wantErr:   false,
		},
		{
			name:      "Nested Array with Objects",
			delimiter: ".",
			input:     []byte(`{"users":[{"name":"alice"},{"name":"bob"}]}`),
			want:      []byte(`{"users.0.name":"alice","users.1.name":"bob"}`),
			wantErr:   false,
		},
		{
			name:      "Mixed Types",
			delimiter: ".",
			input:     []byte(`{"name":"test","count":42,"active":true,"data":null}`),
			want:      []byte(`{"active":true,"count":42,"data":null,"name":"test"}`),
			wantErr:   false,
		},
		{
			name:      "Custom Delimiter",
			delimiter: "_",
			input:     []byte(`{"user":{"name":"jane"}}`),
			want:      []byte(`{"user_name":"jane"}`),
			wantErr:   false,
		},
		{
			name:      "Empty Object",
			delimiter: ".",
			input:     []byte(`{}`),
			want:      []byte(`{}`),
			wantErr:   false,
		},
		{
			name:      "Invalid JSON",
			delimiter: ".",
			input:     []byte(`{invalid}`),
			want:      nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := flattener.Flattener{
				Delimeter: tt.delimiter,
			}

			result, err := f.JSON(tt.input)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("JSON() failed: %v", err)
				}

				return
			}

			if tt.wantErr {
				t.Fatal("JSON() succeeded unexpectedly")
			}

			assert.NotNil(t, result, "result is nil")
			assert.Equal(t, tt.want, result, "objects do not match")
		})
	}
}

func TestFlattener_flatten(t *testing.T) {
	tests := []struct {
		name      string
		delimiter string
		prefix    string
		input     any
		want      map[string]any
	}{
		{
			name:      "Primitive String",
			delimiter: ".",
			prefix:    "key",
			input:     "value",
			want:      map[string]any{"key": "value"},
		},
		{
			name:      "Primitive Number",
			delimiter: ".",
			prefix:    "count",
			input:     42,
			want:      map[string]any{"count": 42},
		},
		{
			name:      "Primitive Boolean",
			delimiter: ".",
			prefix:    "active",
			input:     true,
			want:      map[string]any{"active": true},
		},
		{
			name:      "Primitive Null",
			delimiter: ".",
			prefix:    "data",
			input:     nil,
			want:      map[string]any{"data": nil},
		},
		{
			name:      "Map with Empty Prefix",
			delimiter: ".",
			prefix:    "",
			input:     map[string]any{"name": "john", "age": 30},
			want:      map[string]any{"name": "john", "age": 30},
		},
		{
			name:      "Map with Prefix",
			delimiter: ".",
			prefix:    "user",
			input:     map[string]any{"name": "john", "age": 30},
			want:      map[string]any{"user.name": "john", "user.age": 30},
		},
		{
			name:      "Nested Map",
			delimiter: ".",
			prefix:    "",
			input:     map[string]any{"user": map[string]any{"name": "john"}},
			want:      map[string]any{"user.name": "john"},
		},
		{
			name:      "Array with Empty Prefix",
			delimiter: ".",
			prefix:    "",
			input:     []any{"a", "b", "c"},
			want:      map[string]any{"0": "a", "1": "b", "2": "c"},
		},
		{
			name:      "Array with Prefix",
			delimiter: ".",
			prefix:    "items",
			input:     []any{"a", "b", "c"},
			want:      map[string]any{"items.0": "a", "items.1": "b", "items.2": "c"},
		},
		{
			name:      "Array of Objects",
			delimiter: ".",
			prefix:    "users",
			input:     []any{map[string]any{"name": "alice"}, map[string]any{"name": "bob"}},
			want:      map[string]any{"users.0.name": "alice", "users.1.name": "bob"},
		},
		{
			name:      "Custom Delimiter",
			delimiter: "_",
			prefix:    "user",
			input:     map[string]any{"name": "jane"},
			want:      map[string]any{"user_name": "jane"},
		},
		{
			name:      "Empty Map",
			delimiter: ".",
			prefix:    "",
			input:     map[string]any{},
			want:      map[string]any{},
		},
		{
			name:      "Empty Array",
			delimiter: ".",
			prefix:    "items",
			input:     []any{},
			want:      map[string]any{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := flattener.Flattener{
				Delimeter: tt.delimiter,
			}

			result := make(map[string]any)
			f.Flatten(tt.prefix, tt.input, result)

			assert.Equal(t, tt.want, result, "flattened result does not match")
		})
	}
}

func TestFlattener_JSON_Sorting(t *testing.T) {
	f := flattener.Flattener{
		Delimeter: ".",
	}

	// Test with input that would have random order if not sorted
	input := []byte(`{"z":{"b":"value1","a":"value2"},"a":{"z":"value3","a":"value4"}}`)

	result, err := f.JSON(input)
	assert.NoError(t, err)

	// The keys should be in alphabetical order
	expected := []byte(`{"a.a":"value4","a.z":"value3","z.a":"value2","z.b":"value1"}`)
	assert.Equal(t, expected, result)
}

func TestFlattener_YAML(t *testing.T) {
	tests := []struct {
		name      string
		delimiter string
		input     []byte
		want      string
		wantErr   bool
	}{
		{
			name:      "Simple YAML",
			delimiter: ".",
			input:     []byte("name: james\n"),
			want:      "name: james\n",
			wantErr:   false,
		},
		{
			name:      "Nil Input",
			delimiter: ".",
			input:     nil,
			want:      "",
			wantErr:   true,
		},
		{
			name:      "Nested YAML",
			delimiter: ".",
			input:     []byte("user:\n  name: john\n  age: 30\n"),
			want:      "user.age: 30\nuser.name: john\n",
			wantErr:   false,
		},
		{
			name:      "YAML with Array",
			delimiter: ".",
			input:     []byte("items:\n  - a\n  - b\n  - c\n"),
			want:      "items.0: a\nitems.1: b\nitems.2: c\n",
			wantErr:   false,
		},
		{
			name:      "Custom Delimiter",
			delimiter: "_",
			input:     []byte("user:\n  name: jane\n"),
			want:      "user_name: jane\n",
			wantErr:   false,
		},
		{
			name:      "Invalid YAML",
			delimiter: ".",
			input:     []byte("invalid: yaml: structure: :\n"),
			want:      "",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := flattener.Flattener{
				Delimeter: tt.delimiter,
			}

			result, err := f.YAML(tt.input)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("YAML() failed: %v", err)
				}

				return
			}

			if tt.wantErr {
				t.Fatal("YAML() succeeded unexpectedly")
			}

			assert.NotNil(t, result, "result is nil")
			assert.Equal(t, tt.want, string(result), "YAML output does not match")
		})
	}
}

func TestFlattener_YAML_Sorting(t *testing.T) {
	f := flattener.Flattener{
		Delimeter: ".",
	}

	input := []byte("z:\n  b: value1\n  a: value2\na:\n  z: value3\n  a: value4\n")

	result, err := f.YAML(input)
	assert.NoError(t, err)

	expected := "a.a: value4\na.z: value3\nz.a: value2\nz.b: value1\n"
	assert.Equal(t, expected, string(result))
}
