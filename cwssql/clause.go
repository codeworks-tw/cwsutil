package cwssql

import "github.com/codeworks-tw/cwsutil/cwsbase"

type WhereCaluse map[string][]any

func (w WhereCaluse) Eq(key string, value any) WhereCaluse {
	w[cwsbase.ToSnakeCase(key)+" = ?"] = []any{value}
	return w
}

func (w WhereCaluse) Ne(key string, value any) WhereCaluse {
	w[cwsbase.ToSnakeCase(key)+" != ?"] = []any{value}
	return w
}

func (w WhereCaluse) Gt(key string, value any) WhereCaluse {
	w[cwsbase.ToSnakeCase(key)+" > ?"] = []any{value}
	return w
}

func (w WhereCaluse) Gte(key string, value any) WhereCaluse {
	w[cwsbase.ToSnakeCase(key)+" >= ?"] = []any{value}
	return w
}

func (w WhereCaluse) Lt(key string, value any) WhereCaluse {
	w[cwsbase.ToSnakeCase(key)+" < ?"] = []any{value}
	return w
}

func (w WhereCaluse) Lte(key string, value any) WhereCaluse {
	w[cwsbase.ToSnakeCase(key)+" <= ?"] = []any{value}
	return w
}

func (w WhereCaluse) In(key string, values ...any) WhereCaluse {
	if len(values) > 0 {
		w[cwsbase.ToSnakeCase(key)+" IN (?)"] = values
	}
	return w
}

func (w WhereCaluse) Nin(key string, values ...any) WhereCaluse {
	if len(values) > 0 {
		w[cwsbase.ToSnakeCase(key)+" NOT IN (?)"] = values
	}
	return w
}

func (w WhereCaluse) Like(key string, value any) WhereCaluse {
	w[cwsbase.ToSnakeCase(key)+" LIKE ?"] = []any{value}
	return w
}

func (w WhereCaluse) Between(key string, left any, right any) WhereCaluse {
	w[cwsbase.ToSnakeCase(key)+" BETWEEN ? AND ?"] = []any{left, right}
	return w
}

func (w WhereCaluse) IsNull(key string) WhereCaluse {
	w[cwsbase.ToSnakeCase(key)+" IS NULL"] = []any{}
	return w
}

func (w WhereCaluse) IsNotNull(key string) WhereCaluse {
	w[cwsbase.ToSnakeCase(key)+" IS NOT NULL"] = []any{}
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
				query += " OR "
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

func IsNull(key string) WhereCaluse {
	return WhereCaluse{}.IsNull(key)
}

func IsNotNull(key string) WhereCaluse {
	return WhereCaluse{}.IsNotNull(key)
}

func And(clauses ...WhereCaluse) WhereCaluse {
	return WhereCaluse{}.And(clauses...)
}

func Or(clauses ...WhereCaluse) WhereCaluse {
	return WhereCaluse{}.Or(clauses...)
}
