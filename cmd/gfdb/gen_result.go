package gfdb

import (
	"github.com/spf13/cobra"
	"github.com/tama1029/gfdb/gen"
)

func GenResultCmd() *cobra.Command {
	var host string
	var port int
	var user string
	var pass string
	var database string
	var outputd string

	cmd := &cobra.Command{
		Use:   "result",
		Short: "result from database",
		RunE: func(cmd *cobra.Command, args []string) error {
			gs, err := gen.NewGenResult(host, user, pass, database, port, outputd)
			if err != nil {
				return err
			}
			err = gs.Generate(database)
			if err != nil {
				return err
			}

			return nil
		},
	}
	cmd.Flags().SortFlags = false
	cmd.Flags().StringVar(&host, "host", "localhost", "database host")
	cmd.Flags().IntVar(&port, "port", 3306, "database port")
	cmd.Flags().StringVar(&user, "user", "admin", "database user")
	cmd.Flags().StringVar(&database, "database", "", "database name")
	cmd.Flags().StringVar(&pass, "pass", "", "database password")
	cmd.Flags().StringVar(&outputd, "outputd", "example", "output directory path")

	return cmd
}
