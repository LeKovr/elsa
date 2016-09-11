package flow

import (
	"context"
	"log"
	"net/http"

	"github.com/urfave/negroni"
)

const Key string = "flowStatus"

// -----------------------------------------------------------------------------

type Middleware struct {
	Log *log.Logger
}

func New(logger *log.Logger) *Middleware {
	return &Middleware{Log: logger}
}

// ServeHTTP is the middleware handler
func (mw *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ctx := r.Context()
	var data int
	mw.Log.Printf("debug: flow init: %v", data)
	ctx = context.WithValue(ctx, Key, &data)
	r = r.WithContext(ctx)
	next(w, r)
}

// -----------------------------------------------------------------------------

type MiddlewareHandler struct {
	Log     *log.Logger
	handler negroni.Handler
}

func NewHandler(logger *log.Logger, handler negroni.Handler) *MiddlewareHandler {
	return &MiddlewareHandler{Log: logger, handler: handler}
}

// ServeHTTP is the middleware handler
func (mwh *MiddlewareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if !Prohibited(r) {
		mwh.handler.ServeHTTP(w, r, next)
	} else {
		next(w, r)
	}
}

// -----------------------------------------------------------------------------

func Finish(r *http.Request) {
	ctx := r.Context()
	data := ctx.Value(Key)
	flag := data.(*int)
	*flag = 1
}

func Prohibit(r *http.Request) {
	ctx := r.Context()
	data := ctx.Value(Key)
	flag := data.(*int)
	*flag = 2
}

func FinishHandler(logger *log.Logger) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		Finish(r)
		next(w, r)
	}
}

func Finished(r *http.Request) bool {
	ctx := r.Context()
	data := ctx.Value(Key)
	flag := data.(*int)
	return *flag > 0
}

func Prohibited(r *http.Request) bool {
	ctx := r.Context()
	data := ctx.Value(Key)
	flag := data.(*int)
	return *flag > 1
}
