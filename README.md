# readz
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
