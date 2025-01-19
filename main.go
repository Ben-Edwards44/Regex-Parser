package main


import (
	"os"
	"fmt"
	"bufio"
)


const RED = "\033[31m"
const GRAY = "\033[38;5;238m"
const GREEN = "\033[32m"
const END_COLOUR = "\033[0m"

const INDENT = "\t\t"


func read_input(prompt string) string {
	fmt.Print(prompt)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	input := scanner.Text()

	if scanner.Err() != nil {panic(scanner.Err())}

	return input
}


func display_accept(accept bool) {
	text := RED + "rejected"
	if accept {text = GREEN + "accepted"}

	fmt.Println(INDENT + text + END_COLOUR)
}


func main() {
	regex := read_input(GRAY + "regex: " + END_COLOUR)
	root := parse_regex(regex)
	fsm := get_fsm(&root)

	for {
		test_string := read_input(GRAY + "match: " + END_COLOUR)
		accept := fsm.check_accept(test_string, fsm.start)

		if test_string == "quit" {break}

		display_accept(accept)
	}
}