package db

import (
	"fmt"
	"github.com/tama1029/gfdb/client"
)

type TableInfo struct {
	c *client.Client
}

func NewTableInfo(c *client.Client) *TableInfo {
	return &TableInfo{c: c}
}

func (co *TableInfo) GetTableInfo(databaseName string) (map[string]string, error) {
	tableInfo := map[string]string{}

	tableDataTypeQuery := "SELECT TABLE_NAME, TABLE_COMMENT FROM INFORMATION_SCHEMA.TABLES where TABLE_SCHEMA = ? order by TABLE_NAME"

	rows, err := co.c.Db.Query(tableDataTypeQuery, databaseName)
	if err != nil {
		return nil, err
	}
	if rows != nil {
		defer rows.Close()
	} else {
		return nil, err
	}

	for rows.Next() {
		var table string
		var comment string
		rows.Scan(&table, &comment)
		tableInfo[table] = comment
	}
	return tableInfo, nil
}

func (co *TableInfo) GetColumnInfo(databaseName string) (map[string]map[string]map[string]string, map[string][]string, error) {
	tableNamesSorted := map[string][]string{}

	tableDataTypes := make(map[string]map[string]map[string]string)
	tableDataTypeQuery := "SELECT table_name, COLUMN_NAME, COLUMN_KEY, COLUMN_TYPE, DATA_TYPE, IS_NULLABLE, COLUMN_COMMENT FROM INFORMATION_SCHEMA.COLUMNS where TABLE_SCHEMA = ? order by table_name, ordinal_position asc"

	rows, err := co.c.Db.Query(tableDataTypeQuery, databaseName)
	if err != nil {
		return nil, nil, fmt.Errorf("Error selecting from db: " + err.Error())
	}
	if rows != nil {
		defer rows.Close()
	} else {
		return nil, nil, fmt.Errorf("Error selecting from db: " + err.Error())
	}

	nowTableName := ""
	for rows.Next() {
		var table string
		var column string
		var columnKey string
		var columnType string
		var dataType string
		var nullable string
		var comment string
		rows.Scan(&table, &column, &columnKey, &columnType, &dataType, &nullable, &comment)

		if nowTableName != table {
			tableDataTypes[table] = map[string]map[string]string{}
		}
		tableDataTypes[table][column] = map[string]string{"value": dataType, "nullable": nullable, "columnType": columnType, "primary": columnKey, "comment": comment}
		nowTableName = table
		tableNamesSorted[table] = append(tableNamesSorted[table], column)
	}

	return tableDataTypes, tableNamesSorted, nil
}
