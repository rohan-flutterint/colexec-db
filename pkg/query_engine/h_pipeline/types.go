package pipeline

import relalgebra "colexecdb/pkg/query_engine/i_rel_algebra"

type Pipeline struct {
	// attrs, column list.
	attrs []string
	// orders to be executed
	instructions relalgebra.Instructions
}
