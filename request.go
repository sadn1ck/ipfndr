package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"math/rand"
	"net"
	"strings"
)

func BuildQuery(domain string, recordType uint16) []byte {
	encodedDomain := EncodeDomainName(domain)
	var RECURSION_DESIRED uint16 = 1 << 8
	id := uint16(rand.Intn(65535)) //0x8298
	header := DNSHeader{
		Id:            id,
		QuestionCount: 1,
		Flags:         RECURSION_DESIRED,
		AnswerCount:   0,
		NsCount:       0,
		ARCount:       0,
	}

	question := DNSQuestion{
		QName:  encodedDomain,
		QType:  recordType,
		QClass: CLASS_IN,
	}

	// log.Print("request header -> ", header)
	headerByteString := HeaderToBytes(header)
	questionByteString := QuestionToBytes(question)
	query := append(headerByteString, questionByteString...)
	return query

}

func SendQuery(query []byte) ([]byte, error) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		log.Fatalln(err)
		return []byte{}, err
	}
	conn.Write(query)
	res := make([]byte, 1024)
	n, err := conn.Read(res)
	response := make([]byte, n)
	copy(response, res)
	if err != nil {
		log.Fatalln(err)
		return []byte{}, err
	}
	return response, nil
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
	binary.Write(buf, binary.BigEndian, question.QType)
	binary.Write(buf, binary.BigEndian, question.QClass)
	return append(question.QName, buf.Bytes()...)
}
