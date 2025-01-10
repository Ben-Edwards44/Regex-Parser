package main


import "fmt"


func read_input(prompt string) string {
	fmt.Print(prompt)

	var input string

	_, err := fmt.Scanln(&input)
	if (err != nil) {panic(err)}
	

	return input
}


func main() {
	regex := read_input("Enter regex: ")
	root := parse_regex(regex)
	fsm := get_fsm(&root)

	quit := false
	for !quit {
		test_string := read_input("Enter string: ")
		accept := fsm.check_accept(test_string, fsm.start)

		if accept {
			fmt.Println("Accepted")
		} else {
			fmt.Println("Rejected")
		}
	}
}