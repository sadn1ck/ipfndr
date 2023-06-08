package main

import (
	"log"
	"strings"
	"unicode"
)

func GetAnswers(packet DNSPacket) string {
	for _, answer := range packet.Answers {
		if answer.RecordType == TYPE_A {
			return (DataToIp(answer.Data))
		}
	}
	return ""
}

func GetNameserverIp(packet DNSPacket) string {
	for _, answer := range packet.Additional {
		if answer.RecordType == TYPE_A {
			return (DataToIp(answer.Data))
		}
	}
	return ""
}

func GetNameServer(packet DNSPacket) string {
	res := ""
	for _, answer := range packet.Authorities {
		var domain []byte
		if answer.RecordType == TYPE_NS {
			for _, v := range answer.Data {
				if unicode.IsSpace(rune(v)) {
					domain = append(domain, ' ')
				} else {
					domain = append(domain, v)
				}
			}
		}
		res = strings.Join(strings.Split(strings.TrimSpace(string(domain)), " "), ".")
	}
	return res
}

func GetCname(packet DNSPacket) string {
	for _, answer := range packet.Answers {
		if answer.RecordType == TYPE_CNAME {
			return string(answer.Data)
		}
	}
	return ""
}

func Resolve(domain string, record int) string {
	nameserver := "198.41.0.4" // e.root-servers.net
	ip := ""
	nsIp := ""
	nsDomain := ""
	cName := ""
	for {
		log.Println("Querying", nameserver, "for", domain)
		packet := BuildQuery(domain, TYPE_A, nameserver)
		ip = GetAnswers(packet)
		nsIp = GetNameserverIp(packet)
		nsDomain = GetNameServer(packet)
		cName = GetCname(packet)
		// log.Println("ip", ip, "nsIp", nsIp, "nsDomain", nsDomain, "cName", cName)
		if len(ip) > 0 {
			// base case
			return string(ip)
		} else if len(nsIp) > 0 {
			nameserver = nsIp
		} else if len(nsDomain) > 0 {
			nameserver = Resolve(nsDomain, TYPE_A)
		} else if len(cName) > 0 {
			domain = cName
			return Resolve(domain, TYPE_CNAME)
		} else {
			log.Fatalln("No nameserver found")
		}
	}
}
