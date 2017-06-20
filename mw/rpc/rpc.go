package rpc

import (
	"log"
	"net/http"
	"strings"

	rpc "github.com/gorilla/rpc/v2"
	json "github.com/gorilla/rpc/v2/json2"
)

// -----------------------------------------------------------------------------

// Middleware is a struct that has a ServeHTTP method
type Middleware struct {
	Log    *log.Logger
	Server *rpc.Server
	Prefix string
}

// -----------------------------------------------------------------------------

func New(logger *log.Logger, prefix string, service interface{}) *Middleware {

	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(service, "")

	return &Middleware{Log: logger, Prefix: prefix, Server: s}
}

// -----------------------------------------------------------------------------
func (mw *Middleware) Add(name string, service interface{}) {
	mw.Server.RegisterService(service, name)
}

// -----------------------------------------------------------------------------

// ServeHTTP is the middleware handler
func (mw *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	url := r.URL.RequestURI()
	inGame := strings.HasPrefix(url, mw.Prefix)

	if inGame {
		mw.Log.Printf("debug: Requested rpc for: %s / %s", url, r.Header.Get("Content-Type"))
		mw.Server.ServeHTTP(w, r)

		// TODO: prepare json for next middleware?
		return
	}
	next(w, r)
}
