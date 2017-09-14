package git

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
)

func (c *Cmd) Start() error {
	comb := strings.Join(append([]string{c.Name}, c.Args...), " ")
	return c.Session.Start(comb)
}

func (c *Cmd) UploadPackTX() (*PackList, error) {
	pl := &PackList{}
	err := c.Start()
	if err != nil {
		return nil, err
	}

	if err = pl.Read(c.out); err != nil {
		return nil, err
	}
	N := rand.Int31n(int32(len(pl.order)))

	fmt.Printf("randomly selected: %s", pl.order[N])

	opt := Options{"multi-ack", "side-band-64k", "ofs-delta"}
	for _, v := range []string{
		opt.Want(pl.order[N].Hash),
		//opt.Want("b0b57e6a7848ddfc2bd9e3de6e7fc458cd6d2dce"),
		opt.Flush(),
		opt.Done(),
	} {
		fmt.Fprint(c.in, v)
	}

	lr, err := nextMsgReader(c.out)
	io.Copy(os.Stdout, lr)
	return pl, err
}

func (c *Cmd) UploadPackRX() (ph *packhdr, err error) {
	t := make([]byte, 1)
	lr, err := nextMsgReader(c.out)
	pr, pw := io.Pipe()
	go func() {
		io.Copy(os.Stdout, lr)
		defer pw.Close()
		for {
			lr, err = nextMsgReader(c.out)
			_, err = lr.Read(t)
			if err != nil {
				return
			}
			switch t[0] {
			case 0x0:
				fmt.Println("channel 0")
				io.Copy(os.Stdout, lr)
			case 0x1:
				fmt.Println("channel 1")
				_, err = io.Copy(pw, lr)
			case 0x2:
				fmt.Println("channel 2")
				io.Copy(os.Stdout, lr)
			default:
				printerr(fmt.Errorf("bad channel: %d", t[0]))
			}

		}
	}()
	if err != nil {
		return
	}
	ph = &packhdr{}
	return ph, ph.ReadBinary(bufio.NewReader(pr))
}
