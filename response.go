package main

import (
	"bytes"
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

// 🟢
func ParseResponseHeaders(buf []byte, offset *int) DNSHeader {
	var header DNSHeader
	data := buf[:12]
	*offset += 12
	err := binary.Read(bytes.NewBuffer(data), binary.BigEndian, &header)
	if err != nil {
		log.Fatalln("brother binary read error ->", err)
	}
	return header
}

// 🟢
func ParseResponseQuestion(buf []byte, offset *int) DNSQuestion {
	var q DNSQuestion
	q.QName = []byte(DecodeDomainName(buf, offset))
	data := buf[*offset : *offset+4]
	*offset += 4
	q.QType = binary.BigEndian.Uint16(data[:2])
	q.QClass = binary.BigEndian.Uint16(data[2:4])
	return q
}

// 🔴
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
	record.Data = buf[*offset : *offset+int(dataLen)]
	*offset += int(dataLen)
	return record
}
