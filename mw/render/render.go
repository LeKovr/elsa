// render is a elsa middleware for render
package render

import (
	"log"
	"net/http"
	"strings"

	"gopkg.in/unrolled/render.v1"

	"github.com/LeKovr/elsa/mw/flow"
	"github.com/LeKovr/elsa/struct/page"
	"github.com/LeKovr/go-base/dirtree"
)

// -----------------------------------------------------------------------------

// Flags defines local application flags
type Flags struct {
	Prefix    string `long:"render_refix"  default:"/"             description:"URL prefix to handle"`
	Ext       string `long:"render_ext"    default:".tmpl"         description:"Template extention"`
	Templates string `long:"render_path"   default:"templates/render" description:"Dir where templates are"`
}

// -----------------------------------------------------------------------------

// Middleware is a struct that has a ServeHTTP method
type Middleware struct {
	Log    *log.Logger
	Config *Flags
	Tree   *dirtree.Tree
	Engine *render.Render
}

// -----------------------------------------------------------------------------

func New(logger *log.Logger, cfg *Flags, show404 bool) *Middleware {
	tree, _ := dirtree.New(logger, cfg.Templates, cfg.Ext)

	ren := render.New(render.Options{
		IndentJSON: true,
		Layout:     "layout",
		Directory:  cfg.Templates,
	})

	return &Middleware{Log: logger, Config: cfg, Tree: tree, Engine: ren}
}

// -----------------------------------------------------------------------------

// The middleware handler
func (mw *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	if flow.Finished(r) {
		mw.Log.Printf("debug: Flow is finished")
		next(w, r)
		return
	}

	ctx := r.Context()
	data := ctx.Value("data")
	var attr page.Attr
	var slug string

	if data != nil {
		mw.Log.Printf("debug: Got data: %+v", data)

		attr = data.(page.Attr)    //(data)
		if attr.Ext() == ".json" { // we serve .json
			mw.Log.Printf("debug: Show data as json: %+v", data)
			mw.Engine.JSON(w, http.StatusOK, attr.Data())
			flow.Finish(r)
			next(w, r)
			return
		}
		if attr.Ext() == mw.Config.Ext { // we serve this ext
			slug = attr.Template()
		} else if attr.Ext() != "" { // other engine will serve
			next(w, r)
			return
		}
		// Ext() == "" => any engine may serve, check request uri
	}
	// data == nil => check request uri

	if slug == "" { // check request uri
		url := r.URL.Path
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

	if ok {
		tmpl := strings.TrimPrefix(strings.TrimSuffix(node.File.Path, mw.Config.Ext), "/")
		mw.Log.Printf("debug: Use template: %s", tmpl)
		if err := mw.Engine.HTML(w, http.StatusOK, tmpl, attr); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		flow.Finish(r)
		return
	}
	next(w, r)
}
