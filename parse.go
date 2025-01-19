package main


import "fmt"


const START = 0;
const EXP = 1;
const STATEMENT = 2;
const TERM = 3;
const MODIFIER = 4;
const ITEM = 5;
const GROUP = 6;
const CHAR = 7;
const TERMINAL = 8;


var MODIFIERS [3]string = [3]string{"*", "+", "?"}
var RESERVED_CHARS [7]string = [7]string{"*", "+", "?", "(", ")", "|"}


type symbol struct {
	symbol_type int
	regex_string string
	children []*symbol
}


func new_symbol(sym_type int, regex_string string) *symbol {
	sym := symbol{sym_type, regex_string, []*symbol{}}

	return &sym
}


func (sym *symbol) add_child(child *symbol) {
	sym.children = append(sym.children, child)
}


func (sym *symbol) get_terminal_string() string {
	//should only be called for symbol with a single child
	if sym.symbol_type == TERMINAL {
		return sym.regex_string
	} else {
		return sym.children[0].get_terminal_string()
	}
}


func (sym *symbol) get_replacements() [][]*symbol {
	/*
	regex context free language:

	START :== EXP
	EXP :== STATEMENT | STATEMENT "|" EXP
	STATEMENT :== TERM EXP | TERM
	TERM :== ITEM MODIFIER | ITEM
	MODIFIER :== "*" | "+" | "?"
	ITEM :== CHAR | GROUP
	GROUP :== "(" EXP ")"
	CHAR :== "a" | "b" | "c" | ...
	*/

	possible_replaces := [][]*symbol{}

	switch sym.symbol_type {
	case START:
		possible_replaces = append(possible_replaces, []*symbol{new_symbol(EXP, sym.regex_string)})
	case EXP:
		contains_or := false
		for i, x := range sym.regex_string {
			if x == '|' {
				contains_or = true

				statement_string := sym.regex_string[:i]
				exp_string := sym.regex_string[i + 1:]

				possible_replaces = append(possible_replaces, []*symbol{new_symbol(STATEMENT, statement_string), new_symbol(EXP, exp_string)})
			}
		}

		if !contains_or {possible_replaces = append(possible_replaces, []*symbol{new_symbol(STATEMENT, sym.regex_string)})}  //just the statement
	case STATEMENT:
		possible_replaces = append(possible_replaces, []*symbol{new_symbol(TERM, sym.regex_string)})  //just the term
	
		for i := 1; i < len(sym.regex_string); i++ {
			term_string := sym.regex_string[:i]
			exp_string := sym.regex_string[i:]
	
			possible_replaces = append(possible_replaces, []*symbol{new_symbol(TERM, term_string), new_symbol(EXP, exp_string)})
		}
	case TERM:
		possible_replaces = append(possible_replaces, []*symbol{new_symbol(ITEM, sym.regex_string)})  //just the item
	
		for i := 1; i < len(sym.regex_string); i++ {
			item_string := sym.regex_string[:i]
			modifier_string := sym.regex_string[i:]
	
			possible_replaces = append(possible_replaces, []*symbol{new_symbol(ITEM, item_string), new_symbol(MODIFIER, modifier_string)})
		}
	case MODIFIER:
		for _, i := range MODIFIERS {
			if i == sym.regex_string {possible_replaces = append(possible_replaces, []*symbol{new_symbol(TERMINAL, i)})}
		}
	case ITEM:
		possible_replaces = append(possible_replaces, []*symbol{new_symbol(GROUP, sym.regex_string)})
		possible_replaces = append(possible_replaces, []*symbol{new_symbol(CHAR, sym.regex_string)})
	case GROUP:
		if len(sym.regex_string) > 2 && sym.regex_string[0] == '(' && sym.regex_string[len(sym.regex_string) - 1] == ')' {
			exp_string := sym.regex_string[1 : len(sym.regex_string) - 1]
			possible_replaces = append(possible_replaces, []*symbol{new_symbol(EXP, exp_string)})
		}
	case CHAR:
		if is_regex_char(sym.regex_string) {
			possible_replaces = append(possible_replaces, []*symbol{new_symbol(TERMINAL, sym.regex_string)})
		}
	}

	return possible_replaces
}


func is_regex_char(regex string) bool {
	if len(regex) != 1 {return false}

	for _, i := range RESERVED_CHARS {
		if i == regex {return false}
	}

	return true
}


func print_tree(current_sym *symbol, depth int) {
	//for debugging only
	indent := ""
	for i := 0; i < depth; i++ {indent += "-"}

	for _, i := range current_sym.children {
		if i.symbol_type != TERMINAL {
			fmt.Printf("%s%v\n", indent, i.symbol_type)
			print_tree(i, depth + 1)
		} else {
			fmt.Printf("%s%s\n", indent, i.regex_string)
		}
	}
}


func build_tree(current_symbol *symbol) bool {
	if current_symbol.symbol_type == TERMINAL {return true}

	replacements := current_symbol.get_replacements()

	for i := 0; i < len(replacements); i++ {
		replace_symbols := replacements[i]

		is_valid_tree := true
		
		//range copies the structs in the replace_symbols slice, so we need to use the index instead
		for i := 0; i < len(replace_symbols); i++ {
			sym := replace_symbols[i]
			valid := build_tree(sym)

			if !valid {is_valid_tree = false}
		}

		if is_valid_tree {
			//actually add the symbols to the tree
			for i := 0; i < len(replace_symbols); i++ {
				sym := replace_symbols[i]
				current_symbol.add_child(sym)
			}

			return true
		}
	}

	return false
}


func parse_regex(regex string) symbol {
	root_symbol := symbol{START, regex, []*symbol{}}

	is_valid_regex := build_tree(&root_symbol)
	if !is_valid_regex {panic("Invalid regular expression")}

	return root_symbol
}