package git

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

const (
	_ PackType = iota
	PackCommit
	PackTree
	PackBlob
	PackTag
	_
	PackOffDelta
	PackRefDelta
)

type PackType uint64

type PackList struct {
	head   Hashref
	opts   Options
	Refmap map[string]Hashref
	order  []Hashref
	last   Hashref
	up     map[Hashref]Hashref
}

func (p packhdr) Objects() []zobject {
	return p.objects
}

type Info struct {
	PackType
	Size int
}

type Hashref struct {
	Hash, Ref string
}

func (p PackType) Describe() string {
	switch p {
	case PackCommit:
		return "commit"
	case PackTree:
		return "tree"
	case PackBlob:
		return "blob"
	case PackTag:
		return "tag"
	}
	return ""
}

func (p *PackList) Print() {
	for i, v := range p.order {
		fmt.Printf("%04d %s (%s)\n", i, v.Hash, v.Ref)
	}
	e := Hashref{}
	for k := p.order[len(p.order)-1]; ; {
		v := p.up[k]
		if v == e {
			fmt.Println("done")
			break
		}
		fmt.Printf("%-20s (%5s) -> %-20s (%5s)\n", k.Hash, k.Ref, v.Hash, v.Ref)
		k = v
	}
}

func (p *PackList) Read(r io.Reader) error {
	lbr, err := nextMsgReader(r)
	if err != nil {
		return err
	}
	if p.head.Hash, err = lbr.ReadString(' '); err != nil {
		return err
	}
	if p.head.Ref, err = lbr.ReadString('\x00'); err != nil {
		return err
	}
	if p.head.Ref != "HEAD\x00" {
		return fmt.Errorf("expected HEAD\x00, got %q", p.head.Ref)
	}
	p.head.Ref = clean(p.head.Ref)

	rawopts, err := lbr.ReadString('\n')
	if err != nil {
		return err
	}
	p.opts = Options(strings.Split(rawopts[:len(rawopts)-1], " "))
	p.Refmap = make(map[string]Hashref)
	p.up = make(map[Hashref]Hashref)
	p.order = make([]Hashref, 0)

	for {
		var hr Hashref
		if lbr, err = nextMsgReader(r); err != nil {
			break
		}
		if hr.Hash, err = lbr.ReadString(' '); err != nil {
			break
		}
		if hr.Ref, err = lbr.ReadString('\n'); err != nil {
			break
		}
		hr.Hash, hr.Ref = clean(hr.Hash), clean(hr.Ref)
		if r := hr.Ref; len(r) > 3 && r[len(r)-3:] == "^{}" {
			hr.Ref = hr.Ref[:len(r)-3]
			p.up[hr] = p.last
		}
		p.Refmap[hr.Hash] = hr
		p.Refmap[hr.Ref] = hr
		p.order = append(p.order, hr)
		p.last = hr
	}
	if err != nil && err != io.EOF {
		return err
	} else if err == io.EOF {
		err = nil
	}
	return nil
}

func (a *Info) ReadBinary(r io.Reader) (err error) {
	var peek []byte
	switch r := r.(type) {
	case *bytes.Buffer:
		peek = r.Bytes()[:10]
	case io.Reader:
		return fmt.Errorf("info must have buffered reader, have %T", r)
	}
	if err != nil {
		return
	}
	u, n := a.Uvarint(peek)
	_, err = r.Read(peek[:n])
	if err != nil {
		return err
	}
	a.Size = int(u)
	return nil
}
func (a *Info) Uvarint(buf []byte) (uint64, int) {
	var x uint64
	s := uint64(0)
	a.PackType = PackType((buf[0] >> 4) & 0x7)
	fmt.Println("The type is", a.PackType)
	for i, b := range buf {
		if i == 0 {
			b = b &^ 0x70
		}
		if b < 0x80 {
			if i > 9 || i == 9 && b > 1 {
				return 0, -(i + 1) // overflow
			}
			return x | uint64(b)<<s, i + 1
		}
		x |= uint64(b&0x7f) << s
		if i > 0 {
			s += 3
		}
		s += 4
	}
	return 0, 0
}
func (a *Info) WriteBinary(w io.Writer) error {
	panic("no")
}
