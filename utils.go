package home

import "io/fs"

var nopFS fs.FS = &nopFSImpl{}

type nopFSImpl struct{}

func (f *nopFSImpl) Open(path string) (fs.File, error) {
	return nil, &fs.PathError{
		Op:   "open",
		Path: path,
		Err:  fs.ErrNotExist,
	}
}
