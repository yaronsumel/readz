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
	"io/ioutil"
	"math/rand"
	"sync"
	"testing"
)

func getDummyReader(size int64) (*bytes.Buffer, []byte) {
	var (
		bcpy = make([]byte, size)
		b    = make([]byte, size)
	)
	for k := range b {
		b[k] = byte(rand.Int())
	}
	copy(bcpy, b)
	return bytes.NewBuffer(b), bcpy
}

func TestNewUnifiedReader(t *testing.T) {

	var (
		f, b = getDummyReader(2 * 1e+6)
		rs   = NewReReader(100)
	)

	nr, err := rs.Read(f)
	if err != nil {
		t.Fatal(err)
	}

	byt, err := ioutil.ReadAll(nr)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(byt, b) {
		t.Fatal("bytes are not equal")
	}

	if !bytes.Equal(rs.Bytes(), b[:100]) {
		t.Fatal("bytes are not equal")
	}

}

func TestNewReaderSplitter(t *testing.T) {

	t.Run("one reader", func(t *testing.T) {

		var (
			n       int64 = 100
			f, b          = getDummyReader(n)
			rs            = NewReaderSplitter(f, "reader1")
			ctx, fn       = context.WithCancel(context.Background())
		)

		defer fn()
		defer rs.Close()
		go rs.Pipe(ctx)

		byt, err := ioutil.ReadAll(rs.Reader("reader1"))
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(byt, b) {
			t.Fatal("bytes are not equal")
		}

		if rs.Written() != n {
			t.Fatal("number of written bytes are not equal")
		}

		if rs.Error() != nil {
			t.Fatal(rs.Error())
		}

	})

	t.Run("two readers", func(t *testing.T) {

		var (
			n       int64 = 2 * 1e+6
			f, b          = getDummyReader(n)
			rs            = NewReaderSplitter(f, "reader1", "reader2")
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
				t.Fatal(err)
			}
			if !bytes.Equal(byt, b) {
				t.Fatal("bytes are not equal")
			}
		}()

		go func() {
			defer wg.Done()
			byt, err := ioutil.ReadAll(rs.Reader("reader2"))
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(byt, b) {
				t.Fatal("bytes are not equal")
			}
		}()

		wg.Wait()

		if rs.Written() != n {
			t.Fatal("number of written bytes are not equal")
		}

		if rs.Error() != nil {
			t.Fatal(rs.Error())
		}

	})

	t.Run("two readers - one failure", func(t *testing.T) {

		var (
			n       int64 = 4 * 1e+6
			f, b          = getDummyReader(n)
			rs            = NewReaderSplitter(f, "reader1", "reader2")
			ctx, fn       = context.WithCancel(context.Background())
		)

		go rs.Pipe(ctx)
		defer rs.Close()

		wg := sync.WaitGroup{}
		wg.Add(2)

		// read a bit and cancel the context
		go func() {
			defer wg.Done()
			r := rs.Reader("reader1")
			b := make([]byte, 100)
			_, err := r.Read(b)
			if err != nil {
				t.Fatal(err)
			}
			fn()
		}()

		go func() {
			defer wg.Done()
			byt, err := ioutil.ReadAll(rs.Reader("reader2"))
			if err == nil {
				t.Fatal("err is empty", err)
			}
			if bytes.Equal(byt, b) {
				t.Fatal("bytes are equal")
			}
		}()

		wg.Wait()

		if rs.Written() == n {
			t.Fatal("number of written bytes are equal")
		}

		if rs.Error() != nil {
			t.Fatal("error should not be empty")
		}

	})

}
