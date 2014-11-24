package main

import "archive/zip"
import "fmt"
import "io"
import "os"

func main() {
	if len(os.Args) < 2 {
		return
	}
	reader, readerErr := os.Open(os.Args[1])
	if readerErr != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[1], readerErr.Error())
		return
	}
	defer reader.Close()
	finfo, finfoErr := reader.Stat()
	if finfoErr != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[1], finfoErr.Error())
		return
	}
	zipReader, zipReaderErr := zip.NewReader(reader, finfo.Size())
	if zipReaderErr != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[1], zipReaderErr.Error())
		return
	}
	for _, f := range zipReader.File {
		zipFileReader, zipFileReaderErr := f.Open()
		if zipFileReaderErr != nil {
			fmt.Fprintf(os.Stderr, "%s: %s: %s\n",
				os.Args[1],
				f.Name,
				zipFileReaderErr.Error())
			continue
		}
		unzipWriter, unzipWriterErr := os.Create(f.Name)
		if unzipWriterErr != nil {
			fmt.Fprintf(os.Stderr, "%s: %s: %s\n",
				os.Args[1],
				f.Name,
				unzipWriterErr.Error())
		} else {
			_, err := io.Copy(unzipWriter, zipFileReader)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s: %s\n",
					os.Args[1],
					f.Name,
					err.Error())
			} else {
				fmt.Println(f.Name)
			}
			zipFileReader.Close()
		}
		unzipWriter.Close()
	}
}
