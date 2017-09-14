package git

import (
	"fmt"
	"io"

	"golang.org/x/crypto/ssh"
)

type Git struct {
	*ssh.Client
}

type Cmd struct {
	Name    string
	Args    []string
	in      io.WriteCloser
	out     io.Reader
	Session *ssh.Session
}

func (g *Git) Command(name string, args ...string) (*Cmd, error) {
	switch name {
	case "git-upload-pack":
		if n := len(args) != 1; n {
			return nil, fmt.Errorf("git-upload-pack: expecting 1 arg, got %d", len(args))
		}
		sess, err := g.NewSession()
		if err != nil {
			return nil, err
		}
		cmd := &Cmd{Name: name, Args: args, Session: sess}
		if cmd.in, err = sess.StdinPipe(); err != nil {
			return nil, err
		}
		if cmd.out, err = sess.StdoutPipe(); err != nil {
			return nil, err
		}
		//sess.Stderr = os.Stderr
		return cmd, nil
	}
	return nil, fmt.Errorf("bad git command: %s", name)
}
