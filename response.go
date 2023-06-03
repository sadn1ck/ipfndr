package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"strings"
)

func DecompressDomainName(buf *bytes.Buffer, len int) string {
	// val, _ := buf.ReadByte()
	// var ptrBytes int = (len & 63) + int(val)

	// current := buf.Cap() - buf.Len()
	// log.Println("current ->", current)
	// for i := 0; i < current; i++ {
	// 	buf.UnreadByte()
	// }
	// log.Println("after unread", buf.Bytes())
	// return ""
	// result := DecodeDomainName(buf)
	// return result
	return ""
}

func DecodeDomainName(buf *bytes.Buffer) string {
	var parts []string
	c, _ := buf.ReadByte()
	for c != 0 {
		if c&192 != 0 {
			part := DecompressDomainName(buf, int(c))
			log.Println("part decomp->", part)
			parts = append(parts, part)
			break
		} else {
			var part []byte = buf.Next(int(c))
			parts = append(parts, string(part))
			c, _ = buf.ReadByte()
		}
	}
	return strings.Join(parts, ".")
}

// 🟢
func ParseResponseHeaders(buf *bytes.Buffer) DNSHeader {
	var header DNSHeader
	data := buf.Next(12)
	err := binary.Read(bytes.NewBuffer(data), binary.BigEndian, &header)
	if err != nil {
		log.Fatalln("brother binary read error ->", err)
	}
	return header
}

// 🟢
func ParseResponseQuestion(buf *bytes.Buffer) DNSQuestion {
	var q DNSQuestion
	q.QName = []byte(DecodeDomainName(buf))
	binary.Read(buf, binary.BigEndian, &q.QType)
	binary.Read(buf, binary.BigEndian, &q.QClass)
	return q
}

// 🔴
func ParseResponseRecord(buf *bytes.Buffer) DNSRecord {
	log.Println("buf for record ->", buf.Bytes())
	var record DNSRecord
	record.Name = []byte(DecodeDomainName(buf))
	// data := buf.Next(10)
	// binary.Read(bytes.NewBuffer(data), binary.BigEndian, &record.RecordType)
	// binary.Read(bytes.NewBuffer(data), binary.BigEndian, &record.ClassType)
	// binary.Read(bytes.NewBuffer(data), binary.BigEndian, &record.Ttl)
	// var dataLen byte
	// binary.Read(buf, binary.BigEndian, &dataLen)
	// record.Data = buf.Next(int(dataLen))
	// log.Println("type ->", record.RecordType, "class ->", record.ClassType, "ttl ->", record.Ttl, "data ->", record.Data)
	return record
}
