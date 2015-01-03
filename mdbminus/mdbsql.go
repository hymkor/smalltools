package main

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"strings"
)

func mdbSql(db *sql.DB, sqlStr string, writer io.Writer) error {
	fmt.Fprintln(writer, sqlStr)
	if strings.HasPrefix(strings.TrimSpace(strings.ToLower(sqlStr)), "select") {
		rows, err := db.Query(sqlStr)
		if err != nil {
			return err
		}
		defer rows.Close()
		cols, err := rows.Columns()
		if err != nil {
			return err
		}
		values := make([]interface{}, len(cols), len(cols))
		pValues := make([]interface{}, len(cols), len(cols))
		for i := 0; i < len(cols); i++ {
			if i > 0 {
				fmt.Fprint(writer, ",")
			}
			fmt.Fprint(writer, cols[i])
			pValues[i] = &values[i]
		}
		fmt.Fprintln(writer)
		for rows.Next() {
			if err := rows.Scan(pValues...); err != nil {
				return err
			}
			for i, cell := range values {
				if i > 0 {
					fmt.Fprint(writer, ",")
				}
				fmt.Fprintf(writer, "%v", cell)
			}
			fmt.Fprintln(writer)
		}
	} else {
		_, resultErr := db.Exec(os.Args[2])
		if resultErr != nil {
			return fmt.Errorf("%s: %s", sqlStr, resultErr)
		}
		fmt.Fprintln(writer, "done")
	}
	return nil
}
