package examples

import (
	"io/fs"
	"time"
)

type SimpleFS struct {
	files map[string]*file
}

type file struct {
	name    string
	data    []byte
	modTime time.Time
	isDir   bool
}

func NewSimpleFS() *SimpleFS {
	return &SimpleFS{
		files: make(map[string]*file),
	}
}

func (sfs *SimpleFS) AddFile(name string, data []byte) {
	sfs.files[name] = &file{
		name:    name,
		data:    data,
		modTime: time.Now(),
		isDir:   false,
	}
}

func (sfs *SimpleFS) Open(name string) (fs.File, error) {
	f, ok := sfs.files[name]
	if !ok {
		return nil, fs.ErrNotExist
	}
	return &openFile{file: f, offset: 0}, nil
}

type openFile struct {
	file   *file
	offset int
}

func (of *openFile) Stat() (fs.FileInfo, error) {
	return of.file, nil
}

func (of *openFile) Read(b []byte) (int, error) {
	if of.offset >= len(of.file.data) {
		return 0, nil
	}
	n := copy(b, of.file.data[of.offset:])
	of.offset += n
	return n, nil
}

func (of *openFile) Close() error {
	return nil
}

func (f *file) Name() string       { return f.name }
func (f *file) Size() int64        { return int64(len(f.data)) }
func (f *file) Mode() fs.FileMode  { return 0444 }
func (f *file) ModTime() time.Time { return f.modTime }
func (f *file) IsDir() bool        { return f.isDir }
func (f *file) Sys() interface{}   { return nil }

func ExampleFS() fs.FS {
	sfs := NewSimpleFS()
	sfs.AddFile("hello.txt", []byte("Hello, World!"))
	sfs.AddFile("data.json", []byte(`{"key": "value"}`))
	return sfs
}
