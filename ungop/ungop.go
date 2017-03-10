package main

import "archive/zip"
import "path"
import "flag"
import "fmt"
import "io"
import "os"

var listFlag = flag.Bool("l", false, "listing")
var directory = flag.String("d", "", "expand directory")

func deepMkdir(folder string) error {
	finfo, err := os.Stat(folder)
	if err == nil {
		if finfo.IsDir() {
			return nil
		} else {
			return fmt.Errorf("%s: Not Directory", folder)
		}
	} else {
		parent := path.Dir(folder)
		if err := deepMkdir(parent); err != nil {
			return err
		}
		if _, err2 := os.Stat(folder); err2 != nil {
			return os.Mkdir(folder, 0666)
		} else {
			return nil
		}
	}
}

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
	files := map[string]bool{}
	for _, fname := range args[1:] {
		files[fname] = true
	}
	for _, f := range zipReader.File {
		if len(files) > 0 && !files[f.Name] {
			continue
		}
		if *listFlag {
			dos_dt := f.ModifiedDate
			dos_year := (dos_dt >> 9)
			dos_month := (dos_dt >> 5 ) & 0x0F
			dos_day := (dos_dt & 0x1F)

			fmt.Printf("%s (1980+%02d-%02d-%02d) %s\n",
				f.ModTime(),
				dos_year,
				dos_month,
				dos_day,
				f.Name)
			continue
		}
		if f.FileInfo().IsDir() {
			if err := deepMkdir(f.Name); err != nil {
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
			if err := deepMkdir(path.Dir(f.Name)); err != nil {
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
