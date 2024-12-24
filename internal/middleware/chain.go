package middleware

import (
	"net/http"
)

type Constructor = func(http.Handler) http.Handler

type Chain struct {
	middlewares []Constructor
}

func NewChain(middlewares ...Constructor) Chain {
	return Chain{
		middlewares: middlewares,
	}
}

func (c Chain) Then(fn http.HandlerFunc) http.Handler {
	var ret http.Handler = fn
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		ret = c.middlewares[i](ret)
	}
	return ret
}
