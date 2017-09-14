package git

import (
	"fmt"
	"io"
	"os"
)

//wire9 tree head[,toSpace] size[,toNull] leaf[size,[]*leaf]

func (t tree) Leaves() []*leaf {
	return t.leaf
}

type Tree struct {
	data map[string]*Object
	tree
}

func (t *tree) Hash() Hash {
	panic("fuck")
}

func (t *tree) Data() []byte {
	return nil
}

func TreeFromFile(name string) (*tree, error) {
	fd, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	t := &tree{}
	if err = t.ReadBinary(fd); err != nil {
		return nil, err
	}
	return t, nil
}

func (z *tree) ReadBinary(r io.Reader) (err error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		if r, ok := r.(error); ok {
			err = r
			return
		}
		panic(r)
	}()
	if z == nil {
		return fmt.Errorf("ReadBinary: z nil")
	}
	if err := z.head.ReadBinary(r); err != nil {
		return err
	}
	if err := z.size.ReadBinary(r); err != nil {
		return err
	}

	z.leaf = make([]*leaf, 0)
	for err == nil {
		l := leaf{}
		if err = (&l).ReadBinary(r); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		z.leaf = append(z.leaf, &l)
	}

	return nil
}

func (z *leaf) String() (s string) {
	return fmt.Sprintf("mode: %06s file: %s hash: %x\n",
		z.mode, z.file, z.hash)
}
func (z *tree) String() (s string) {
	for _, v := range z.leaf {
		s += v.String()
	}
	return s
}

//wi	re9 hdr name[,toSpace] value[,toNull]
//wi	re9 kv name[,toSpace] value[,toLine]
//wi	re9 id name[,toSpace] email[,toSpace] date[,toLine]
//wi	re9 commit hdr[,kv] tree[,kv] parent[,kv] author[,id] commiter[,id] comment[,toLine]
