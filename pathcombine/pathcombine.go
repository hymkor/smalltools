package main

import "fmt"
import "os"
import "github.com/zetamatta/nyagos/src/dos"

func main(){
	fmt.Println(dos.Join(os.Args[1:]...))
}
