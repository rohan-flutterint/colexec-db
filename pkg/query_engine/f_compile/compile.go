package compile

import (
	batch "colexecdb/pkg/query_engine/b_batch"
	process "colexecdb/pkg/query_engine/c_process"
	"colexecdb/pkg/query_engine/d_parser"
	planner "colexecdb/pkg/query_engine/e_planner"
	scope "colexecdb/pkg/query_engine/g_scope"
	relalgebra "colexecdb/pkg/query_engine/i_rel_algebra"
	"colexecdb/pkg/query_engine/i_rel_algebra/projection"
	"colexecdb/pkg/storage_engine"
	"context"
	"errors"
)

// New is used to new an object of compile
func New(sql string, ctx context.Context, proc *process.Process, stmt parser.Statement) *Compile {
	c := &Compile{}
	c.Ctx = ctx
	c.sql = sql
	c.Process = proc
	c.stmt = stmt
	return c
}

// Compile is the entrance of the compute-execute-layer.
// It generates a scope (logic pipeline) for a query plan.
func (c *Compile) Compile(ctx context.Context, pn planner.Plan, fill func(any, *batch.Batch) error) (err error) {

	c.Ctx = c.Process.Ctx
	c.fill = fill
	c.pn = pn

	c.scope, err = c.compileScope(ctx, pn)
	return nil
}

func (c *Compile) compileScope(ctx context.Context, pn planner.Plan) ([]*scope.Scope, error) {
	switch qry := pn.(type) {
	case *planner.QueryPlan:
		rs := scope.Scope{
			Magic:        Normal,
			Plan:         pn,
			Instructions: make(relalgebra.Instructions, 0),
		}
		rs.Instructions = append(rs.Instructions, relalgebra.Instruction{
			Op: relalgebra.Projection,
			Arg: &projection.Argument{
				Es: qry.Params,
			},
		})
		rs.DataSource = &Source{
			Reader:     storage_engine.NewMergeReader(),
			Attributes: []string{"Id", "Age"},
		}

		rs.Process = c.Process

		return []*scope.Scope{&rs}, nil

	case *planner.DDLPlan:
		switch qry.Type {
		case planner.DdlCreateTable:
			return []*scope.Scope{{
				Magic: CreateTable,
				Plan:  pn,
			}}, nil
		}
	}
	return nil, errors.New("unimplemented")
}
