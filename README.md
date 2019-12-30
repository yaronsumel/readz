# readz [![Go Report Card](https://goreportcard.com/badge/github.com/yaronsumel/readz)](https://goreportcard.com/report/github.com/yaronsumel/readz) [![GoDoc](https://godoc.org/github.com/yaronsumel/readz?status.svg)](https://godoc.org/github.com/yaronsumel/readz) [![Actions Status](https://github.com/yaronsumel/readz/workflows/Go/badge.svg)](https://github.com/yaronsumel/readz/actions)

small util that helps you read stuff from your readers

### ReReader 

ReReader gives you the ability to read data from io.Reader and to
get another reader which has the same data. the meaning is that you can
re-read the same data you have just read and use it again without data loss.

```go

// create new rereader with 100 byes buffer 
rereader := readz.NewReReader(100)

// reader from r into your buffer and return fresh reader with all data
newReader, err := rereader.Read(r)
if err != nil {
	//...
}

// do something with your bytes
rs.Bytes()

byt, err := ioutil.ReadAll(newReader)
if err != nil {
  //...
}

```

### ReaderSplitter

ReaderSplitter allows you to read data from one reader into multiple pipes.
data will copied to pipes writers and exposed back by the Reader method.

```go

var (
    rs            = readz.NewReaderSplitter(f, "reader1", "reader2")
    ctx, fn       = context.WithCancel(context.Background())
)

defer fn()
defer rs.Close()
go rs.Pipe(ctx)

wg := sync.WaitGroup{}
wg.Add(2)

go func() {
	defer wg.Done()
	byt, err := ioutil.ReadAll(rs.Reader("reader1"))
	if err != nil {
		t.Error(err)
	}
// do something with byt
}()

go func() {
	defer wg.Done()
	byt, err := ioutil.ReadAll(rs.Reader("reader2"))
	if err != nil {
		t.Error(err)
	}
// do something with byt
}()

wg.Wait()

```

### Examples

Check out the examples [directory](https://github.com/yaronsumel/readz/blob/master/example) 

### License

```
MIT License

Copyright (c) 2019 Yaron Sumel

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

> ##### Written and Maintained by [@YaronSumel](https://twitter.com/yaronsumel) #####
