package codec

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"

	"google.golang.org/protobuf/proto"
)

// Marshal / Unmarshal 基础封装
func isNilMessage(m proto.Message) bool {
	if m == nil {
		return true
	}
	v := reflect.ValueOf(m)
	return v.Kind() == reflect.Ptr && v.IsNil()
}

func Marshal(m proto.Message) ([]byte, error) {
	if isNilMessage(m) {
		return nil, fmt.Errorf("nil proto message")
	}
	return proto.Marshal(m)
}

func Unmarshal(b []byte, m proto.Message) error {
	if isNilMessage(m) {
		return fmt.Errorf("nil proto message")
	}
	return proto.Unmarshal(b, m)
}

// Frame：长度前缀（uint32 big endian）+ payload
// 适合 TCP、文件流、消息队列里做边界分隔
func WriteFrame(w io.Writer, m proto.Message) error {
	b, err := Marshal(m)
	if err != nil {
		return err
	}

	var hdr [4]byte
	binary.BigEndian.PutUint32(hdr[:], uint32(len(b)))

	// bufio.NewWriter 可以由调用方自己包，这里保持简单
	if _, err := w.Write(hdr[:]); err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

func ReadFrame(r *bufio.Reader, m proto.Message) error {
	var hdr [4]byte
	if _, err := io.ReadFull(r, hdr[:]); err != nil {
		return err
	}

	n := binary.BigEndian.Uint32(hdr[:])
	if n == 0 {
		return fmt.Errorf("empty frame")
	}
	// 可选：做一个最大长度限制，避免 OOM
	if n > 32<<20 { // 32MB
		return fmt.Errorf("frame too large: %d", n)
	}

	buf := make([]byte, n)
	if _, err := io.ReadFull(r, buf); err != nil {
		return err
	}
	return Unmarshal(buf, m)
}
