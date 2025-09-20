package cwssql

import "github.com/codeworks-tw/cwsutil/cwsbase"

// JSONContains checks if JSON column contains a value at the given path
// For PostgreSQL: column->'path' @> value
// For MySQL: JSON_CONTAINS(column, value, 'path')
func (j WhereCaluse) JSONContains(key string, path string, value any) WhereCaluse {
	j[cwsbase.ToSnakeCase(key)+"->'"+path+"' @> ?"] = []any{value}
	return j
}

// JSONBContains checks if JSONB column contains a value (PostgreSQL specific)
func (j WhereCaluse) JSONBContains(key string, value any) WhereCaluse {
	j[cwsbase.ToSnakeCase(key)+" @> ?"] = []any{value}
	return j
}

// JSONBContainedBy checks if JSONB column is contained by a value (PostgreSQL specific)
func (j WhereCaluse) JSONBContainedBy(key string, value any) WhereCaluse {
	j[cwsbase.ToSnakeCase(key)+" <@ ?"] = []any{value}
	return j
}

// JSONExtract extracts value from JSON at path and compares with value
// For PostgreSQL: column->'path' = value
// For MySQL: JSON_EXTRACT(column, 'path') = value
func (j WhereCaluse) JSONExtract(key string, path string, value any) WhereCaluse {
	j[cwsbase.ToSnakeCase(key)+"->'"+path+"' = ?"] = []any{value}
	return j
}

// JSONExtractText extracts text value from JSON at path (PostgreSQL specific)
// column->>'path' = value
func (j WhereCaluse) JSONExtractText(key string, path string, value any) WhereCaluse {
	j[cwsbase.ToSnakeCase(key)+"->>'"+path+"' = ?"] = []any{value}
	return j
}

// JSONArrayContains checks if JSON array contains a specific value
// For PostgreSQL: column ? 'value'
func (j WhereCaluse) JSONArrayContains(key string, value any) WhereCaluse {
	j[cwsbase.ToSnakeCase(key)+" ? ?"] = []any{value}
	return j
}

// JSONArrayContainsAny checks if JSON array contains any of the specified values
// For PostgreSQL: column ?| array['value1', 'value2']
func (j WhereCaluse) JSONArrayContainsAny(key string, values ...any) WhereCaluse {
	j[cwsbase.ToSnakeCase(key)+" ?| (?)"] = values
	return j
}

// JSONArrayContainsAll checks if JSON array contains all of the specified values
// For PostgreSQL: column ?& array['value1', 'value2']
func (j WhereCaluse) JSONArrayContainsAll(key string, values ...any) WhereCaluse {
	j[cwsbase.ToSnakeCase(key)+" ?& (?)"] = values
	return j
}

// JSONPath checks if JSON path exists
// For PostgreSQL: column #> '{path,subpath}' IS NOT NULL
func (j WhereCaluse) JSONPath(key string, path string) WhereCaluse {
	j[cwsbase.ToSnakeCase(key)+" #> '{"+path+"}' IS NOT NULL"] = []any{}
	return j
}

// JSONPathExists checks if JSON path exists using the #? operator (PostgreSQL specific)
func (j WhereCaluse) JSONPathExists(key string, path string) WhereCaluse {
	j[cwsbase.ToSnakeCase(key)+" #? '"+path+"'"] = []any{}
	return j
}

// JSONLength checks the length of a JSON array or object
// For PostgreSQL: json_array_length(column) = length
// For MySQL: JSON_LENGTH(column) = length
func (j WhereCaluse) JSONLength(key string, length int) WhereCaluse {
	j["json_array_length("+cwsbase.ToSnakeCase(key)+") = ?"] = []any{length}
	return j
}

// JSONType checks the type of a JSON value
// For PostgreSQL: json_typeof(column) = 'type'
// For MySQL: JSON_TYPE(column) = 'type'
func (j WhereCaluse) JSONType(key string, jsonType string) WhereCaluse {
	j["json_typeof("+cwsbase.ToSnakeCase(key)+") = ?"] = []any{jsonType}
	return j
}

// JSONValid checks if the JSON is valid
// For MySQL: JSON_VALID(column) = true
func (j WhereCaluse) JSONValid(key string) WhereCaluse {
	j["JSON_VALID("+cwsbase.ToSnakeCase(key)+") = ?"] = []any{true}
	return j
}

// JSONSearch searches for a value in JSON (MySQL specific)
// JSON_SEARCH(column, 'one', 'value') IS NOT NULL
func (j WhereCaluse) JSONSearch(key string, value any) WhereCaluse {
	j["JSON_SEARCH("+cwsbase.ToSnakeCase(key)+", 'one', ?) IS NOT NULL"] = []any{value}
	return j
}

// Standalone functions following the same pattern as clause.go

// JSONContains creates a new JSON contains clause
func JSONContains(key string, path string, value any) WhereCaluse {
	return WhereCaluse{}.JSONContains(key, path, value)
}

// JSONBContains creates a new JSONB contains clause
func JSONBContains(key string, value any) WhereCaluse {
	return WhereCaluse{}.JSONBContains(key, value)
}

// JSONBContainedBy creates a new JSONB contained by clause
func JSONBContainedBy(key string, value any) WhereCaluse {
	return WhereCaluse{}.JSONBContainedBy(key, value)
}

// JSONExtract creates a new JSON extract clause
func JSONExtract(key string, path string, value any) WhereCaluse {
	return WhereCaluse{}.JSONExtract(key, path, value)
}

// JSONExtractText creates a new JSON extract text clause
func JSONExtractText(key string, path string, value any) WhereCaluse {
	return WhereCaluse{}.JSONExtractText(key, path, value)
}

// JSONArrayContains creates a new JSON array contains clause
func JSONArrayContains(key string, value any) WhereCaluse {
	return WhereCaluse{}.JSONArrayContains(key, value)
}

// JSONArrayContainsAny creates a new JSON array contains any clause
func JSONArrayContainsAny(key string, values ...any) WhereCaluse {
	return WhereCaluse{}.JSONArrayContainsAny(key, values...)
}

// JSONArrayContainsAll creates a new JSON array contains all clause
func JSONArrayContainsAll(key string, values ...any) WhereCaluse {
	return WhereCaluse{}.JSONArrayContainsAll(key, values...)
}

// JSONPath creates a new JSON path clause
func JSONPath(key string, path string) WhereCaluse {
	return WhereCaluse{}.JSONPath(key, path)
}

// JSONPathExists creates a new JSON path exists clause
func JSONPathExists(key string, path string) WhereCaluse {
	return WhereCaluse{}.JSONPathExists(key, path)
}

// JSONLength creates a new JSON length clause
func JSONLength(key string, length int) WhereCaluse {
	return WhereCaluse{}.JSONLength(key, length)
}

// JSONType creates a new JSON type clause
func JSONType(key string, jsonType string) WhereCaluse {
	return WhereCaluse{}.JSONType(key, jsonType)
}

// JSONValid creates a new JSON valid clause
func JSONValid(key string) WhereCaluse {
	return WhereCaluse{}.JSONValid(key)
}

// JSONSearch creates a new JSON search clause
func JSONSearch(key string, value any) WhereCaluse {
	return WhereCaluse{}.JSONSearch(key, value)
}

// JSONAnd creates a new AND combination of JSON clauses
func JSONAnd(clauses ...WhereCaluse) WhereCaluse {
	return WhereCaluse{}.And(clauses...)
}

// JSONOr creates a new OR combination of JSON clauses
func JSONOr(clauses ...WhereCaluse) WhereCaluse {
	return WhereCaluse{}.Or(clauses...)
}
