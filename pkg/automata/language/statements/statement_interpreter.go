package statements

import (
	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/automata/language/expressions"
)

type SymbolicInterpreter struct {
	context     *z3.Context
	valuations  expressions.Valuations
	expressions expressions.SymbolicInterpreter
}

func NewSymbolicInterpreter(
	context *z3.Context, valuations expressions.Valuations,
) SymbolicInterpreter {
	return SymbolicInterpreter{
		context:    context,
		valuations: valuations,
		expressions: expressions.NewSymbolicInterpreter(
			context, valuations,
		),
	}
}

func (interpreter SymbolicInterpreter) Interpret(statement Statement) {
	switch cast := any(statement).(type) {
	case Block:
		interpreter.Block(cast)
	case Assignment:
		interpreter.Assignment(cast)
	}
}

func (interpreter SymbolicInterpreter) Block(block Block) {
	for idx := range block.statements {
		interpreter.Interpret(block.statements[idx])
	}
}

func (interpreter SymbolicInterpreter) Assignment(assignment Assignment) {
	interpreter.valuations.Assign(assignment.lhs.Symbol(), assignment.rhs)
}