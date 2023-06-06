package main

const (
	TYPE_A     = 1
	TYPE_NS    = 2
	TYPE_CNAME = 5
	TYPE_SOA   = 6
	TYPE_PTR   = 12
	TYPE_MX    = 15
	TYPE_TXT   = 16
	TYPE_AAAA  = 28
	TYPE_SRV   = 33
	//
	CLASS_IN = 1
)

type DNSHeader struct {
	Id             uint16
	Flags          uint16
	QuestionCount  uint16
	AnswerCount    uint16
	AuthorityCount uint16 // name server records count
	AddnCount      uint16 // additional records
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

type DNSPacket struct {
	Header      DNSHeader
	Questions   []DNSQuestion
	Answers     []DNSRecord
	Authorities []DNSRecord
	Additional  []DNSRecord
}
