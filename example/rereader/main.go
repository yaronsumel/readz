// MIT License
//
// Copyright (c) 2019 Yaron Sumel
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/yaronsumel/readz"
)

// read file and check its content type by reading the first 512 bytes and use rereader to read all the file
func main() {

	dir, err := os.Getwd()
	if err != nil {
		log.Panicln(err)
	}

	f, err := os.Open(dir + "/main.go")
	if err != nil {
		log.Panicln(err)
	}

	defer f.Close()

	reReader := readz.NewReReader(512)

	nr, err := reReader.Read(f)
	if err != nil {
		log.Panicln(err)
	}

	switch http.DetectContentType(reReader.Bytes()) {
	case "text/plain; charset=utf-8":
		byt, err := ioutil.ReadAll(nr)
		if err != nil {
			log.Panicln(err)
		}
		log.Println(string(byt))
	}

}
