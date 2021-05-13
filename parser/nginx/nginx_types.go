//go:generate pigeon -o nginx.go nginx.peg

package nginx

type Directive struct {
	Name       string
	Parameters []string
	Comment    bool
	Blank      bool
	Children   []Directive
}

func toString(v interface{}) string {
	s, _ := v.(string)
	return s
}

func toStringSlice(v interface{}) []string {
	var ss []string
	for _, vv := range toIfaceSlice(v) {
		s := toString(vv)
		if s == "" {
			continue
		}
		ss = append(ss, s)
	}
	return ss
}

func toIfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	vv, _ := v.([]interface{})
	return vv
}

func toDirectiveSlice(v interface{}) []Directive {
	if v == nil {
		return nil
	}
	vv, _ := v.([]interface{})
	var d []Directive
	for _, vvv := range vv {
		vvvv, ok := vvv.(Directive)
		if ok {
			d = append(d, vvvv)
		}
	}
	return d
}
