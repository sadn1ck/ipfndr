package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	fmt.Print("Enter domain: ")
	reader := bufio.NewReader(os.Stdin)
	// ReadString will block until the delimiter is entered
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("An error occured while reading input. Please try again", err)
		return
	}

	// remove the delimeter from the string
	input = strings.TrimSuffix(input, "\n")
	buf, err := SendQuery(BuildQuery(input, TYPE_A))
	if err != nil {
		log.Fatalln(err)
	} else {
		packet := ParseDNSPacket(buf)
		fmt.Println("IP: ", DataToIP(packet.Answers[0].Data))
	}
}
