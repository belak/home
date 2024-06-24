package internal

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/monoculum/formam/v3"
)

var formDecoder *formam.Decoder = formam.NewDecoder(&formam.DecoderOptions{
	TagName: "form",
})

func Bind(r *http.Request, v interface{}) ([]string, bool) {
	var ret []string

	_ = r.ParseForm()

	err := formDecoder.Decode(r.PostForm, v)
	if err != nil {
		return []string{err.Error()}, false
	}

	err = validator.New().Struct(v)
	if err == nil {
		return nil, true
	}

	switch iErr := err.(type) {
	case validator.ValidationErrors:
		for _, item := range iErr {
			ret = append(ret, item.Error())
		}
	default:
		ret = append(ret, err.Error())
	}

	return ret, false
}
