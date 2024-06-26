package automata

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/automata/language"
	"github.com/Brandhoej/gobion/pkg/symbols"
	"github.com/stretchr/testify/assert"
)

func Test_Name(t *testing.T) {
	// Arrange
	context := z3.NewContext(z3.NewConfig())

	symbols := symbols.NewSymbolsMap[string](
		symbols.NewSymbolsFactory(),
	)
	x, y := symbols.Insert("x"), symbols.Insert("y")
	variables := language.NewVariablesMap()
	variables.Declare(x, language.IntegerSort)
	variables.Declare(y, language.IntegerSort)
	constraint := language.Disjunction(
		language.NewBinary(
			language.NewVariable(x),
			language.GreaterThanEqual,
			language.NewInteger(2),
		),
		language.Disjunction(
			language.NewBinary(
				language.NewVariable(y),
				language.LessThanEqual,
				language.NewInteger(1),
			),
			language.NewBinary(
				language.NewVariable(y),
				language.GreaterThanEqual,
				language.NewInteger(3),
			),
		),
	)
	invariant := NewInvariant(constraint)

	valuations := language.NewValuationsMap()
	solver := NewInterpreter(context, variables)

	for i := 0; i < 1000; i++ {
		// Act
		xVal, yVal := rand.Intn(1000)-500, rand.Intn(1000)-500
		valuations.Assign(x, language.NewInteger(xVal))
		valuations.Assign(y, language.NewInteger(yVal))
		satisfiable := invariant.IsSatisfiable(valuations, solver)

		// Assert
		expected := ((xVal >= 2) || (yVal <= 1 || yVal >= 3))
		if satisfiable != expected {
			assert.Equal(t, expected, satisfiable, fmt.Sprintf("Counter example with [x=%v, y=%v]", xVal, yVal))
		}
	}
}
