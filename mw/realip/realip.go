package realip

import (
	"context" // go 1.7 here
	"log"
	"net"
	"net/http"
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

	ip, src := userIP(r)
	mw.Log.Printf("trace: Fetched ip %s from %s", ip, src)
	ctx := r.Context()
	ctx = context.WithValue(ctx, mw.IPField, &ip)
	r = r.WithContext(ctx)
	next(w, r)
}

func userIP(r *http.Request) (ip string, ipSource string) {
	ip = r.Header.Get("X-Real-Ip")
	if ip != "" {
		return ip, "real-ip"
	}
	ip = r.Header.Get("X-Forwarded-For")
	if ip != "" {
		return ip, "fwd-for"
	}
	ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	return ip, "rem-addr"
}
