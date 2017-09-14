package fs

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
)

var (
	NameFunc = keyPath
	HashFunc = sha1.New
)

func keyPath(k Key) string {
	if len(k) < 2 {
		return ""
	}
	return filepath.Clean(fmt.Sprintf("%s/%s", k[:2], k[2:]))
}

// Fs is a key value store
type Fs struct {
	// NameFunc maps Keys to file names on disk. The names
	// are used to Load and Store data on the filesystem.
	NameFunc func(Key) string
	// HashFunc returns a new hash.Hash algorithm used to address
	// the underlying file contents
	HashFunc func() hash.Hash
	dir      string
}

// New creates a filesystem in path with a SHA1
// key value store.
func New(path string) *Fs {
	return &Fs{
		dir:      filepath.Clean(path),
		NameFunc: NameFunc,
		HashFunc: HashFunc,
	}
}

func (f *Fs) Store() (*File, error) {
	if f.dir == "" {
		return nil, fmt.Errorf("path not set")
	}
	return NewFile(f.dir, f.NameFunc), nil
}

// NewFile returns a *File satisfying io.WriteCloser. A File's data
// may not be committed to the filesystem until a Sync() or Close()
// occur.
func (f *Fs) NewFile() (*File, error) {
	return f.Store()
}

// Open resolves k using Fs.NameFunc and returns a io.ReadCloser
// of the underlying file contents.
func (f *Fs) Open(k Key) (r io.ReadCloser, err error) {
	return f.Load(k)
}

func (f *Fs) Load(k Key) (r io.ReadCloser, err error) {
	fd, err := os.Open(filepath.Join(f.dir, f.NameFunc(k)))
	if err != nil {
		return nil, err
	}
	zr, err := zlib.NewReader(fd)
	if err != nil {
		return nil, err
	}
	return zr, err
}

type Key string

type File struct {
	dir string
	buf *bytes.Buffer
	h   hash.Hash
	io.Writer

	namefn func(Key) string
}

func NewFile(path string, namefn func(Key) string) *File {
	f := &File{
		dir:    filepath.Clean(path),
		buf:    new(bytes.Buffer),
		h:      sha1.New(),
		namefn: namefn,
	}
	f.Writer = io.MultiWriter(f.buf, f.h)
	return f
}

func (f *File) Write(p []byte) (n int, err error) {
	n, err = f.Writer.Write(p)
	return
}

// Sync commits the File's contents to the filesystem
func (f *File) Sync() error {
	return f.Flush()
}

func (f *File) Flush() error {
	k, err := f.Key()
	if err != nil {
		return err
	}
	if _, err := os.Stat(f.dir); err != nil {
		err = os.MkdirAll(f.dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	path := filepath.Join(f.dir, f.namefn(k))
	if _, err := os.Stat(path); err == nil {
		// The file already exists, there's no point in overwriting it
		// with the same contents
		return nil
	}

	fd, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fd.Close()
	zfd, err := zlib.NewWriterLevel(fd, zlib.BestCompression)
	if err != nil {
		return err
	}
	_, err = f.buf.WriteTo(zfd)
	zfd.Close()
	return err
}

// Close calls Sync to flush the file to the filesystem
func (f *File) Close() error {
	return f.Flush()
}

// Key returns the file's key. It is an error to call
// this method before calling Sync() or Close().
func (f *File) Key() (Key, error) {
	return Key(fmt.Sprintf("%40x", f.h.Sum(nil))), nil
}
