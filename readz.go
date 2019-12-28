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

package readz

import (
	"bytes"
	"context"
	"io"
)

// ReReader gives you the ability to read data from io.Reader and to
// get another reader which has the same data. the meaning is that you can
// re-read the same data you have just read and use it again without data loss.
type ReReader struct {
	data []byte
}

// NewReReader creates new ReReader instance with data slice at the size specified.
func NewReReader(size int64) *ReReader {
	return &ReReader{
		data: make([]byte, size),
	}
}

// Read reads from r into data. creates new multi-reader and return it or error
func (br *ReReader) Read(r io.Reader) (io.Reader, error) {
	_, err := r.Read(br.data)
	if err != nil {
		return nil, err
	}
	return io.MultiReader(bytes.NewReader(br.data), r), nil
}

// Bytes returns the data read from r
func (br *ReReader) Bytes() []byte {
	return br.data
}

// ReaderSplitter allows you to read data from one reader into multiple pipes.
// data will copied to pipes writers and exposed back by the Reader method.
type ReaderSplitter struct {
	reader  io.Reader
	writers []io.Writer
	pipes   map[string]struct {
		pr *io.PipeReader
		pw *io.PipeWriter
	}
	n                            int64
	err                          error
	closedReaders, closedWriters bool
}

// NewReaderSplitter creates new ReaderSplitter with the right amount of pipes
func NewReaderSplitter(r io.Reader, pipeNames ...string) *ReaderSplitter {
	rp := &ReaderSplitter{
		reader: r,
		pipes: map[string]struct {
			pr *io.PipeReader
			pw *io.PipeWriter
		}{},
	}
	for _, name := range pipeNames {
		pr, pw := io.Pipe()
		rp.pipes[name] = struct {
			pr *io.PipeReader
			pw *io.PipeWriter
		}{pr, pw}
		rp.writers = append(rp.writers, pw)
	}
	return rp
}

// Close all writers and readers that still open
func (rs *ReaderSplitter) Close() {
	if !rs.closedReaders {
		rs.closeReaders()
	}
}

// closeWriters all writers
func (rs *ReaderSplitter) closeWriters() {
	for _, v := range rs.pipes {
		v.pw.Close()
	}
	rs.closedWriters = true
}

// closeReaders all readers
func (rs *ReaderSplitter) closeReaders() {
	for _, v := range rs.pipes {
		v.pr.Close()
	}
	rs.closedReaders = true
}

// Pipe read from rs.reader into pipe writers
func (rs *ReaderSplitter) Pipe(ctx context.Context) (n int64, err error) {
	defer rs.closeWriters()
	done := make(chan struct{})
	go func() {
		n, err = io.Copy(io.MultiWriter(rs.writers...), rs.reader)
		done <- struct{}{}
	}()
	select {
	case <-done:
	case <-ctx.Done():
		rs.closeReaders()
	}
	return
}

// Reader returns reader by its name
func (rs *ReaderSplitter) Reader(name string) io.Reader {
	return rs.pipes[name].pr
}
