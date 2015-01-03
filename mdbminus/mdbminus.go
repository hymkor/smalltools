package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-adodb"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s MDBFILENAME SQL\n", os.Args[0])
		return
	}
	if _, err := os.Stat(os.Args[1]); err != nil {
		fmt.Fprintf(os.Stderr, "%s: Not Found(%s)", os.Args[1], err)
		return
	}
	db, dbErr := sql.Open("adodb",
		"Provider=Microsoft.Jet.OLEDB.4.0;Data Source="+os.Args[1])
	if dbErr != nil {
		fmt.Fprintln(os.Stderr, dbErr)
		return
	}
	defer db.Close()
	if err := mdbSql(db, os.Args[2], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
