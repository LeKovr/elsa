// page holds web page attributes
package page

type Attr interface {
	Ext() string
	Template() string
	Data() interface{}
}
