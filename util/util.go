package util

import "database/sql"

func ToNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func ToNullInt64(i int) sql.NullInt64 {
	return sql.NullInt64{Int64: int64(i), Valid: i != 0}
}

func TuNullStringArray(arr []string) []sql.NullString {
	var a []sql.NullString

	for _, val := range arr {
		a = append(a, sql.NullString{String: val, Valid: val != ""})
	}

	return a
}
