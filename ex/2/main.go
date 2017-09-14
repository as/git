// Example git driver
package main

import (
	"crypto/sha1"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/as/git"
	"github.com/as/mute"
)

const (
	Prefix     = "git: "
	BufferSize = 1024 * 1024
	Debug      = false // true false
)

var arg struct {
	passphrase string
	sshkeypath string
	s, u, c    string
	n, v       bool
}

var f *flag.FlagSet

var (
	authfmt = "proto=pass service=ssh server=%s user=%s"
)

func init() {
	f = flag.NewFlagSet("main", flag.ContinueOnError)

	f.StringVar(&arg.passphrase, "kp", "", "")
	f.StringVar(&arg.sshkeypath, "k", "", "")
	f.StringVar(&arg.s, "s", os.Getenv("ssh"), "")
	f.StringVar(&arg.u, "u", os.Getenv("user"), "")
	f.StringVar(&arg.c, "c", "", "")
	f.BoolVar(&arg.n, "n", false, "")
	f.BoolVar(&arg.v, "v", false, "")

	err := mute.Parse(f, os.Args[1:])
	if err != nil {
		printerr(err)
		os.Exit(1)
	}
}

func clean(s string) string {
	return strings.TrimSpace(s)
}

var N int32

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {

	//
	// Walk the tree and print objects found within
	/*
		fmt.Println("TREE")
		for _, v := range t.Leaves() {
			x := fmt.Sprintf("%x", v.hash)
			file := fmt.Sprintf("testdata/.git/objects/%02s/%38s", x[:2], x[2:])
			data, err := ioutil.ReadFile(file)
			no(err)
			fmt.Printf("\n\n#####\nname: %s hash: %s content:\n\n%q\n", v.file, x, string(data))
		}
	*/

	Git := clientinit()
	//c, err := Git.Command("git-upload-pack", "garethjensen/goconquer")
	c, err := Git.Command("git-upload-pack", "as/structslice")
	pl, err := c.UploadPackTX()
	no(err)
	pl.Print()

	pack, err := c.UploadPackRX()
	no(err)
	for _, v := range pack.Objects() {
		data := []byte(v.String())
		x := fmt.Sprintf("%x", sha1.Sum(data))
		dir := fmt.Sprintf("testdata/.git/objects/%02s", x[:2])
		file := fmt.Sprintf("%38s", x[2:])
		os.MkdirAll(dir, 0770)
		ioutil.WriteFile(fmt.Sprintf("%s/%s", dir, file), data, 0666)
		fmt.Printf("computed: %s hashref: %s\n", x, pl.Refmap[x])
	}

	fmt.Println("REFMAP")
	for k, v := range pl.Refmap {
		if v.Ref == "refs/heads/master" {
			fmt.Printf("testdata/.git/refs/heads/master key=%s value=%v\n", k, v)
		}
	}

	//t, err := git.TreeFromFile("../data/.git/objects/76/20e8bce130075206e51cdfb70a8b915731e3dd")
	t, err := git.TreeFromFile("testdata/.git/objects/db/af7f6df43c70f1452755f3ce9a4d86e3c14e6f") //as/structslice
	no(err)
	fmt.Println(t)

	fmt.Println("TREE")
	for _, v := range t.Leaves() {
		o, err := v.Eval()
		fmt.Println(err)
		switch t := o.(type) {
		case *git.Blob:
			fmt.Printf("name: %s\n", t.Name())
			fmt.Printf("hash: %s\n", t.Hash())
			n := fmt.Sprintf("../data/%s", t.Name())
			ioutil.WriteFile(n, t.Data(), 0666)
			fmt.Println("wrote file:", n)
			fmt.Println()
		default:
			fmt.Printf("%T\n", t)
		}
	}

}

func clientinit() git.Git {
	addr := dialsplit(arg.s)
	if arg.s == "" || addr == nil {
		printerr("set $ssh")
		os.Exit(1)
	}
	conf := &ssh.ClientConfig{
		User: "git",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeysCallback(GetSigners),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	addrsvc := addr.addr + ":" + addr.svc
	client, err := ssh.Dial(addr.net, addrsvc, conf)
	no(err)
	return git.Git{client}

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

type Addr struct {
	net, addr, svc string
}

func dialsplit(dial string) *Addr {
	s := strings.Split(dial, "!")
	a := new(Addr)
	n := len(s)
	switch {
	case n == 3:
		a.svc = s[2]
		fallthrough
	case n == 2:
		a.addr = s[1]
		fallthrough
	case n == 1:
		a.net = s[0]
		return a
	}
	printerr("bad dialer:", dial)
	return nil
}
func GetSigners() ([]ssh.Signer, error) {
	signers := make([]ssh.Signer, 0)
	buf, err := ioutil.ReadFile(arg.sshkeypath)
	no(err)

	xxx, err := ssh.ParsePrivateKeyWithPassphrase(buf, []byte(arg.passphrase))
	no(err)

	signers = append(signers, xxx)

	return signers, nil
}
