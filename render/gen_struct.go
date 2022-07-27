package render

import (
	"fmt"
	"github.com/gertd/go-pluralize"
	"go/format"
	"strings"
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
		sR := r.structRender(columnTypes, columnsSorted, singleEntityName)
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

func (r *GenStructRender) structRender(columnTypes map[string]map[string]string, columnsSorted []string, pluralEntityName string) string {
	dbTypes := r.generateMysqlTypes(columnTypes, columnsSorted, 0, false, true, false)
	return fmt.Sprintf("type %s %s\n}\n",
		pluralEntityName,
		dbTypes)
}

func (r *GenStructRender) dtoRender(entityName string, columnsSorted []string) string {
	dtoBody := fmt.Sprintf("r := &result.%s{}", entityName)

	for _, key := range columnsSorted {
		if key == "deleted_at" {
			continue
		}
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
	dtosBody += "for _, e := range es {\n"
	dtosBody += fmt.Sprintf("res = append(res, ToResult%s(e))\n", singleEntityName)
	dtosBody += "}\n\n"
	dtosBody += "return res\n"
	return fmt.Sprintf("func ToResult%s(es []*%s) []*result.%s {\n%s}\n\n", pluralEntityName, singleEntityName, singleEntityName, dtosBody)
}

func (_ *GenStructRender) generateMysqlTypes(obj map[string]map[string]string, columnsSorted []string, depth int, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool) string {
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

		valueType = mysqlTypeToGoType(key, mysqlType["value"], mysqlType["columnType"], nullable, gureguTypes, false)

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
