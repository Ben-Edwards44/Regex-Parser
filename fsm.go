package main


import "strconv"


type state struct {
	is_start bool
	is_accept bool
}


type transition struct {
	from *state
	to *state

	is_epsilon bool

	condition rune
}


type finite_state_machine struct {
	transitions []transition

	start *state
	accept *state
}


func (fsm *finite_state_machine) check_accept(sequence string, current_node *state) bool {
	if len(sequence) == 0 && current_node.is_accept {return true}  //we may need to search other epsilon transitions to find an accept state
	
	next_char := rune(0)
	if len(sequence) > 0 {next_char = rune(sequence[0])}

	for _, i := range fsm.transitions {
		if i.from == current_node && (i.is_epsilon || i.condition == next_char) {
			next_state := i.to

			var works bool
			if i.is_epsilon {
				works = fsm.check_accept(sequence, next_state)
			} else if i.condition == next_char {
				works = fsm.check_accept(sequence[1:], next_state)
			}

			if works {return true}
		}
	}

	return false
}


func (fsm *finite_state_machine) add_ghost_ends() {
	new_start := state{true, false}
	new_end := state{false, true}

	fsm.start.is_start = false
	fsm.accept.is_accept = false

	start_join_trans := transition{&new_start, fsm.start, true, rune(0)}
	end_join_trans := transition{fsm.accept, &new_end, true, rune(0)}

	fsm.transitions = append(fsm.transitions, start_join_trans)
	fsm.transitions = append(fsm.transitions, end_join_trans)

	fsm.start = &new_start
	fsm.accept = &new_end
}


func single_char(char rune) finite_state_machine {
	start := state{true, false}
	end := state{false, true}
	trans := transition{&start, &end, false, char}

	return finite_state_machine{[]transition{trans}, &start, &end}
}


func concatenate(left finite_state_machine, right finite_state_machine) finite_state_machine {
	a := left.accept
	b := right.start

	a.is_accept = false
	b.is_start = false

	trans := transition{a, b, true, rune(0)}

	transitions := left.transitions
	transitions = append(transitions, right.transitions...)
	transitions = append(transitions, trans)

	return finite_state_machine{transitions, left.start, right.accept}
}


func one_or_more(fsm finite_state_machine) finite_state_machine {
	repeat_trans := transition{fsm.accept, fsm.start, true, rune(0)}

	fsm.transitions = append(fsm.transitions, repeat_trans)

	return fsm
}


func zero_or_more(fsm finite_state_machine) finite_state_machine {
	fsm = one_or_more(fsm)
	fsm.add_ghost_ends()
	
	zero_occurence_trans := transition{fsm.start, fsm.accept, true, rune(0)}
	fsm.transitions = append(fsm.transitions, zero_occurence_trans)

	return fsm
}


func zero_or_one(fsm finite_state_machine) finite_state_machine {
	fsm.add_ghost_ends()

	zero_occurence_trans := transition{fsm.start, fsm.accept, true, rune(0)}
	fsm.transitions = append(fsm.transitions, zero_occurence_trans)
	
	return fsm
}


func apply_modifier(fsm finite_state_machine, mod_char string) finite_state_machine {
	switch mod_char {
	case "+": return one_or_more(fsm)
	case "*": return zero_or_more(fsm)
	case "?": return zero_or_one(fsm)
	default: panic("Invalid modifier: " + mod_char)
	}
}


func get_fsm(tree_symbol *symbol) finite_state_machine {
	switch tree_symbol.symbol_type {
	case START:
		return get_fsm(tree_symbol.children[0])
	case EXP:
		if len(tree_symbol.children) == 1 {return get_fsm(tree_symbol.children[0])}  //no additional expression

		left_fsm := get_fsm(tree_symbol.children[0])
		right_fsm := get_fsm(tree_symbol.children[1])

		return concatenate(left_fsm, right_fsm)
	case TERM:
		if len(tree_symbol.children) == 1 {return get_fsm(tree_symbol.children[0])}  //no modifier

		fsm := get_fsm(tree_symbol.children[0])
		mod_char := tree_symbol.children[1].get_terminal_string()

		return apply_modifier(fsm, mod_char)
	case ITEM:
		return get_fsm(tree_symbol.children[0])
	case GROUP:
		return get_fsm(tree_symbol.children[0])
	case CHAR:
		char := tree_symbol.get_terminal_string()

		return single_char(rune(char[0]))
	default:
		panic("Invalid symbol type for finite state machine creation: " + strconv.Itoa(tree_symbol.symbol_type))
	}
}