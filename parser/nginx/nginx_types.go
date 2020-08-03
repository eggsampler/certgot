//go:generate pigeon -o nginx.go nginx.peg

package nginx

type SimpleDirective struct {
	Name       string
	Parameters []string
}

type CommentDirective string

type BlockDirective struct {
	Name       string
	Parameters []string
	Children   []interface{}
}

// TODO: stringer interface which also indents?
