package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-adodb"
)

func main(){
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr,"Usage: %s MDBFILENAME SQL\n",os.Args[0])
		return
	}
	if _,err := os.Stat(os.Args[1]) ; err != nil {
		fmt.Fprintf(os.Stderr,"%s: Not Found(%s)",os.Args[1],err)
		return
	}
	db,dbErr := sql.Open("adodb",
		"Provider=Microsoft.Jet.OLEDB.4.0;Data Source="+os.Args[1])
	if dbErr != nil {
		fmt.Fprintln(os.Stderr,dbErr)
		return
	}
	defer db.Close()
	fmt.Println(os.Args[2])
	if strings.HasPrefix(strings.ToLower(os.Args[2]),"select") {
		rows,err := db.Query(os.Args[2])
		if err != nil {
			fmt.Fprintln(os.Stderr,err)
			return
		}
		defer rows.Close()
		cols,err := rows.Columns()
		if err != nil {
			fmt.Fprintln(os.Stderr,err)
			return
		}
		values := make([]interface{},len(cols),len(cols))
		pValues := make([]interface{},len(cols),len(cols))
		for i := 0 ; i < len(cols) ; i++ {
			if i > 0 {
				fmt.Print(",")
			}
			fmt.Print(cols[i])
			pValues[i] = &values[i]
		}
		fmt.Println()
		for rows.Next() {
			if err := rows.Scan(pValues...) ; err != nil {
				fmt.Fprintln(os.Stderr,err)
				return
			}
			for i,cell := range values {
				if i > 0 {
					fmt.Print(",")
				}
				fmt.Printf("%v",cell)
			}
			fmt.Println()
		}
	}else{
		_,resultErr := db.Exec(os.Args[2])
		if resultErr != nil {
			fmt.Fprintf(os.Stderr,"%s: %s\n",os.Args[2],resultErr.Error())
			return
		}
		fmt.Println("done")
	}
}
