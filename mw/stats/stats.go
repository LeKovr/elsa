package stats

import (
	"context" // go 1.7 here
	"log"
	"net/http"
	"strings"

	stats_lib "github.com/thoas/stats"
)

// -----------------------------------------------------------------------------

// Flags defines local application flags
type Flags struct {
	Prefix string `long:"stats_prefix" default:"/stats"             description:"URL prefix to handle"`
}

// -----------------------------------------------------------------------------

// Middleware is a struct that has a ServeHTTP method
type Middleware struct {
	Log    *log.Logger
	Config *Flags
	Engine *stats_lib.Stats
}

// -----------------------------------------------------------------------------

type Info struct{ engine *stats_lib.Stats }

func (i Info) Ext() string {
	return ".json"
}
func (i Info) Template() string {
	return ""
}
func (i Info) Data() interface{} { //} *stats_lib.Data {
	return i.engine.Data()
}

// -----------------------------------------------------------------------------

func New(logger *log.Logger, cfg *Flags) *Middleware {
	s := stats_lib.New()
	return &Middleware{Log: logger, Config: cfg, Engine: s}
}

// -----------------------------------------------------------------------------

// The middleware handler
func (mw *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	url := r.URL.RequestURI()
	inGame := strings.HasPrefix(url, mw.Config.Prefix)

	if inGame {
		mw.Log.Printf("debug: Generate stats for: %s", url)
		ctx := r.Context()
		data := Info{engine: mw.Engine}
		ctx = context.WithValue(ctx, "data", &data)
		r = r.WithContext(ctx)
	}
	next(w, r)
}
