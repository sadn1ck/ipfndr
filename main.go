package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
)

const TYPE_A = 1
const CLASS_IN = 1

type DNSHeader struct {
	Id            uint16
	Flags         uint16
	QuestionCount uint16
	AnswerCount   uint16
	NsCount       uint16 // name server records count
	ARCount       uint16 // additional records
}

type DNSQuestion struct {
	QName  []byte
	QType  uint16
	QClass uint16
}

type DNSRecord struct {
	Name       []byte
	RecordType uint16
	ClassType  uint16
	Ttl        int
	Data       []byte
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

	log.Print("request header -> ", header)
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
	conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
	response := make([]byte, 1024)
	_, err = conn.Read(response)
	if err != nil {
		log.Fatalln(err)
		return []byte{}, err
	}
	return response, nil
}

func ParseResponseHeaders(buf *bytes.Buffer) DNSHeader {
	var header DNSHeader
	data := buf.Next(12)
	err := binary.Read(bytes.NewBuffer(data), binary.BigEndian, &header)
	if err != nil {
		log.Fatalln("brother binary read error ->", err)
	}
	return header
}

func DecodeCompressedDomainName(buf *bytes.Buffer, len int) string {
	var domain string
	return domain
}

func DecodeDomainNameFromQuestion(buf *bytes.Buffer) string {
	var parts []string
	r, _ := buf.ReadByte()
	log.Println("v -> ", r)
	for r != 0 {
		if r&0b11000000 == 0b11000000 {
			// todo: handle compressed domain name
			// pointer counting
		} else {
			var part []byte = buf.Next(int(r))
			log.Println("part -> ", part)
			parts = append(parts, string(part))
			r, _ = buf.ReadByte()
		}

	}
	return strings.Join(parts, ".")
}

/**
In order to reduce the size of messages, the domain system utilizes a
compression scheme which eliminates the repetition of domain names in a
message by replacing it with a pointer to a location in the message.

The pointer takes the form of a two octet sequence:

	+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	| 1  1|                OFFSET                   |
	+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+

The first two bits are ones.  This allows a pointer to be distinguished
from a label, since the label must begin with two zero bits because
labels are restricted to 63 octets or less. The OFFSET field specifies an offset from
the start of the message.
*/

func ParseResponseQuestion(buf *bytes.Buffer) DNSQuestion {
	var q DNSQuestion
	// TODO: impl compressed domain decoding
	log.Println(DecodeDomainNameFromQuestion(buf))
	return q
}

func main() {
	// todo: error handling for send/build query
	response, err := SendQuery(BuildQuery("anikd.com", TYPE_A))
	if err != nil {
		log.Fatalln(err)
	} else {
		buf := bytes.NewBuffer(response)
		header := ParseResponseHeaders(buf)
		log.Print("response header -> ", header)
		question := ParseResponseQuestion(buf)
		log.Print("response question -> ", question)
	}
}
