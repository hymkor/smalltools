package main

import "archive/zip"
import "path"
import "flag"
import "fmt"
import "io"
import "os"

var listFlag = flag.Bool("l", false, "listing")
var directory = flag.String("d", "", "expand directory")

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		return
	}
	zipFileName := args[0]
	reader, readerErr := os.Open(zipFileName)
	if readerErr != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", zipFileName, readerErr.Error())
		return
	}
	defer reader.Close()
	finfo, finfoErr := reader.Stat()
	if finfoErr != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", zipFileName, finfoErr.Error())
		return
	}
	zipReader, zipReaderErr := zip.NewReader(reader, finfo.Size())
	if zipReaderErr != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", zipFileName, zipReaderErr.Error())
		return
	}
	if len(*directory) > 0 {
		os.Chdir(*directory)
	}
	files := map[string]struct{}{}
	for _, fname := range args[1:] {
		files[fname] = struct{}{}
	}
	for _, f := range zipReader.File {
		if len(files) > 0 {
			if _, ok := files[f.Name]; !ok {
				continue
			}
		}
		if *listFlag {
			fmt.Println(f.Name)
			continue
		}
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(f.Name, 0777); err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s\n", f.Name, err.Error())
			}
			fmt.Fprintln(os.Stdout, f.Name)
			continue
		}
		zipFileReader, zipFileReaderErr := f.Open()
		if zipFileReaderErr != nil {
			fmt.Fprintf(os.Stderr, "%s: %s: %s\n",
				zipFileName,
				f.Name,
				zipFileReaderErr.Error(),
			)
		} else {
			if err := os.MkdirAll(path.Dir(f.Name), 0777); err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s\n", f.Name, err.Error())
			} else if unzipWriter, unzipWriterErr := os.Create(f.Name); unzipWriterErr != nil {
				fmt.Fprintf(os.Stderr, "%s: %s: %s\n",
					zipFileName,
					f.Name,
					unzipWriterErr.Error())
			} else {
				_, err := io.Copy(unzipWriter, zipFileReader)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s: %s: %s\n",
						zipFileName,
						f.Name,
						err.Error())
				} else {
					fmt.Println(f.Name)
				}
				unzipWriter.Close()
			}
		}
		zipFileReader.Close()
	}
}
