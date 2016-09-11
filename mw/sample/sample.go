package sample

import (
	"context" // go 1.7 here
	"log"
	"net/http"

	"github.com/LeKovr/elsa/struct/pagesample"
)

// -----------------------------------------------------------------------------

// Flags defines local application flags
type Flags struct {
	Title string `long:"sample_title" default:"Sample page" description:"Sample page title"`
}

// -----------------------------------------------------------------------------

type Info struct{ pagesample.Attr }

func (i Info) Ext() string {
	return ""
}
func (i Info) Template() string {
	return ""
}
func (i Info) Data() interface{} { //} *stats_lib.Data {
	return i.Attr
}

// -----------------------------------------------------------------------------

// Middleware is a struct that has a ServeHTTP method
type Middleware struct {
	Log      *log.Logger
	Config   *Flags
	writable bool
}

// -----------------------------------------------------------------------------

func New(logger *log.Logger, cfg *Flags, writable bool) *Middleware {
	return &Middleware{Log: logger, Config: cfg, writable: writable}
}

// -----------------------------------------------------------------------------

// ServeHTTP is the middleware handler
func (mw *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// Log moose status
	mw.Log.Printf("debug: APP: %v", mw.writable)

	ctx := r.Context()

	if mw.writable {
		// Save data in context
		p := Info{pagesample.Attr{Title: mw.Config.Title, Created: "just now", Body: "Hello here"}}
		ctx = context.WithValue(ctx, "data", &p)
		r = r.WithContext(ctx)
	} else {
		// Read data from context
		d := ctx.Value("data")
		p := d.(*Info).Data()
		mw.Log.Printf("debug: Got from context: %v", p.(pagesample.Attr).Title)
	}
	next(w, r)
}
