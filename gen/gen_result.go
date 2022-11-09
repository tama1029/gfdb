package gen

import (
	"github.com/tama1029/gfdb/client"
	"github.com/tama1029/gfdb/db"
	"github.com/tama1029/gfdb/output"
	"github.com/tama1029/gfdb/render"
)

type GenResult struct {
	tableInfo       *db.TableInfo
	GenResultRender *render.GenResultRender
	GenResultOutput *output.GenOutput
}

func NewGenResult(host, user, pass, database string, port int, outputd string) (*GenResult, error) {
	c, err := client.NewClient(user, pass, host, database, port)
	if err != nil {
		return nil, err
	}
	ti := db.NewTableInfo(c)
	gsr := render.NewGenResultRender()
	gso := output.NewGenOutput(outputd)
	return &GenResult{
		tableInfo:       ti,
		GenResultRender: gsr,
		GenResultOutput: gso,
	}, nil
}

func (g GenResult) Generate(database string) error {
	tableDataTypes, tableNamesSorted, err := g.tableInfo.GetColumnInfo(database)
	if err != nil {
		return err
	}

	renders, err := g.GenResultRender.RenderFacade(tableDataTypes, tableNamesSorted)
	if err != nil {
		return err
	}

	err = g.GenResultOutput.ToFile(renders, "result")
	if err != nil {
		return err
	}
	return nil
}
