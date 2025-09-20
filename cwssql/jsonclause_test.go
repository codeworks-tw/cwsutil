package cwssql

import (
	"reflect"
	"testing"
)

func TestJSONContains(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		path     string
		value    any
		expected map[string][]any
	}{
		{
			name:     "basic JSON contains",
			key:      "userData",
			path:     "profile.name",
			value:    "John",
			expected: map[string][]any{"user_data->'profile.name' @> ?": {"John"}},
		},
		{
			name:     "JSON contains with camelCase key",
			key:      "userMetaData",
			path:     "settings.theme",
			value:    "dark",
			expected: map[string][]any{"user_meta_data->'settings.theme' @> ?": {"dark"}},
		},
		{
			name:     "JSON contains with numeric value",
			key:      "config",
			path:     "version",
			value:    2,
			expected: map[string][]any{"config->'version' @> ?": {2}},
		},
		{
			name:     "JSON contains with boolean value",
			key:      "flags",
			path:     "enabled",
			value:    true,
			expected: map[string][]any{"flags->'enabled' @> ?": {true}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test method on existing WhereCaluse
			clause := WhereCaluse{}
			result := clause.JSONContains(tt.key, tt.path, tt.value)
			if !reflect.DeepEqual(map[string][]any(result), tt.expected) {
				t.Errorf("JSONContains() = %v, want %v", result, tt.expected)
			}

			// Test standalone function
			standaloneResult := JSONContains(tt.key, tt.path, tt.value)
			if !reflect.DeepEqual(map[string][]any(standaloneResult), tt.expected) {
				t.Errorf("JSONContains() standalone = %v, want %v", standaloneResult, tt.expected)
			}
		})
	}
}

func TestJSONBContains(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    any
		expected map[string][]any
	}{
		{
			name:     "JSONB contains string",
			key:      "tags",
			value:    "important",
			expected: map[string][]any{"tags @> ?": {"important"}},
		},
		{
			name:     "JSONB contains with camelCase key",
			key:      "userTags",
			value:    "admin",
			expected: map[string][]any{"user_tags @> ?": {"admin"}},
		},
		{
			name:     "JSONB contains object",
			key:      "metadata",
			value:    map[string]string{"type": "user"},
			expected: map[string][]any{"metadata @> ?": {map[string]string{"type": "user"}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test method on existing WhereCaluse
			clause := WhereCaluse{}
			result := clause.JSONBContains(tt.key, tt.value)
			if !reflect.DeepEqual(map[string][]any(result), tt.expected) {
				t.Errorf("JSONBContains() = %v, want %v", result, tt.expected)
			}

			// Test standalone function
			standaloneResult := JSONBContains(tt.key, tt.value)
			if !reflect.DeepEqual(map[string][]any(standaloneResult), tt.expected) {
				t.Errorf("JSONBContains() standalone = %v, want %v", standaloneResult, tt.expected)
			}
		})
	}
}

func TestJSONBContainedBy(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    any
		expected map[string][]any
	}{
		{
			name:     "JSONB contained by",
			key:      "permissions",
			value:    []string{"read", "write", "admin"},
			expected: map[string][]any{"permissions <@ ?": {[]string{"read", "write", "admin"}}},
		},
		{
			name:     "JSONB contained by with camelCase key",
			key:      "userPermissions",
			value:    []string{"read"},
			expected: map[string][]any{"user_permissions <@ ?": {[]string{"read"}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test method on existing WhereCaluse
			clause := WhereCaluse{}
			result := clause.JSONBContainedBy(tt.key, tt.value)
			if !reflect.DeepEqual(map[string][]any(result), tt.expected) {
				t.Errorf("JSONBContainedBy() = %v, want %v", result, tt.expected)
			}

			// Test standalone function
			standaloneResult := JSONBContainedBy(tt.key, tt.value)
			if !reflect.DeepEqual(map[string][]any(standaloneResult), tt.expected) {
				t.Errorf("JSONBContainedBy() standalone = %v, want %v", standaloneResult, tt.expected)
			}
		})
	}
}

func TestJSONExtract(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		path     string
		value    any
		expected map[string][]any
	}{
		{
			name:     "JSON extract string",
			key:      "config",
			path:     "database.host",
			value:    "localhost",
			expected: map[string][]any{"config->'database.host' = ?": {"localhost"}},
		},
		{
			name:     "JSON extract with camelCase key",
			key:      "appConfig",
			path:     "server.port",
			value:    8080,
			expected: map[string][]any{"app_config->'server.port' = ?": {8080}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test method on existing WhereCaluse
			clause := WhereCaluse{}
			result := clause.JSONExtract(tt.key, tt.path, tt.value)
			if !reflect.DeepEqual(map[string][]any(result), tt.expected) {
				t.Errorf("JSONExtract() = %v, want %v", result, tt.expected)
			}

			// Test standalone function
			standaloneResult := JSONExtract(tt.key, tt.path, tt.value)
			if !reflect.DeepEqual(map[string][]any(standaloneResult), tt.expected) {
				t.Errorf("JSONExtract() standalone = %v, want %v", standaloneResult, tt.expected)
			}
		})
	}
}

func TestJSONExtractText(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		path     string
		value    any
		expected map[string][]any
	}{
		{
			name:     "JSON extract text",
			key:      "profile",
			path:     "user.name",
			value:    "John Doe",
			expected: map[string][]any{"profile->>'user.name' = ?": {"John Doe"}},
		},
		{
			name:     "JSON extract text with camelCase key",
			key:      "userProfile",
			path:     "contact.email",
			value:    "john@example.com",
			expected: map[string][]any{"user_profile->>'contact.email' = ?": {"john@example.com"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test method on existing WhereCaluse
			clause := WhereCaluse{}
			result := clause.JSONExtractText(tt.key, tt.path, tt.value)
			if !reflect.DeepEqual(map[string][]any(result), tt.expected) {
				t.Errorf("JSONExtractText() = %v, want %v", result, tt.expected)
			}

			// Test standalone function
			standaloneResult := JSONExtractText(tt.key, tt.path, tt.value)
			if !reflect.DeepEqual(map[string][]any(standaloneResult), tt.expected) {
				t.Errorf("JSONExtractText() standalone = %v, want %v", standaloneResult, tt.expected)
			}
		})
	}
}

func TestJSONArrayContains(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    any
		expected map[string][]any
	}{
		{
			name:     "JSON array contains string",
			key:      "tags",
			value:    "important",
			expected: map[string][]any{"tags ? ?": {"important"}},
		},
		{
			name:     "JSON array contains with camelCase key",
			key:      "userTags",
			value:    "admin",
			expected: map[string][]any{"user_tags ? ?": {"admin"}},
		},
		{
			name:     "JSON array contains number",
			key:      "numbers",
			value:    42,
			expected: map[string][]any{"numbers ? ?": {42}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test method on existing WhereCaluse
			clause := WhereCaluse{}
			result := clause.JSONArrayContains(tt.key, tt.value)
			if !reflect.DeepEqual(map[string][]any(result), tt.expected) {
				t.Errorf("JSONArrayContains() = %v, want %v", result, tt.expected)
			}

			// Test standalone function
			standaloneResult := JSONArrayContains(tt.key, tt.value)
			if !reflect.DeepEqual(map[string][]any(standaloneResult), tt.expected) {
				t.Errorf("JSONArrayContains() standalone = %v, want %v", standaloneResult, tt.expected)
			}
		})
	}
}

func TestJSONArrayContainsAny(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		values   []any
		expected map[string][]any
	}{
		{
			name:     "JSON array contains any strings",
			key:      "tags",
			values:   []any{"important", "urgent"},
			expected: map[string][]any{"tags ?| (?)": {"important", "urgent"}},
		},
		{
			name:     "JSON array contains any with camelCase key",
			key:      "userTags",
			values:   []any{"admin", "user"},
			expected: map[string][]any{"user_tags ?| (?)": {"admin", "user"}},
		},
		{
			name:     "JSON array contains any numbers",
			key:      "scores",
			values:   []any{100, 95, 90},
			expected: map[string][]any{"scores ?| (?)": {100, 95, 90}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test method on existing WhereCaluse
			clause := WhereCaluse{}
			result := clause.JSONArrayContainsAny(tt.key, tt.values...)
			if !reflect.DeepEqual(map[string][]any(result), tt.expected) {
				t.Errorf("JSONArrayContainsAny() = %v, want %v", result, tt.expected)
			}

			// Test standalone function
			standaloneResult := JSONArrayContainsAny(tt.key, tt.values...)
			if !reflect.DeepEqual(map[string][]any(standaloneResult), tt.expected) {
				t.Errorf("JSONArrayContainsAny() standalone = %v, want %v", standaloneResult, tt.expected)
			}
		})
	}
}

func TestJSONArrayContainsAll(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		values   []any
		expected map[string][]any
	}{
		{
			name:     "JSON array contains all strings",
			key:      "permissions",
			values:   []any{"read", "write"},
			expected: map[string][]any{"permissions ?& (?)": {"read", "write"}},
		},
		{
			name:     "JSON array contains all with camelCase key",
			key:      "userPermissions",
			values:   []any{"admin", "user"},
			expected: map[string][]any{"user_permissions ?& (?)": {"admin", "user"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test method on existing WhereCaluse
			clause := WhereCaluse{}
			result := clause.JSONArrayContainsAll(tt.key, tt.values...)
			if !reflect.DeepEqual(map[string][]any(result), tt.expected) {
				t.Errorf("JSONArrayContainsAll() = %v, want %v", result, tt.expected)
			}

			// Test standalone function
			standaloneResult := JSONArrayContainsAll(tt.key, tt.values...)
			if !reflect.DeepEqual(map[string][]any(standaloneResult), tt.expected) {
				t.Errorf("JSONArrayContainsAll() standalone = %v, want %v", standaloneResult, tt.expected)
			}
		})
	}
}

func TestJSONPath(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		path     string
		expected map[string][]any
	}{
		{
			name:     "JSON path exists",
			key:      "config",
			path:     "database,host",
			expected: map[string][]any{"config #> '{database,host}' IS NOT NULL": {}},
		},
		{
			name:     "JSON path exists with camelCase key",
			key:      "appConfig",
			path:     "server,port",
			expected: map[string][]any{"app_config #> '{server,port}' IS NOT NULL": {}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test method on existing WhereCaluse
			clause := WhereCaluse{}
			result := clause.JSONPath(tt.key, tt.path)
			if !reflect.DeepEqual(map[string][]any(result), tt.expected) {
				t.Errorf("JSONPath() = %v, want %v", result, tt.expected)
			}

			// Test standalone function
			standaloneResult := JSONPath(tt.key, tt.path)
			if !reflect.DeepEqual(map[string][]any(standaloneResult), tt.expected) {
				t.Errorf("JSONPath() standalone = %v, want %v", standaloneResult, tt.expected)
			}
		})
	}
}

func TestJSONPathExists(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		path     string
		expected map[string][]any
	}{
		{
			name:     "JSON path exists with #? operator",
			key:      "data",
			path:     "$.user.profile",
			expected: map[string][]any{"data #? '$.user.profile'": {}},
		},
		{
			name:     "JSON path exists with camelCase key",
			key:      "userData",
			path:     "$.profile.settings",
			expected: map[string][]any{"user_data #? '$.profile.settings'": {}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test method on existing WhereCaluse
			clause := WhereCaluse{}
			result := clause.JSONPathExists(tt.key, tt.path)
			if !reflect.DeepEqual(map[string][]any(result), tt.expected) {
				t.Errorf("JSONPathExists() = %v, want %v", result, tt.expected)
			}

			// Test standalone function
			standaloneResult := JSONPathExists(tt.key, tt.path)
			if !reflect.DeepEqual(map[string][]any(standaloneResult), tt.expected) {
				t.Errorf("JSONPathExists() standalone = %v, want %v", standaloneResult, tt.expected)
			}
		})
	}
}

func TestJSONLength(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		length   int
		expected map[string][]any
	}{
		{
			name:     "JSON length check",
			key:      "items",
			length:   5,
			expected: map[string][]any{"json_array_length(items) = ?": {5}},
		},
		{
			name:     "JSON length check with camelCase key",
			key:      "userItems",
			length:   10,
			expected: map[string][]any{"json_array_length(user_items) = ?": {10}},
		},
		{
			name:     "JSON length check zero",
			key:      "emptyArray",
			length:   0,
			expected: map[string][]any{"json_array_length(empty_array) = ?": {0}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test method on existing WhereCaluse
			clause := WhereCaluse{}
			result := clause.JSONLength(tt.key, tt.length)
			if !reflect.DeepEqual(map[string][]any(result), tt.expected) {
				t.Errorf("JSONLength() = %v, want %v", result, tt.expected)
			}

			// Test standalone function
			standaloneResult := JSONLength(tt.key, tt.length)
			if !reflect.DeepEqual(map[string][]any(standaloneResult), tt.expected) {
				t.Errorf("JSONLength() standalone = %v, want %v", standaloneResult, tt.expected)
			}
		})
	}
}

func TestJSONType(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		jsonType string
		expected map[string][]any
	}{
		{
			name:     "JSON type string",
			key:      "value",
			jsonType: "string",
			expected: map[string][]any{"json_typeof(value) = ?": {"string"}},
		},
		{
			name:     "JSON type array",
			key:      "items",
			jsonType: "array",
			expected: map[string][]any{"json_typeof(items) = ?": {"array"}},
		},
		{
			name:     "JSON type with camelCase key",
			key:      "userData",
			jsonType: "object",
			expected: map[string][]any{"json_typeof(user_data) = ?": {"object"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test method on existing WhereCaluse
			clause := WhereCaluse{}
			result := clause.JSONType(tt.key, tt.jsonType)
			if !reflect.DeepEqual(map[string][]any(result), tt.expected) {
				t.Errorf("JSONType() = %v, want %v", result, tt.expected)
			}

			// Test standalone function
			standaloneResult := JSONType(tt.key, tt.jsonType)
			if !reflect.DeepEqual(map[string][]any(standaloneResult), tt.expected) {
				t.Errorf("JSONType() standalone = %v, want %v", standaloneResult, tt.expected)
			}
		})
	}
}

func TestJSONValid(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected map[string][]any
	}{
		{
			name:     "JSON valid check",
			key:      "data",
			expected: map[string][]any{"JSON_VALID(data) = ?": {true}},
		},
		{
			name:     "JSON valid check with camelCase key",
			key:      "jsonData",
			expected: map[string][]any{"JSON_VALID(json_data) = ?": {true}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test method on existing WhereCaluse
			clause := WhereCaluse{}
			result := clause.JSONValid(tt.key)
			if !reflect.DeepEqual(map[string][]any(result), tt.expected) {
				t.Errorf("JSONValid() = %v, want %v", result, tt.expected)
			}

			// Test standalone function
			standaloneResult := JSONValid(tt.key)
			if !reflect.DeepEqual(map[string][]any(standaloneResult), tt.expected) {
				t.Errorf("JSONValid() standalone = %v, want %v", standaloneResult, tt.expected)
			}
		})
	}
}

func TestJSONSearch(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    any
		expected map[string][]any
	}{
		{
			name:     "JSON search string",
			key:      "data",
			value:    "search_term",
			expected: map[string][]any{"JSON_SEARCH(data, 'one', ?) IS NOT NULL": {"search_term"}},
		},
		{
			name:     "JSON search with camelCase key",
			key:      "searchData",
			value:    "test",
			expected: map[string][]any{"JSON_SEARCH(search_data, 'one', ?) IS NOT NULL": {"test"}},
		},
		{
			name:     "JSON search number",
			key:      "numbers",
			value:    42,
			expected: map[string][]any{"JSON_SEARCH(numbers, 'one', ?) IS NOT NULL": {42}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test method on existing WhereCaluse
			clause := WhereCaluse{}
			result := clause.JSONSearch(tt.key, tt.value)
			if !reflect.DeepEqual(map[string][]any(result), tt.expected) {
				t.Errorf("JSONSearch() = %v, want %v", result, tt.expected)
			}

			// Test standalone function
			standaloneResult := JSONSearch(tt.key, tt.value)
			if !reflect.DeepEqual(map[string][]any(standaloneResult), tt.expected) {
				t.Errorf("JSONSearch() standalone = %v, want %v", standaloneResult, tt.expected)
			}
		})
	}
}

func TestJSONAnd(t *testing.T) {
	tests := []struct {
		name     string
		clauses  []WhereCaluse
		expected int // expected number of keys in result
	}{
		{
			name: "JSON AND with multiple clauses",
			clauses: []WhereCaluse{
				JSONContains("data", "user.name", "John"),
				JSONType("config", "object"),
			},
			expected: 1, // Should create one combined key
		},
		{
			name: "JSON AND with single clause",
			clauses: []WhereCaluse{
				JSONValid("data"),
			},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := JSONAnd(tt.clauses...)
			if len(result) != tt.expected {
				t.Errorf("JSONAnd() result length = %v, want %v", len(result), tt.expected)
			}
		})
	}
}

func TestJSONOr(t *testing.T) {
	tests := []struct {
		name     string
		clauses  []WhereCaluse
		expected int // expected number of keys in result
	}{
		{
			name: "JSON OR with multiple clauses",
			clauses: []WhereCaluse{
				JSONContains("data", "user.name", "John"),
				JSONContains("data", "user.name", "Jane"),
			},
			expected: 1, // Should create one combined key
		},
		{
			name: "JSON OR with single clause",
			clauses: []WhereCaluse{
				JSONValid("data"),
			},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := JSONOr(tt.clauses...)
			if len(result) != tt.expected {
				t.Errorf("JSONOr() result length = %v, want %v", len(result), tt.expected)
			}
		})
	}
}

func TestJSONClauseChaining(t *testing.T) {
	t.Run("chain multiple JSON operations", func(t *testing.T) {
		// Test chaining multiple JSON operations
		clause := WhereCaluse{}
		result := clause.JSONContains("data", "user.name", "John").JSONType("config", "object")

		// Should have 2 keys in the map
		if len(result) != 2 {
			t.Errorf("Chained JSON operations result length = %v, want 2", len(result))
		}

		// Check if both operations are present
		expectedKeys := []string{
			"data->'user.name' @> ?",
			"json_typeof(config) = ?",
		}
		for _, key := range expectedKeys {
			if _, exists := result[key]; !exists {
				t.Errorf("Expected key %s not found in result", key)
			}
		}
	})

	t.Run("chain JSON with regular operations", func(t *testing.T) {
		// Test chaining JSON operations with regular where clauses
		clause := WhereCaluse{}
		result := clause.Eq("status", "active").JSONContains("metadata", "type", "user")

		// Should have 2 keys in the map
		if len(result) != 2 {
			t.Errorf("Mixed operations result length = %v, want 2", len(result))
		}

		// Check if both operations are present
		expectedKeys := []string{
			"status = ?",
			"metadata->'type' @> ?",
		}
		for _, key := range expectedKeys {
			if _, exists := result[key]; !exists {
				t.Errorf("Expected key %s not found in result", key)
			}
		}
	})
}