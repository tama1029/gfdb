package render

import (
	"fmt"
	"github.com/gertd/go-pluralize"
	"go/format"
)

type GenGoaResultTypeRender struct{}

func NewGenGoaResultTypeRender() *GenGoaResultTypeRender {
	return &GenGoaResultTypeRender{}
}

func (r *GenGoaResultTypeRender) RenderFacade(tableDataTypes map[string]map[string]map[string]string, tableNamesSorted map[string][]string, tableAndComment map[string]string) (map[string][]byte, error) {
	renders := map[string][]byte{}
	plu := pluralize.NewClient()
	for tableName, columnsSorted := range tableNamesSorted {
		columnTypes := tableDataTypes[tableName]

		single := plu.Singular(tableName)
		singleEntityName := fmtFieldName(stringifyFirstChar(single))
		multiEntityName := fmtFieldName(stringifyFirstChar(tableName))
		pR := r.packageRender("design")
		sR := r.structRender(columnTypes, columnsSorted, singleEntityName, tableName, tableAndComment[tableName], single)
		ssR := r.listRender(columnTypes, columnsSorted, singleEntityName, tableName, tableAndComment[tableName], single, multiEntityName)
		src := pR + sR + ssR
		formatted, err := format.Source([]byte(src))
		if err != nil {
			return nil, fmt.Errorf("error formatting: %s, was formatting\n%s", err, src)
		}
		renders[single] = formatted
	}
	return renders, nil
}

func (_ *GenGoaResultTypeRender) packageRender(packageName string) string {
	return fmt.Sprintf("package %s\n", packageName)
}

func (r *GenGoaResultTypeRender) structRender(
	columnTypes map[string]map[string]string,
	columnsSorted []string,
	pluralEntityName string,
	tableName string,
	tableComment string,
	singleTableName string) string {
	identifier := fmt.Sprintf("application/vnd.%s+json", singleTableName)
	dbTypes := r.generateMysqlTypes(columnTypes, columnsSorted, 0, false, false, false, tableComment, pluralEntityName)
	return fmt.Sprintf("var %sResult = ResultType(\"%s\", func() { %s })\n",
		pluralEntityName,
		identifier,
		dbTypes)
}

func (r *GenGoaResultTypeRender) listRender(
	columnTypes map[string]map[string]string,
	columnsSorted []string,
	pluralEntityName string,
	tableName string,
	tableComment string,
	singleTableName string,
	multiEntityName string) string {
	header := fmt.Sprintf("\nvar List%sResult = ResultType(\"application/vnd.list_%s+json\", func() {\n", multiEntityName, tableName)
	footer := "})\n"

	body := ""
	body += fmt.Sprintf("Attribute(\"%s\", ArrayOf(%sResult))\n", tableName, pluralEntityName)
	body += "Attribute(\"total\", Int64)\n"
	body += fmt.Sprintf("Required(\"%s\", \"total\")\n", tableName)
	body += "View(\"default\", func() {\n"
	body += fmt.Sprintf("Attribute(\"%s\")\n", tableName)
	body += fmt.Sprintf("Attribute(\"total\")\n")
	body += "})\n"

	return header + body + footer
}

//var ListThemesResult = ResultType("application/vnd.list_themes+json", func() {
//	Attribute("themes", ArrayOf(ThemeResult))
//	Attribute("total", Int64)
//	Required("themes", "total")
//	View("default", func() {
//		Attribute("themes")
//		Attribute("total")
//	})
//})
func (_ *GenGoaResultTypeRender) generateMysqlTypes(obj map[string]map[string]string, columnsSorted []string, depth int, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool, tableComment string, pluralEntityName string) string {
	//dummy := ""
	header := ""
	header += fmt.Sprintf("Description(\"%s\")\n", tableComment)
	header += fmt.Sprintf("CreateFrom(result.%s{})\n", pluralEntityName)
	header += fmt.Sprintf("Extend(%s)\n\n", pluralEntityName)
	structure := ""
	structure += "View(\"default\", func() {"
	//var i interface{} = "\"aiaiiaiai\""
	for _, key := range columnsSorted {
		//fieldName := fmtFieldName(stringifyFirstChar(key))
		//comment := mysqlType["comment"]

		//RequiredAttribute("id", UInt64, "id", func() {
		//	Example(123)
		//})

		structure += fmt.Sprintf("Attribute(\"%s\")\n", key)
	}
	structure += "})"
	return header + structure
}
