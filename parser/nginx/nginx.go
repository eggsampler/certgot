package nginx

type SimpleDirective struct {
	Name       string
	Parameters []string
}

type BlockDirective struct {
	Name             string
	Parameters       []string
	SimpleDirectives map[string]SimpleDirective
	BlockDirectives  map[string]BlockDirective
}
