package render

import (
	"fmt"
	"github.com/gertd/go-pluralize"
	"go/format"
	"strings"
)

type GenResultRender struct{}

func NewGenResultRender() *GenResultRender {
	return &GenResultRender{}
}

func (r *GenResultRender) RenderFacade(tableDataTypes map[string]map[string]map[string]string, tableNamesSorted map[string][]string) (map[string][]byte, error) {
	renders := map[string][]byte{}
	plu := pluralize.NewClient()
	for tableName, columnsSorted := range tableNamesSorted {
		columnTypes := tableDataTypes[tableName]

		single := plu.Singular(tableName)
		singleEntityName := fmtFieldName(stringifyFirstChar(single))
		pR := r.packageRender("result")
		sR := r.structRender(columnTypes, columnsSorted, singleEntityName)
		src := pR + sR
		formatted, err := format.Source([]byte(src))
		if err != nil {
			return nil, fmt.Errorf("error formatting: %s, was formatting\n%s", err, src)
		}
		renders[single] = formatted
	}
	return renders, nil
}

func (_ *GenResultRender) packageRender(packageName string) string {
	return fmt.Sprintf("package %s\n", packageName)
}

func (r *GenResultRender) structRender(columnTypes map[string]map[string]string, columnsSorted []string, pluralEntityName string) string {
	dbTypes := r.generateMysqlTypes(columnTypes, columnsSorted, 0, false, false, false)
	return fmt.Sprintf("type %s %s\n}\n",
		pluralEntityName,
		dbTypes)
}

func (_ *GenResultRender) generateMysqlTypes(obj map[string]map[string]string, columnsSorted []string, depth int, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool) string {
	structure := "struct {"

	for _, key := range columnsSorted {
		mysqlType := obj[key]
		nullable := false
		//if mysqlType["nullable"] == "YES" {
		//	nullable = true
		//}

		primary := ""
		if mysqlType["primary"] == "PRI" {
			primary = ";primary_key"
		}

		// Get the corresponding go value type for this mysql type
		var valueType string
		// If the guregu (https://github.com/guregu/null) CLI option is passed use its types, otherwise use go's sql.NullX

		valueType = mysqlTypeToGoType(key, mysqlType["value"], mysqlType["columnType"], nullable, gureguTypes, true)

		fieldName := fmtFieldName(stringifyFirstChar(key))
		var annotations []string
		if gormAnnotation == true {
			annotations = append(annotations, fmt.Sprintf("gorm:\"column:%s%s\"", key, primary))
		}
		if jsonAnnotation == true {
			annotations = append(annotations, fmt.Sprintf("json:\"%s\"", key))
		}
		comment := mysqlType["comment"]
		if len(annotations) > 0 {
			// add colulmn comment
			structure += fmt.Sprintf("\n%s %s `%s`  // %s", fieldName, valueType, strings.Join(annotations, " "), comment)
			//structure += fmt.Sprintf("\n%s %s `%s`", fieldName, valueType, strings.Join(annotations, " "))
		} else {
			structure += fmt.Sprintf("\n%s %s  // %s", fieldName, valueType, comment)
		}
	}
	return structure
}
