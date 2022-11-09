package gen

import (
	"github.com/tama1029/gfdb/client"
	"github.com/tama1029/gfdb/db"
	"github.com/tama1029/gfdb/output"
	"github.com/tama1029/gfdb/render"
)

type GenGoaResultType struct {
	tableInfo              *db.TableInfo
	GenGoaResultTypeRender *render.GenGoaResultTypeRender
	GenGoaResultTypeOutput *output.GenOutput
}

func NewGenGoaResultType(host, user, pass, database string, port int, outputd string) (*GenGoaResultType, error) {
	c, err := client.NewClient(user, pass, host, database, port)
	if err != nil {
		return nil, err
	}
	ti := db.NewTableInfo(c)
	gsr := render.NewGenGoaResultTypeRender()
	gso := output.NewGenOutput(outputd)
	return &GenGoaResultType{
		tableInfo:              ti,
		GenGoaResultTypeRender: gsr,
		GenGoaResultTypeOutput: gso,
	}, nil
}

func (g GenGoaResultType) Generate(database string) error {
	tableDataTypes, tableNamesSorted, err := g.tableInfo.GetColumnInfo(database)
	if err != nil {
		return err
	}
	tableAndComment, err := g.tableInfo.GetTableInfo(database)
	if err != nil {
		return err
	}

	renders, err := g.GenGoaResultTypeRender.RenderFacade(tableDataTypes, tableNamesSorted, tableAndComment)
	if err != nil {
		return err
	}

	err = g.GenGoaResultTypeOutput.ToFile(renders, "goa_result_type")
	if err != nil {
		return err
	}
	return nil
}
