package main

import "fmt"

func main() {
	// fmt.Print("Enter domain: ")
	// reader := bufio.NewReader(os.Stdin)
	// // ReadString will block until the delimiter is entered
	// input, err := reader.ReadString('\n')
	// if err != nil {
	// 	fmt.Println("An error occured while reading input. Please try again", err)
	// 	return
	// }

	// // remove the delimeter from the string
	// input = strings.TrimSuffix(input, "\n")
	fmt.Println("Resolved IP for anikd.com:", Resolve("anikd.com", TYPE_A))
	fmt.Println("-----")
	fmt.Println("Resolved IP for arnabsen.dev:", Resolve("arnabsen.dev", TYPE_A))
}
