package git

import "bytes"

type Ref struct {
	name string
	hash Hash
}

// ParseRef parses a named hash reference from the byte buffer.
func ParseRef(b *bytes.Buffer) (r Ref, err error) {
	var line string
	if line, err = b.ReadString(' '); err != nil {
		return
	}
	r.name = line
	r.hash, err = ParseHashString(string(b.Next(40)))
	return
}
