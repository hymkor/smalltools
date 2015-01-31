package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-adodb"
	"github.com/zetamatta/nyagos/src/conio"
)

var optionE = flag.String("e", "", "SQL")
var optionF = flag.String("f", "", "Script")

func script(db *sql.DB, fname string) error {
	reader, readerErr := os.Open(fname)
	if readerErr != nil {
		return readerErr
	}
	defer reader.Close()
	scnr := bufio.NewScanner(reader)
	for scnr.Scan() {
		text := scnr.Text()
		if text[0] == '@' {
			continue
		}
		if err := mdbSql(db, text, os.Stdout); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}
	return nil
}

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Fprintf(os.Stderr,
			"Usage: %s MDBFILENAME [SCRIPTFILENAME]\n"+
				"       %s -e SQL MDBFILENAME\n",
			os.Args[0], os.Args[0])
		return
	}
	if _, err := os.Stat(args[0]); err != nil {
		fmt.Fprintf(os.Stderr, "%s: Not Found(%s)", args[0], err)
		return
	}
	db, dbErr := sql.Open("adodb",
		"Provider=Microsoft.Jet.OLEDB.4.0;Data Source="+args[0])
	if dbErr != nil {
		fmt.Fprintln(os.Stderr, dbErr)
		return
	}
	defer db.Close()
	if optionE != nil && *optionE != "" {
		if err := mdbSql(db, *optionE, os.Stdout); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}
	useShell := false
	if optionF != nil && *optionF != "" {
		if err := script(db, *optionF); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		useShell = true
	}
	if len(args) >= 2 {
		if err := script(db, args[2]); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		useShell = true
	}
	if !useShell {
		for {
			text, result := conio.ReadLinePromptStr("SQL> ")
			if !result {
				break
			}
			text = strings.TrimSpace(text)
			if text == "" {
				continue
			}
			if err := mdbSql(db, text, os.Stdout); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
			}
		}
	}
}
