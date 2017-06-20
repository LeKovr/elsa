// pagedata holds web page attributes
package pagedata

// import "html/template"

type Simple struct {
	Meta map[string]interface{}
}

func (pd *Simple) Set(key string, value interface{}) string {
	pd.Meta[key] = value
	return ""
}

func (pd *Simple) Get(key string) interface{} {
	return pd.Meta[key]
}

func New() *Simple {
	s := Simple{Meta: map[string]interface{}{}}
	return &s
}
