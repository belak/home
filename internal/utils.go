package internal

import (
	"errors"
	"net/url"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

var ErrUnimplemented = errors.New("unimplemented")

func GetFunctionName(i interface{}) string {
	// Get the full value of the function
	ret := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()

	return cleanFuncName(ret)
}

func CurrentFunc(skip int) string {
	counter, _, _, success := runtime.Caller(skip)

	if !success {
		return "unknown"
	}

	return cleanFuncName(runtime.FuncForPC(counter).Name())
}

func cleanFuncName(in string) string {
	// Use basename to strip off the path of the package
	ret := filepath.Base(in)

	// Strip off the -fm which seems to be used for receivers
	ret = strings.TrimSuffix(ret, "-fm")

	return ret
}

func CleanURL(in string) string {
	u, err := url.Parse(in)
	if err != nil {
		return "/"
	}

	if u.Host != "" {
		return "/"
	}

	if u.Path == "" {
		return "/"
	}

	return u.Path
}
