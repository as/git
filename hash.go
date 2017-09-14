package git

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
)

var (
	ErrHashLen = errors.New("hash not 20 bytes")
)

type Hash string

func (h Hash) String() string {
	return string(h)
}

func (h Hash) Bytes() ([]byte, error) {
	return hex.DecodeString(string(h))
}

func HashData(p []byte) (Hash, error) {
	return ParseHash(sha1.Sum(p)), nil
}

func ParseHashString(h string) (Hash, error) {
	n := len(h)
	if n != 40 {
		return "", ErrHashLen
	}
	return Hash(h), nil
}

func ParseHash(p [20]byte) Hash {
	return Hash(fmt.Sprintf("%40x", p[:]))
}
