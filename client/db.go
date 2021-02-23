package client

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

type Client struct {
	Db *sql.DB
}

func NewClient(user, password, host, database string, port int) (*Client, error) {
	c := &Client{}

	db, err := sql.Open("mysql", user+":"+password+"@tcp("+host+":"+strconv.Itoa(port)+")/"+database+"?&parseTime=True")
	if err != nil {
		return nil, fmt.Errorf("Error opening mysql db: " + err.Error())
	}
	c.Db = db
	return c, nil
}
