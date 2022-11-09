package gen

import (
	"github.com/tama1029/gfdb/client"
	"github.com/tama1029/gfdb/db"
	"github.com/tama1029/gfdb/output"
	"github.com/tama1029/gfdb/render"
)

type GenGoaType struct {
	tableInfo        *db.TableInfo
	GenGoaTypeRender *render.GenGoaTypeRender
	GenGoaTypeOutput *output.GenOutput
}

func NewGenGoaType(host, user, pass, database string, port int, outputd string) (*GenGoaType, error) {
	c, err := client.NewClient(user, pass, host, database, port)
	if err != nil {
		return nil, err
	}
	ti := db.NewTableInfo(c)
	gsr := render.NewGenGoaTypeRender()
	gso := output.NewGenOutput(outputd)
	return &GenGoaType{
		tableInfo:        ti,
		GenGoaTypeRender: gsr,
		GenGoaTypeOutput: gso,
	}, nil
}

func (g GenGoaType) Generate(database string) error {
	tableDataTypes, tableNamesSorted, err := g.tableInfo.GetColumnInfo(database)
	if err != nil {
		return err
	}

	renders, err := g.GenGoaTypeRender.RenderFacade(tableDataTypes, tableNamesSorted)
	if err != nil {
		return err
	}

	err = g.GenGoaTypeOutput.ToFile(renders, "goa_type")
	if err != nil {
		return err
	}
	return nil
}
