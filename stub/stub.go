package main

import "fmt"
import "os"
import "strings"

import "github.com/zetamatta/nyagos/Src/dos"

func main(){
	exeName,_ := dos.GetModuleFileName()
	logName := strings.Replace(exeName,".exe",".out",-1)
	logFp,logFpErr := os.Create(logName)
	if logFpErr != nil {
		fmt.Fprintf(os.Stderr,"%s: %s\n",os.Args[0],logFpErr.Error())
		return
	}
	for i,arg1 := range os.Args {
		fmt.Fprintf(logFp,"[%d]='%s'\n",i,arg1)
	}
	logFp.Close()
}
