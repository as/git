package git

import (
	"bytes"
	"io"
	"strconv"
	"time"
)

type ID struct {
	Name  string
	Email string
	Date  time.Time
}

type Commit struct {
	Tree      *Tree
	Parent    *Commit
	Author    ID
	Committer ID
	Comment   string
}

func ParseCommit(r io.Reader) (c *Commit, err error) {
	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(r); err != nil {
		return nil, err
	}
	var asize string
	var size int

	_, asize, err = ParseKeyValue(buf, ' ', "commit", '\x00')
	if err != nil {
		return nil, err
	}
	if size, err = strconv.Atoi(string(asize)); err != nil {
		return nil, err
	}
	size = size
	return nil, nil
}

func ParseKeyValue(buf *bytes.Buffer, delim byte, key string, delim2 byte) (k string, v string, err error) {
	if k, err = cut(buf, delim); err != nil || k != key {
		return "", "", err
	}
	if v, err = cut(buf, delim); err != nil {
		return "", "", err
	}
	return k, v, err
}
