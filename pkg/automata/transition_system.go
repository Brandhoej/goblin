package automata

type TransitionSystem struct {
	solver    *ConstraintSolver
	automaton *Automaton
}

func NewTransitionSystem(automaton *Automaton) *TransitionSystem {
	return &TransitionSystem{
		automaton: automaton,
	}
}

func (system *TransitionSystem) Initial() State {
	key := system.automaton.initial
	location, _ := system.automaton.Location(key)
	return NewState(key, location.invariant.constraint)
}

// Returns all states from the state.
func (system *TransitionSystem) Outgoing(state State) (successors []State) {
	if location, exists := system.automaton.Location(state.location); exists {
		// We have found an inconsistency where the location is disabled.
		// Meaning that even enabled edges wont be traversable.
		if !location.IsEnabled(system.solver) {
			return successors
		}
	} else {
		panic("State is in an unkown location")
	}

	edges := system.automaton.Outgoing(state.location)
	for _, edge := range edges {
		// Check if we can even traverse the edge.
		if !edge.IsEnabled(system.solver) {
			continue
		}

		// We can traverse the edge so we create a new and updated state.
		state := edge.Traverse(state)
		successors = append(successors, state)
	}
	return successors
}

func (system *TransitionSystem) Reachability(solver *ConstraintSolver, search SearchStrategy, goals ...State) Trace {
	return search.For(
		func(state State) bool {
			// We have reached a goal when the locations are the same
			// and the goal contains (Meaning that more valuations or the same) are possible.
			for _, goal := range goals {
				if goal.SubsetOf(state, solver) {
					return true
				}
			}

			return false
		},
		system.Initial(),
	)
}
