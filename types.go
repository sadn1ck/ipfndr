package main

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
