package embedfs

import (
	"errors"
	"io/fs"
	"path"
	"path/filepath"
	"strings"
)

// Dir 路径
type Dir struct {
	fs.FS
	Dir string
}

// Open 实现接口方法
func (d Dir) Open(name string) (fs.File, error) {
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) {
		return nil, errors.New("http: invalid character in file path")
	}
	return d.FS.Open(filepath.Join(d.Dir, filepath.FromSlash(path.Clean("/"+name))))
}
