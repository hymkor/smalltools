package main

import "fmt"
import "os"
import "path/filepath"

func main(){
	fmt.Println(filepath.Join(os.Args[1:]...))
}
