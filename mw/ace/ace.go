// ace is a elsa middleware for ace
package ace

import (
	"log"
	"net/http"
	"strings"

	"github.com/yosssi/ace"

	"github.com/LeKovr/elsa/mw/flow"
	"github.com/LeKovr/elsa/pagedata"
	"github.com/LeKovr/elsa/struct/page"
	"github.com/LeKovr/go-base/dirtree"
)

// -----------------------------------------------------------------------------

// Flags defines local application flags
type Flags struct {
	Prefix    string `long:"ace_prefix" default:"/"             description:"URL prefix to handle"`
	Ext       string `long:"ace_ext"    default:".ace"          description:"Template extention"`
	Templates string `long:"ace_path"   default:"templates/ace" description:"Dir where templates are"`
	Error403  string `long:"ace_403"    default:".error403"     description:"Template for 403 error"`
	Error404  string `long:"ace_404"    default:".error404"     description:"Template for 404 error"`
}

// -----------------------------------------------------------------------------

// Middleware is a struct that has a ServeHTTP method
type Middleware struct {
	Log    *log.Logger
	Config *Flags
	Tree   *dirtree.Tree
	IsLast bool // middleware is the last in chain, show errors and show 404 if no template
}

// -----------------------------------------------------------------------------

func New(logger *log.Logger, cfg *Flags, show404 bool) *Middleware {
	tree, _ := dirtree.New(logger, cfg.Templates, cfg.Ext)
	return &Middleware{Log: logger, Config: cfg, Tree: tree, IsLast: show404}
}

// -----------------------------------------------------------------------------

// The middleware handler
func (mw *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	var slug, tmplPath string
	var attr page.Attr

	if flow.Prohibited(r) {
		if mw.IsLast {
			mw.Log.Print("debug: Show 403")
			w.WriteHeader(http.StatusForbidden)
			tmplPath = mw.Config.Error403
		} else {
			next(w, r)
			return
		}
	} else if flow.Finished(r) {
		mw.Log.Print("debug: Flow is finished")
		next(w, r)
		return
	} else {
		// Regular request
		ctx := r.Context()
		data := ctx.Value("data")
		if data != nil {
			mw.Log.Print("debug: Show data from context")
			attr = data.(page.Attr)          //(data)
			if attr.Ext() == mw.Config.Ext { // we serve this ext
				slug = attr.Template()
			} else if attr.Ext() != "" { // other engine will serve
				next(w, r)
				return
			}
			// Ext() == "" => any engine may serve, check request uri
		}
		// data == nil => check request uri
		url := r.URL.Path
		if slug == "" { // check request uri
			inGame := strings.HasPrefix(url, mw.Config.Prefix)
			if !inGame {
				next(w, r)
				return
			} else {
				slug = strings.TrimPrefix(url, mw.Config.Prefix)
			}
		}
		mw.Log.Printf("debug: Requested page: %s", slug)
		node, ok := mw.Tree.Node(slug, nil)
		if !ok {
			if mw.IsLast {
				mw.Log.Print("debug: Show 404")
				w.WriteHeader(http.StatusNotFound)
				tmplPath = mw.Config.Error404
			} else {
				next(w, r)
				return
			}
		} else {
			tmplPath = node.File.Path
		}
	}

	tmpl := strings.TrimSuffix(tmplPath, mw.Config.Ext)
	mw.Log.Printf("debug: Use template: %s", tmpl)
	tpl, err := ace.Load("layout", tmpl, &ace.Options{BaseDir: mw.Config.Templates})

	mw.Log.Printf("debug: TMPL:%+v", tpl.Tree.Root)
	var pageData interface{}
	if attr != nil {
		pageData = attr.Data
	} else {
		pageData = pagedata.New()
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else if err := tpl.Execute(w, pageData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	flow.Finish(r)
	next(w, r)
}
