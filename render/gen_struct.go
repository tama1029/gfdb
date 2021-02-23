package render

import (
	"fmt"
	"github.com/gertd/go-pluralize"
	"go/format"
	"strconv"
	"strings"
	"unicode"
)

type GenStructRender struct{}

func NewGenStructRender() *GenStructRender {
	return &GenStructRender{}
}

func (r *GenStructRender) RenderFacade(tableDataTypes map[string]map[string]map[string]string, tableNamesSorted map[string][]string) (map[string][]byte, error) {
	renders := map[string][]byte{}
	plu := pluralize.NewClient()
	for tableName, columnsSorted := range tableNamesSorted {
		columnTypes := tableDataTypes[tableName]

		single := plu.Singular(tableName)
		singleEntityName := fmtFieldName(stringifyFirstChar(single))
		pluralEntityName := fmtFieldName(stringifyFirstChar(tableName))
		pR := r.packageRender("entity")
		sR := r.structRender(columnTypes, columnsSorted, pluralEntityName)
		dR := r.dtoRender(singleEntityName, columnsSorted)
		dsR := r.dtosRender(singleEntityName, pluralEntityName)
		src := pR + sR + dR + dsR
		formatted, err := format.Source([]byte(src))
		if err != nil {
			return nil, fmt.Errorf("error formatting: %s, was formatting\n%s", err, src)
		}
		renders[single] = formatted
	}
	return renders, nil
}

func (_ *GenStructRender) packageRender(packageName string) string {
	return fmt.Sprintf("package %s\n", packageName)
}

func (_ *GenStructRender) structRender(columnTypes map[string]map[string]string, columnsSorted []string, pluralEntityName string) string {
	dbTypes := generateMysqlTypes(columnTypes, columnsSorted, 0, false, true, false)
	return fmt.Sprintf("type %s %s\n}\n",
		pluralEntityName,
		dbTypes)
}

func (r *GenStructRender) dtoRender(entityName string, columnsSorted []string) string {
	dtoBody := fmt.Sprintf("r := &result.%s{}", entityName)

	for _, key := range columnsSorted {
		leftFieldName := fmtFieldName(stringifyFirstChar(key))
		rightFieldName := fmtFieldName(stringifyFirstChar(key))
		if rightFieldName == "CreatedAt" || rightFieldName == "UpdatedAt" {
			rightFieldName = rightFieldName + ".Format(time.RFC3339)"
		}
		dtoBody += fmt.Sprintf("\nr.%s=e.%s", leftFieldName, rightFieldName)
	}
	dtoBody += "\nreturn r"
	return fmt.Sprintf("func ToResult%s(e *%s) *result.%s {\n%s}\n\n", entityName, entityName, entityName, dtoBody)
}

func (_ *GenStructRender) dtosRender(singleEntityName, pluralEntityName string) string {
	dtosBody := fmt.Sprintf("var res []*result.%s\n\n", singleEntityName)
	dtosBody += "for _, e := range ee {\n"
	dtosBody += fmt.Sprintf("res = append(res, ToResult%s(&e))\n", singleEntityName)
	dtosBody += "}\n\n"
	dtosBody += "return res\n"
	return fmt.Sprintf("func ToResult%s(ee []*%s) []*result.%s {\n%s}\n\n", pluralEntityName, singleEntityName, singleEntityName, dtosBody)
}

func generateMysqlTypes(obj map[string]map[string]string, columnsSorted []string, depth int, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool) string {
	structure := "struct {"

	for _, key := range columnsSorted {
		mysqlType := obj[key]
		nullable := false
		if mysqlType["nullable"] == "YES" {
			nullable = true
		}

		primary := ""
		if mysqlType["primary"] == "PRI" {
			primary = ";primary_key"
		}

		// Get the corresponding go value type for this mysql type
		var valueType string
		// If the guregu (https://github.com/guregu/null) CLI option is passed use its types, otherwise use go's sql.NullX

		valueType = mysqlTypeToGoType(key, mysqlType["value"], mysqlType["columnType"], nullable, gureguTypes)

		fieldName := fmtFieldName(stringifyFirstChar(key))
		var annotations []string
		if gormAnnotation == true {
			annotations = append(annotations, fmt.Sprintf("gorm:\"column:%s%s\"", key, primary))
		}
		if jsonAnnotation == true {
			annotations = append(annotations, fmt.Sprintf("json:\"%s\"", key))
		}

		if len(annotations) > 0 {
			// add colulmn comment
			comment := mysqlType["comment"]
			structure += fmt.Sprintf("\n%s %s `%s`  // %s", fieldName, valueType, strings.Join(annotations, " "), comment)
			//structure += fmt.Sprintf("\n%s %s `%s`", fieldName, valueType, strings.Join(annotations, " "))
		} else {
			structure += fmt.Sprintf("\n%s %s", fieldName, valueType)
		}
	}
	return structure
}

// mysqlTypeToGoType converts the mysql types to go compatible sql.Nullable (https://golang.org/pkg/database/sql/) types
func mysqlTypeToGoType(key string, mysqlType string, colmunType string, nullable bool, gureguTypes bool) string {
	switch mysqlType {
	case "tinyint", "int", "smallint", "mediumint":
		if nullable {
			if gureguTypes {
				return gureguNullInt
			}
			return sqlNullInt
		}
		return golangInt
	case "bigint":
		if nullable {
			if gureguTypes {
				return gureguNullInt
			}
			return sqlNullInt
		}
		if strings.Contains(colmunType, "unsigned") {
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
		if key == "deleted_at" {
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
	golangInt64      = "int64"
	golangUInt64     = "uint64"
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
)
