package main

import "os"
import "bufio"
import "fmt"
import "strings"

func main() {
	total := 0
	for _, arg1 := range os.Args {
		arg1_ := strings.ToLower(arg1)
		if !strings.HasSuffix(arg1_, ".vb") {
			continue
		}
		if strings.Contains(arg1_, "designer") {
			continue
		}
		file, err := os.Open(arg1)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err.Error())
			continue
		}
		defer file.Close()
		reader := bufio.NewReader(file)
		count := 0
		for {
			_, _, err := reader.ReadLine()
			if err != nil {
				break
			}
			count += 1
			total += 1
		}
		fmt.Printf("%8d : %s\n", count, arg1)
	}
	fmt.Printf("%8d : %s\n", total, "*** Total ***")
}
