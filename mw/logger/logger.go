package logger

import (
	"log"
	"net/http"
	"time"

	"github.com/urfave/negroni"
)

// -----------------------------------------------------------------------------

// Middleware is a struct that has a ServeHTTP method
type Middleware struct {
	Log     *log.Logger
	IPField string // `long:"logger_realip_field" default:"real-ip" description:"Context field for Real ip"`
}

// -----------------------------------------------------------------------------

func New(logger *log.Logger, field string) *Middleware {
	return &Middleware{Log: logger, IPField: field}
}

// -----------------------------------------------------------------------------

// The middleware handler
func (mw *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	start := time.Now()

	ctx := r.Context()
	d := ctx.Value(mw.IPField)
	ip := d.(*string)
	mw.Log.Printf("info: Started handling request: %s %s %s", *ip, r.Method, r.RequestURI)
	next(w, r)

	latency := time.Since(start)
	res := w.(negroni.ResponseWriter)

	// latency.Nanoseconds(),

	mw.Log.Printf("info: Completed handling request: %s %s %s %d %s", *ip, r.Method, r.RequestURI, res.Status(), latency)

}
