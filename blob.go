package git

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"

	"github.com/as/git/fs"
)

//wire9 leaf mode[,toSpace] file[,toNull] hash[20]

// Blob is a git object. It stores file contents.
type Blob struct {
	*leaf
	data []byte
}

func (b *Blob) ReadBinary(r io.Reader) (err error) {
	if err = b.leaf.ReadBinary(r); err != nil {
		return
	}
	return
}

func (b *Blob) Name() string {
	return string(b.file)
}

func (b *Blob) Hash() Hash {
	h, err := HashData(b.Data())
	if err != nil {
		panic(err)
	}
	return h
}

func (b *Blob) Data() []byte {
	return b.data
}

type Repo struct {
	*fs.Fs
}

func NewRepo(path string) *Repo {
	return &Repo{fs.New(path)}
}

func (r *Repo) Eval(key string) (Object, error) {
	rc, err := r.Open(fs.Key(key))
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	buf := bytes.NewBuffer(data)
	var typ, asize string

	if typ, err = cut(buf, '\x20'); err != nil {
		return nil, err
	}
	if asize, err = cut(buf, '\x00'); err != nil {
		return nil, err
	}

	//
	// TODO: Finish Eval, and replace usage of hash with fs.FS
	//

	size, err := strconv.Atoi(string(asize))
	if err != nil {
		return nil, err
	}
	if x := len(buf.Bytes()); size != x {
		return nil, fmt.Errorf("wrong size: %d/%d bytes", size, x)
	}

	switch string(typ) {
	/*
		case "blob":
			fmt.Println("blob")
			return &Blob{data: buf.Bytes()}, nil
		case "tree":
			t, err := TreeFromFile(name)
			if err != nil{
				return nil, err
			}
			return t, nil
	*/
	}
	return nil, fmt.Errorf("no format found for: %v", r)
}

// Eval evaluates a leaf (object) and returns
// and interface to that object.
func (l *leaf) Eval() (Object, error) {
	name := hashpath(l.hash)
	name = "testdata/.git/objects/" + name
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(data)
	var typ, asize string

	if typ, err = cut(buf, '\x20'); err != nil {
		return nil, err
	}
	if asize, err = cut(buf, '\x00'); err != nil {
		return nil, err
	}

	size, err := strconv.Atoi(string(asize))
	if err != nil {
		return nil, err
	}
	if x := len(buf.Bytes()); size != x {
		return nil, fmt.Errorf("wrong size: %d/%d bytes", size, x)
	}

	switch string(typ) {
	case "blob":
		fmt.Println("found a blob")
		return &Blob{leaf: l, data: buf.Bytes()}, nil
	case "tree":
		t, err := TreeFromFile(name)
		if err != nil {
			return nil, err
		}
		return t, nil
	}
	return nil, fmt.Errorf("no format found for: %v", l)
}
