package main

import (
	"encoding/binary"
	"log"
	"strings"
)

func DecompressDomainName(response []byte, offset *int, len int) string {
	// val, _ := buf.ReadByte()
	val := response[*offset]
	*offset++
	var ptrBytes int = (len & 63) + int(val)
	currentOffset := *offset
	*offset = ptrBytes

	domain := DecodeDomainName(response, offset)
	*offset = currentOffset // reset offset post usage
	return domain
}

func DecodeDomainName(buf []byte, offset *int) string {
	var parts []string
	c := buf[*offset]
	*offset++
	for c != 0 {
		if c&192 != 0 {
			part := DecompressDomainName(buf, offset, int(c))
			parts = append(parts, part)
			break
		} else {
			var part []byte = buf[*offset : *offset+int(c)]
			parts = append(parts, string(part))
			*offset += int(c)
			c = buf[*offset]
			*offset++
		}
	}
	return strings.Join(parts, ".")
}

func ParseResponseHeaders(buf []byte, offset *int) DNSHeader {
	var header DNSHeader
	data := buf[:12]
	header.Id = binary.BigEndian.Uint16(data[:2])
	header.Flags = binary.BigEndian.Uint16(data[2:4])
	header.QuestionCount = binary.BigEndian.Uint16(data[4:6])
	header.AnswerCount = binary.BigEndian.Uint16(data[6:8])
	header.AuthorityCount = binary.BigEndian.Uint16(data[8:10])
	header.AddnCount = binary.BigEndian.Uint16(data[10:12])
	// err := binary.Read(bytes.NewBuffer(data), binary.BigEndian, &header)
	// if err != nil {
	// 	log.Fatalln("brother binary read error ->", err)
	// }
	*offset += 12
	return header
}

func ParseResponseQuestion(buf []byte, offset *int) DNSQuestion {
	var q DNSQuestion
	q.QName = []byte(DecodeDomainName(buf, offset))
	data := buf[*offset : *offset+4]
	*offset += 4
	q.QType = binary.BigEndian.Uint16(data[:2])
	q.QClass = binary.BigEndian.Uint16(data[2:4])
	return q
}

func ParseResponseRecord(buf []byte, offset *int) DNSRecord {
	var record DNSRecord
	record.Name = []byte(DecodeDomainName(buf, offset))
	data := buf[*offset : *offset+10]
	*offset += 10
	record.RecordType = binary.BigEndian.Uint16(data[:2])
	record.ClassType = binary.BigEndian.Uint16(data[2:4])
	ttl := binary.BigEndian.Uint32(data[4:8])
	record.Ttl = int(ttl)
	dataLen := binary.BigEndian.Uint16(data[8:10])

	if record.RecordType == TYPE_NS {
		if len(string(record.Data)) > 0 {
			log.Println("NS record ->", record.Data, string(record.Data))
		}
		record.Data = []byte(DecodeDomainName(buf, offset))
	} else if record.RecordType == TYPE_A {
		record.Data = buf[*offset : *offset+int(dataLen)]
		*offset += int(dataLen)
	} else if record.RecordType == TYPE_CNAME {
		record.Data = []byte(DecodeDomainName(buf, offset))
		*offset += int(dataLen)
	} else {
		record.Data = buf[*offset : *offset+int(dataLen)]
		*offset += int(dataLen)
		// log.Println("other record ->", record.Data, string(record.Data))
	}
	return record
}

func ParseDNSPacket(buf []byte) DNSPacket {
	var packet DNSPacket
	var offset int = 0
	packet.Header = ParseResponseHeaders(buf, &offset)
	for i := 0; i < int(packet.Header.QuestionCount); i++ {
		packet.Questions = append(packet.Questions, ParseResponseQuestion(buf, &offset))
	}
	var savedOffset int = offset
	var i uint16
	for i = 0; i < packet.Header.AnswerCount; i++ {
		packet.Answers = append(packet.Answers, ParseResponseRecord(buf, &savedOffset))
	}
	for i = 0; i < packet.Header.AuthorityCount; i++ {
		packet.Authorities = append(packet.Authorities, ParseResponseRecord(buf, &savedOffset))
	}
	for i = 0; i < packet.Header.AddnCount; i++ {
		packet.Additional = append(packet.Additional, ParseResponseRecord(buf, &savedOffset))
	}
	return packet
}
