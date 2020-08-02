//go:generate pigeon -o nginx.go nginx.peg

package nginx

type Directive struct {
	Name          string
	HasParameter  bool
	Parameter     string
	HasDirectives bool
	Directives    []Directive
}

// TODO: stringer interface which also indents?
