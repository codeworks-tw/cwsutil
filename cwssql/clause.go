package cwssql

import (
	"regexp"
	"strings"
)

type WhereCaluse map[string][]any

var re *regexp.Regexp = regexp.MustCompile("(.)([A-Z])")

// var s string = "HelloWorldMyNameIsCarl".replaceAll("(.)(\\p{Lu})", "$1_$2")

func (w WhereCaluse) Eq(key string, value any) WhereCaluse {
	w[strings.ToLower(re.ReplaceAllString(key, "$1_$2"))+" = ?"] = []any{value}
	return w
}

func (w WhereCaluse) Ne(key string, value any) WhereCaluse {
	w[key+" != ?"] = []any{value}
	return w
}

func (w WhereCaluse) Gt(key string, value any) WhereCaluse {
	w[key+" > ?"] = []any{value}
	return w
}

func (w WhereCaluse) Gte(key string, value any) WhereCaluse {
	w[key+" >= ?"] = []any{value}
	return w
}

func (w WhereCaluse) Lt(key string, value any) WhereCaluse {
	w[key+" < ?"] = []any{value}
	return w
}

func (w WhereCaluse) Lte(key string, value any) WhereCaluse {
	w[key+" <= ?"] = []any{value}
	return w
}

func (w WhereCaluse) In(key string, values ...any) WhereCaluse {
	w[key+" IN (?)"] = values
	return w
}

func (w WhereCaluse) Nin(key string, values ...any) WhereCaluse {
	w[key+" NOT IN (?)"] = values
	return w
}

func (w WhereCaluse) Like(key string, value any) WhereCaluse {
	w[key+" LIKE ?"] = []any{value}
	return w
}

func (w WhereCaluse) Between(key string, left any, right any) WhereCaluse {
	w[key+" BETWEEN ? AND ?"] = []any{left, right}
	return w
}
func (w WhereCaluse) And(clauses ...WhereCaluse) WhereCaluse {
	query := ""
	vs := []any{}
	for i, wc := range clauses {
		if i > 0 {
			query += " AND "
		}
		query += "("
		wi := 0
		for k, v := range wc {
			query += k
			if len(v) > 0 {
				vs = append(vs, v...)
			}
			if wi < len(wc)-1 {
				query += " AND "
			}
			wi++
		}
		query += ")"
	}
	w[query] = vs
	return w
}

func (w WhereCaluse) Or(clauses ...WhereCaluse) WhereCaluse {
	query := ""
	vs := []any{}
	for i, wc := range clauses {
		if i > 0 {
			query += " OR "
		}
		query += "("
		wi := 0
		for k, v := range wc {
			query += k
			if len(v) > 0 {
				vs = append(vs, v...)
			}
			if wi < len(wc)-1 {
				query += " AND "
			}
			wi++
		}
		query += ")"
	}
	w[query] = vs
	return w
}

func Eq(key string, value any) WhereCaluse {
	return WhereCaluse{}.Eq(key, value)
}

func Ne(key string, value any) WhereCaluse {
	return WhereCaluse{}.Ne(key, value)
}

func Gt(key string, value any) WhereCaluse {
	return WhereCaluse{}.Gt(key, value)
}

func Gte(key string, value any) WhereCaluse {
	return WhereCaluse{}.Gte(key, value)
}

func Lt(key string, value any) WhereCaluse {
	return WhereCaluse{}.Lt(key, value)
}

func Lte(key string, value any) WhereCaluse {
	return WhereCaluse{}.Lte(key, value)
}

func In(key string, values ...any) WhereCaluse {
	return WhereCaluse{}.In(key, values...)
}

func Nin(key string, values ...any) WhereCaluse {
	return WhereCaluse{}.Nin(key, values...)
}

func Like(key string, value any) WhereCaluse {
	return WhereCaluse{}.Like(key, value)
}

func Between(key string, left any, right any) WhereCaluse {
	return WhereCaluse{}.Between(key, left, right)
}

func And(clauses ...WhereCaluse) WhereCaluse {
	return WhereCaluse{}.And(clauses...)
}

func Or(clauses ...WhereCaluse) WhereCaluse {
	return WhereCaluse{}.Or(clauses...)
}
