package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"strings"
	"time"
)

const TYPE_A = 1
const CLASS_IN = 1

type DNSHeader struct {
	id               uint16
	flags            uint16
	question_count   uint16
	answer_count     uint16
	authority_count  uint16
	additional_count uint16
}

type DNSQuestion struct {
	name       []byte
	query_type uint16
	class_type uint16
}

// converts header to byte string in big endian
// network packets always use big endian
func HeaderToBytes(header DNSHeader) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, header)
	return buf.Bytes()
}

func EncodeDomainName(domain string) []byte {
	var buf bytes.Buffer
	labels := strings.Split(domain, ".")
	for _, label := range labels {
		buf.WriteByte(byte(len(label)))
		buf.WriteString(label)
	}
	// end . at end of domain name
	buf.WriteByte(0)
	return buf.Bytes()
}

func QuestionToBytes(question DNSQuestion) []byte {
	buf := new(bytes.Buffer)
	// only write query type and class type
	binary.Write(buf, binary.BigEndian, question.query_type)
	binary.Write(buf, binary.BigEndian, question.class_type)
	return append(question.name, buf.Bytes()...)
}

func BuildQuery(domain string, recordType uint16) []byte {
	encodedDomain := EncodeDomainName(domain)
	var RECURSION_DESIRED uint16 = 1 << 8
	// id := uint16(rand.Intn(65535)) //0x8298
	header := DNSHeader{
		id:               0x8298,
		question_count:   1,
		flags:            RECURSION_DESIRED,
		answer_count:     0,
		authority_count:  0,
		additional_count: 0,
	}

	question := DNSQuestion{
		name:       encodedDomain,
		query_type: recordType,
		class_type: CLASS_IN,
	}

	headerByteString := HeaderToBytes(header)
	questionByteString := QuestionToBytes(question)
	query := append(headerByteString, questionByteString...)
	// queryByteString := fmt.Sprintf("%q", query)
	// fmt.Println(queryByteString)
	return query

}

func SendQuery(query []byte) {
	fmt.Println("hexencoded query ->", hex.EncodeToString(query))
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		fmt.Println(err)
		return
	}
	conn.Write(query)
	conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
	response := make([]byte, 1024)
	_, err = conn.Read(response)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	SendQuery(BuildQuery("example.com", TYPE_A))
}
