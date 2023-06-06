package main

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"net"
	"strings"
)

func BuildQuery(domain string, recordType uint16, nameserver string) DNSPacket {
	encodedDomain := EncodeDomainName(domain)
	id := uint16(rand.Intn(65535)) //0x8298
	header := DNSHeader{
		Id:             id,
		QuestionCount:  1,
		Flags:          0,
		AnswerCount:    0,
		AuthorityCount: 0,
		AddnCount:      0,
	}

	question := DNSQuestion{
		QName:  encodedDomain,
		QType:  recordType,
		QClass: CLASS_IN,
	}

	headerByteString := HeaderToBytes(header)
	questionByteString := QuestionToBytes(question)
	query := append(headerByteString, questionByteString...)
	packet, _ := SendQueryToNS(query, nameserver)
	return packet
}

func SendQueryToNS(query []byte, dnsServerIp string) (DNSPacket, error) {
	var dnsIpPort string = dnsServerIp + ":53"
	conn, _ := net.Dial("udp", dnsIpPort)
	conn.Write(query)
	// read into 1024 byte buffer
	res := make([]byte, 1024)
	n, _ := conn.Read(res)
	// copy response into new fixed size byte array
	response := make([]byte, n)
	copy(response, res)
	// directly return response
	return ParseDNSPacket(response), nil
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
