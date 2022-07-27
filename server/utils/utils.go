package utils

import (
	"fmt"
	"net"
	"time"
)

var VALID_METHODS []string = []string{
	"UDP",
	"TCP",
	"TCP-HTTP",
	"VPN-DOWN",
	"HOME-FREEZE",
	"UDP-BYPASS",
	"HOME-STOMP",
	"STOP",
}

func Get_msg(connection net.Conn) string {
	buffer := make([]byte, 1024)
	msg, err := connection.Read(buffer)

	if err != nil {
		return "Failure getting message"
	}

	s := string(buffer[:msg])

	return s
}
func CmpSocketMessage(msg string, compare string) bool {
	//	temp := ""
	if len(msg) > 100 || len(msg) < len(compare) {
		return false
	}
	return msg[0:len(compare)] == compare
}
func ScanPort(addr string, port int) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", addr, port), time.Second)

	if err != nil {
		return false
	}
	if conn != nil {
		defer conn.Close()
	}
	return true
}
