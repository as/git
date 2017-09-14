package git

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Addr struct {
	net, addr, svc string
}

type toSpace []byte
type toNull []byte
type toLine []byte
type stringInt int
type Until struct {
	byte
}

func clean(s string) string {
	return strings.TrimSpace(s)
}

// hashpath returns a relative path used to access the
// object denoted by h.
func hashpath(h []byte) string {
	x := fmt.Sprintf("%x", h)
	return fmt.Sprintf("%02s/%38s", x[:2], x[2:])
}

func cut(buf *bytes.Buffer, delim byte) (string, error) {
	x, err := buf.ReadBytes(delim)
	if err != nil {
		return "", err
	}
	n := len(x) - 1
	if n < 0 {
		return "", fmt.Errorf("string too short: %q", x)
	}
	return string(x[:n]), nil
}

func (u Until) Read(r io.Reader) (p []byte, err error) {
	buf := new(bytes.Buffer)
	var b [1]byte
	for err == nil {
		_, err = r.Read(b[:])
		if b[0] == u.byte {
			return buf.Bytes(), nil
		}
		buf.WriteByte(b[0])
	}
	return p, err
}

func (z *toSpace) ReadBinary(r io.Reader) error {
	p, err := Until{0x20}.Read(r)
	*z = []byte(p)
	return err
}
func (z *toNull) ReadBinary(r io.Reader) error {
	p, err := Until{0x00}.Read(r)
	*z = []byte(p)
	return err
}
func (z *toLine) ReadBinary(r io.Reader) error {
	p, err := Until{0x0a}.Read(r)
	*z = []byte(p)
	return err
}

func (z toSpace) WriteBinary(w io.Writer) (err error) {
	i := bytes.IndexRune(z, rune(' '))
	if i == -1 {
		i = len(z)
	}
	_, err = w.Write(z[:i])
	return err
}
func (z toNull) WriteBinary(w io.Writer) (err error) {
	i := bytes.IndexRune(z, rune('\x00'))
	if i == -1 {
		i = len(z)
	}
	_, err = w.Write(z[:i])
	return err
}
func (z toLine) WriteBinary(w io.Writer) (err error) {
	i := bytes.IndexRune(z, rune('\x0a'))
	if i == -1 {
		i = len(z)
	}
	_, err = w.Write(z[:i])
	return err
}
func (s *stringInt) ReadBinary(r io.Reader) (err error) {
	data := make([]byte, 4)
	if _, err := r.Read(data); err != nil {
		return err
	}
	n, err := strconv.ParseInt(string(data), 16, 16)
	if n < 4 || err != nil {
		if err != nil {
			return err
		}
		return fmt.Errorf("stringInt: ReadBinary: need >= 4 bytes, got %d", n)
	}
	*s = stringInt(n - 4)
	return err
}

func (s *stringInt) WriteBinary(w io.Writer) error {
	_, err := w.Write([]byte(fmt.Sprintf("%x", *s)))
	return err
}

func (z zobject) String() string {
	return fmt.Sprintf("%s %d\x00%s", z.info.Describe(), len(z.data), z.data)
}

func nextMsgReader(r io.Reader) (lbr *bufio.Reader, err error) {
	var tmp [4]byte
	_, err = r.Read(tmp[:])
	if err != nil {
		return
	}
	dsize, err := strconv.ParseInt(string(tmp[:]), 16, 16)
	if err != nil {
		return
	}
	return bufio.NewReader(io.LimitReader(r, dsize-4)), nil
}
func printerr(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
}

func println(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
}

func no(err error) {
	if err != nil {
		printerr(err)
		os.Exit(1)
	}
}
