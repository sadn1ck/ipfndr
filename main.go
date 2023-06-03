package main

import (
	"log"
)

func main() {
	// todo: error handling for send/build query
	buf, err := SendQuery(BuildQuery("chroniclehq.com", TYPE_A))
	if err != nil {
		log.Fatalln(err)
	} else {
		var offset int = 0
		log.Println(buf)
		// header := ParseResponseHeaders(buf)
		ParseResponseHeaders(buf, &offset)
		log.Println("offset: ", offset)

		// saved := buf.Bytes()
		// question := ParseResponseQuestion(buf)
		ParseResponseQuestion(buf, &offset)
		log.Println("offset: ", offset)
		// record := ParseResponseRecord(buf)
		ParseResponseRecord(buf, &offset)
		log.Println("offset: ", offset)
	}
}
