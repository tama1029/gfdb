package render

import (
	"strconv"
	"strings"
	"unicode"
)

// mysqlTypeToGoType converts the mysql types to go compatible sql.Nullable (https://golang.org/pkg/database/sql/) types
func mysqlTypeToGoType(key string, mysqlType string, columnType string, nullable bool, gureguTypes bool, isResult bool) string {
	switch mysqlType {
	case "tinyint", "int", "smallint", "mediumint":
		if nullable {
			if gureguTypes {
				return gureguNullInt
			}
			return sqlNullInt
		}
		if strings.Contains(columnType, "tinyint(1)") {
			return golangBool
		}
		return golangInt
	case "bigint":
		if nullable {
			if gureguTypes {
				return gureguNullInt
			}
			return sqlNullInt
		}
		if strings.Contains(columnType, "unsigned") {
			return golangUInt64
		}
		return golangInt64
	case "char", "enum", "varchar", "longtext", "mediumtext", "text", "tinytext", "json":
		if nullable {
			if gureguTypes {
				return gureguNullString
			}
			return sqlNullString
		}
		return "string"
	case "date", "datetime", "time", "timestamp":
		if nullable && gureguTypes {
			return gureguNullTime
		}
		if key == "created_at" && isResult {
			return golangString
		}
		if key == "updated_at" && isResult {
			return golangString
		}
		if key == "deleted_at" {
			if isResult {
				return resultGormDeletedAt
			}
			return gormDeletedAt
		}
		return golangTime
	case "decimal", "double":
		if nullable {
			if gureguTypes {
				return gureguNullFloat
			}
			return sqlNullFloat
		}
		return golangFloat64
	case "float":
		if nullable {
			if gureguTypes {
				return gureguNullFloat
			}
			return sqlNullFloat
		}
		return golangFloat32
	case "binary", "blob", "longblob", "mediumblob", "varbinary":
		return golangByteArray
	}
	return ""
}

func mysqlTypeToGoaType(key string, mysqlType string, columnType string, nullable bool, gureguTypes bool, isResult bool) string {
	switch mysqlType {
	case "tinyint", "int", "smallint", "mediumint":
		if strings.Contains(columnType, "tinyint(1)") {
			return goaBool
		}
		return goaInt
	case "bigint":
		if strings.Contains(columnType, "unsigned") {
			return goaUInt64
		}
		return goaInt64
	case "char", "enum", "varchar", "longtext", "mediumtext", "text", "tinytext", "json":
		return goaString
	case "date", "datetime", "time", "timestamp":
		if key == "created_at" && isResult {
			return goaString
		}
		if key == "updated_at" && isResult {
			return goaString
		}
		if key == "deleted_at" {
			return goaAny
		}
		return goaString
	case "decimal", "double":
		return goaFloat64
	case "float":
		return goaFloat32
	case "binary", "blob", "longblob", "mediumblob", "varbinary":
		return goaBytes
	}
	return ""
}

func mysqlTypeToGoaExample(key string, mysqlType string, columnType string, nullable bool, gureguTypes bool, isResult bool) interface{} {
	//loc, _ := time.LoadLocation("Asia/Tokyo")
	//Example(time.Date(2019, 3, 1, 0, 0, 0, 0, loc).Format(time.RFC3339))
	//Format(FormatDateTime)
	switch mysqlType {
	case "tinyint", "int", "smallint", "mediumint":
		if strings.Contains(columnType, "tinyint(1)") {
			return "Example(true)"
		}
		return "Example(1)"
	case "bigint":
		return "Example(1)"
	case "char", "enum", "varchar", "longtext", "mediumtext", "text", "tinytext", "json":
		return "Example(\"\")"
	case "date", "datetime", "time", "timestamp":
		if key == "created_at" {
			return `loc, _ := time.LoadLocation("Asia/Tokyo")
Example(time.Date(2019, 3, 1, 0, 0, 0, 0, loc).Format(time.RFC3339))
Format(FormatDateTime)
`
		}
		if key == "updated_at" {
			return `loc, _ := time.LoadLocation("Asia/Tokyo")
Example(time.Date(2019, 3, 1, 0, 0, 0, 0, loc).Format(time.RFC3339))
Format(FormatDateTime)
`
		}
		if key == "deleted_at" {
			return `loc, _ := time.LoadLocation("Asia/Tokyo")
Example(time.Date(2019, 3, 1, 0, 0, 0, 0, loc).Format(time.RFC3339))
`
		}
		return goaString
	case "decimal", "double":
		return "Example(1.1)"
	case "float":
		return "Example(1.1)"
	case "binary", "blob", "longblob", "mediumblob", "varbinary":
		return "Example(\"\")"
	}
	return ""
}

func fmtFieldName(s string) string {
	name := lintFieldName(s)
	runes := []rune(name)
	for i, c := range runes {
		ok := unicode.IsLetter(c) || unicode.IsDigit(c)
		if i == 0 {
			ok = unicode.IsLetter(c)
		}
		if !ok {
			runes[i] = '_'
		}
	}
	return string(runes)
}

func lintFieldName(name string) string {
	// Fast path for simple cases: "_" and all lowercase.
	if name == "_" {
		return name
	}

	for len(name) > 0 && name[0] == '_' {
		name = name[1:]
	}

	allLower := true
	for _, r := range name {
		if !unicode.IsLower(r) {
			allLower = false
			break
		}
	}
	if allLower {
		runes := []rune(name)
		if u := strings.ToUpper(name); commonInitialisms[u] {
			copy(runes[0:], []rune(u))
		} else {
			runes[0] = unicode.ToUpper(runes[0])
		}
		return string(runes)
	}

	// Split camelCase at any lower->upper transition, and split on underscores.
	// Check each word for common initialisms.
	runes := []rune(name)
	w, i := 0, 0 // index of start of word, scan
	for i+1 <= len(runes) {
		eow := false // whether we hit the end of a word

		if i+1 == len(runes) {
			eow = true
		} else if runes[i+1] == '_' {
			// underscore; shift the remainder forward over any run of underscores
			eow = true
			n := 1
			for i+n+1 < len(runes) && runes[i+n+1] == '_' {
				n++
			}

			// Leave at most one underscore if the underscore is between two digits
			if i+n+1 < len(runes) && unicode.IsDigit(runes[i]) && unicode.IsDigit(runes[i+n+1]) {
				n--
			}

			copy(runes[i+1:], runes[i+n+1:])
			runes = runes[:len(runes)-n]
		} else if unicode.IsLower(runes[i]) && !unicode.IsLower(runes[i+1]) {
			// lower->non-lower
			eow = true
		}
		i++
		if !eow {
			continue
		}

		// [w,i) is a word.
		word := string(runes[w:i])
		if u := strings.ToUpper(word); commonInitialisms[u] {
			// All the common initialisms are ASCII,
			// so we can replace the bytes exactly.
			copy(runes[w:], []rune(u))

		} else if strings.ToLower(word) == word {
			// already all lowercase, and not the first word, so uppercase the first character.
			runes[w] = unicode.ToUpper(runes[w])
		}
		w = i
	}
	return string(runes)
}

// convert first character ints to strings
func stringifyFirstChar(str string) string {
	first := str[:1]

	i, err := strconv.ParseInt(first, 10, 8)

	if err != nil {
		return str
	}

	return intToWordMap[i] + "_" + str[1:]
}

var commonInitialisms = map[string]bool{
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SSH":   true,
	"TLS":   true,
	"TTL":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
}

var intToWordMap = []string{
	"zero",
	"one",
	"two",
	"three",
	"four",
	"five",
	"six",
	"seven",
	"eight",
	"nine",
}

const (
	golangByteArray  = "[]byte"
	gureguNullInt    = "null.Int"
	sqlNullInt       = "sql.NullInt64"
	golangInt        = "int"
	golangBool       = "bool"
	golangInt64      = "int64"
	golangUInt64     = "uint64"
	golangString     = "string"
	gureguNullFloat  = "null.Float"
	sqlNullFloat     = "sql.NullFloat64"
	golangFloat      = "float"
	golangFloat32    = "float32"
	golangFloat64    = "float64"
	gureguNullString = "null.String"
	sqlNullString    = "sql.NullString"
	gureguNullTime   = "null.Time"
	golangTime       = "time.Time"
	gormDeletedAt    = "gorm.DeletedAt"
	// result
	resultGormDeletedAt = "interface{}"

	// goa
	goaBool    = "Boolean"
	goaInt     = "Int"
	goaInt32   = "Int32"
	goaInt64   = "Int64"
	goaUInt    = "UInt"
	goaUInt32  = "UInt32"
	goaUInt64  = "UInt64"
	goaFloat32 = "Float32"
	goaFloat64 = "Float64"
	goaString  = "String"
	goaBytes   = "Bytes"
	goaAny     = "Any"
)
