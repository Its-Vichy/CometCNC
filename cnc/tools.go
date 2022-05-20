package main

import (
	"net"
	"strconv"
)

func CheckPort(address string, port int) bool {
	conn, err := net.Dial("tcp", address+":"+strconv.Itoa(port))

	if err != nil {
		return false
	}

	conn.Close()
	return true
}
