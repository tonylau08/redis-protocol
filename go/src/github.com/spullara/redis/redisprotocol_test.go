package redis

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func Test_readLong(t *testing.T) {
	result, err := readLong(bufio.NewReader(strings.NewReader("123456789\r\n")))
	if err != nil {
		t.Error("Read error", err)
	} else if 123456789 != result {
		t.Error("Not equal: " + fmt.Sprintf("%d", result))
	}
	result, err = readLong(bufio.NewReader(strings.NewReader("-123456789\r\n")))
	if err != nil {
		t.Error("Read error", err)
	} else if -123456789 != result {
		t.Error("Not equal: " + fmt.Sprintf("%d", result))
	}
}

func Test_readBytes(t *testing.T) {
	result, err := readBytes(bufio.NewReader(strings.NewReader("3\r\nSam\r\n")))
	if err != nil {
		t.Error("Read error", err)
	} else if !bytes.Equal([]byte("Sam"), result) {
		t.Error("Not equal: " + string(result))
	}
}

func Benchmark_freelsBench(b *testing.B) {
	buffer := bytes.NewBuffer(make([]byte, 0))
	buffer.WriteByte(MultiBulkMarker)
	buffer.WriteString("100\r\n")
	for i := 0; i < 100; i++ {
		buffer.WriteByte(BulkMarker)
		buffer.WriteString("6\r\n")
		buffer.WriteString("foobar\r\n")
	}
	br := bytes.NewReader(buffer.Bytes())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reply, err := Receive(br)
		if err != nil {
			b.Error("Failed", err)
		}
		multibulk, ok := reply.(*MultiBulkReply)
		if !ok {
			b.Error("Wrong type", multibulk, reply)
		}
		if len(multibulk.replies) != 100 {
			b.Error("Invalid number of replies", multibulk.replies)
		}
		br.Seek(0, 0)
	}
}
