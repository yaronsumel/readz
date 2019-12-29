package main

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/yaronsumel/readz"
)

// split one file into two files
func main() {

	dir, err := os.Getwd()
	if err != nil {
		log.Panicln(err)
	}

	f, err := os.Open(dir + "/readz.go")
	if err != nil {
		log.Panicln(err)
	}

	var (
		ctx, fn = context.WithCancel(context.Background())
		rs      = readz.NewReaderSplitter(f, "file1", "file2")
		wg      = sync.WaitGroup{}
	)

	defer fn()

	go rs.Pipe(ctx)
	wg.Add(2)

	go func() {
		defer wg.Done()
		handle("file1", rs)
	}()

	go func() {
		defer wg.Done()
		handle("file2", rs)
	}()

	wg.Wait()

}

func handle(fileName string, rs *readz.ReaderSplitter) {
	f, err := ioutil.TempFile("/tmp", fileName+"-*.go")
	if err != nil {
		log.Panicln(err)
		return
	}
	r := rs.Reader(fileName)
	defer r.Close()
	_, err = io.Copy(f, r)
	if err != nil {
		log.Panicln(err)
	}
	fi, err := f.Stat()
	if err != nil {
		log.Panicln(err)
	}
	log.Println("/tmp/" + fi.Name())
}
