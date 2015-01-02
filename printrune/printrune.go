package main

import "fmt"
import "os"

func main() {
	for _, arg1 := range os.Args[1:] {
		for _, r := range arg1 {
			fmt.Printf("[%0X]", r)
		}
		fmt.Print("\n")
	}
}
