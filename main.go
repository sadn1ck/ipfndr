package main

import (
	"bytes"
	"log"
)

func main() {
	// todo: error handling for send/build query
	response, err := SendQuery(BuildQuery("chroniclehq.com", TYPE_A))
	if err != nil {
		log.Fatalln(err)
	} else {
		buf := bytes.NewBuffer(response)
		log.Println("response length-> ", buf.Cap())
		log.Println("response -> ", buf.Bytes())
		// header := ParseResponseHeaders(buf)
		ParseResponseHeaders(buf)
		// log.Println("header -> ", header)

		// saved := buf.Bytes()
		// question := ParseResponseQuestion(buf)
		ParseResponseQuestion(buf)
		// log.Println("question -> ", question.QName)
		// record := ParseResponseRecord(buf)
		ParseResponseRecord(buf)
		// log.Println("response record -> ", record)
	}
}
