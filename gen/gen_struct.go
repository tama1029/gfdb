package gen

import (
	"github.com/tama1029/gfdb/client"
	"github.com/tama1029/gfdb/db"
	"github.com/tama1029/gfdb/output"
	"github.com/tama1029/gfdb/render"
)

type GenStruct struct {
	tableInfo       *db.TableInfo
	genStructRender *render.GenStructRender
	genStructOutput *output.GenStructOutput
}

func NewGenStruct(host, user, pass, database string, port int, outputd string) (*GenStruct, error) {
	c, err := client.NewClient(user, pass, host, database, port)
	if err != nil {
		return nil, err
	}
	ti := db.NewTableInfo(c)
	gsr := render.NewGenStructRender()
	gso := output.NewGenStructOutput(outputd)
	return &GenStruct{
		tableInfo:       ti,
		genStructRender: gsr,
		genStructOutput: gso,
	}, nil
}

func (g GenStruct) Generate(database string) error {
	tableDataTypes, tableNamesSorted, err := g.tableInfo.GetTableInfo(database)
	if err != nil {
		return err
	}

	renders, err := g.genStructRender.RenderFacade(tableDataTypes, tableNamesSorted)
	if err != nil {
		return err
	}

	err = g.genStructOutput.ToFile(renders)
	if err != nil {
		return err
	}
	return nil
}
