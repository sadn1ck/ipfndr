package main

import (
	"strconv"
	"strings"
)

// is it even a project without a utils file?

func DataToIP(data []byte) string {
	parts := [4]string{}
	for i := 0; i < 4; i++ {
		parts[i] = strconv.Itoa(int(data[i]))
	}
	return strings.Join(parts[:], ".")
}
