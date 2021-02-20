package home

import (
	"bytes"
	"io/fs"
)

var nopFS fs.FS = &nopFSImpl{}

type nopFSImpl struct{}

func (f *nopFSImpl) Open(path string) (fs.File, error) {
	return nil, &fs.PathError{
		Op:   "open",
		Path: path,
		Err:  fs.ErrNotExist,
	}
}

func checkIsDir(targetFS fs.FS, targetPath string) (bool, error) {
	info, err := fs.Stat(targetFS, targetPath)
	if err != nil {
		return false, err
	}

	return info.IsDir(), nil
}

type compositeError struct {
	Text   string
	Errors []error
}

func (c *compositeError) Error() string {
	buf := bytes.NewBufferString(c.Text)

	for _, err := range c.Errors {
		buf.WriteString("\n* ")
		buf.WriteString(err.Error())
	}

	return buf.String()
}
