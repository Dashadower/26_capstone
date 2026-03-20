package traceinspector

type node_types string

const (
	node_basic node_types = "basic"
	node_cond  node_types = "cond"
)

type CFGNode struct {
	id        int
	code      string
	node_type node_types
	line_num  int
}

type CFGEdge struct {
	id           int
	from_node_id int
	to_node_id   int
	label        string
}
