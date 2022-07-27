package utils

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	ATTACK_KEY = "testkey"
)

func n_check(ip string, port string, time string, psize string) bool {

	port_int, port_err := strconv.Atoi(port)
	if port_err != nil {
		fmt.Println("Error setting port to integer on verification package")
		return false
	}
	psize_int, psize_err := strconv.Atoi(psize)
	if psize_err != nil {
		fmt.Println("Error setting psize to integer on verification package")
		return false
	}
	time_int, time_err := strconv.Atoi(time)
	if time_err != nil {
		fmt.Println("Error setting time to integer on verification package")
		return false
	}
	if port_int < 0 || port_int > 65535 {
		fmt.Println("Invalid port")
		return false
	} else if time_int < 10 || time_int > 4000 {
		fmt.Println("Invalid Time")
		return false
	} else if psize_int < 50 || psize_int > 1024 {
		fmt.Println("Invalid packet size")
		return false
	}
	return Valid_ip(ip)
}
func Valid_ip(ip string) bool {
	return net.ParseIP(ip) != nil
}

func Valid_User_Combo(user string, pass string) bool {
	file, _ := os.Open("auth/users.txt")
	fscanner := bufio.NewScanner(file)
	for fscanner.Scan() {
		if strings.Contains(fscanner.Text(), ":") {
			spl := strings.Split(fscanner.Text(), ":")
			username := spl[0]
			password := spl[1]

			if user == username && pass == password {
				return true
			} else {
				fmt.Printf("user != %s && pass != %s", username, password)
			}
		}
	}
	return false
}
func Check_params(conn net.Conn, key string, args []string) bool {

	if len(args) < 1 {
		fmt.Println("Invalid args")
		return false
	}
	if key != ATTACK_KEY {
		conn.Write([]byte("Invalid API Key"))
		fmt.Printf("%s Attempted to send attack without authentication", conn.RemoteAddr())
		return false
	}
	switch strings.ToLower(args[0]) { // method []args
	case "udp":
		if len(args) != 5 {
			fmt.Println(len(args))
			return false
		}
		ip := args[1]
		port := args[2]
		psize := args[3]
		time := args[4]
		return n_check(ip, port, time, psize) == true
	case "tcp-http": //ip IP -port PORT -time TIME -protocol Udp/tcp -osilayer 4/7 -threads 200
		if len(args) != 5 {
			fmt.Println(len(args))
			return false
		}
		ip := args[1]
		port := args[2]
		time := args[3]
		protocol := strings.ToUpper(args[len(args)-1])
		if protocol != "UDP" && protocol != "TCP" {
			fmt.Println("Invalid protocol in verification package, expected 4/7")
			return false
		}

		return n_check(ip, port, time, "1024") == true
	case "udp-bypass":
		if len(args) != 7 {
			return false
		}
		rpc, rpc_err := strconv.Atoi(args[1])
		if rpc_err != nil {
			fmt.Println("Error setting REQUESTS PER CONNECTION to integer.")
			return false
		}
		threads, threads_err := strconv.Atoi(args[2])
		if threads_err != nil {
			fmt.Println("Error setting threads to integer.")
			return false
		}

		ip := args[3]
		port := args[4]
		psize := args[5]
		time := args[6]

		if rpc > 10 || rpc < 0 {
			fmt.Println("RPC Must be 1-10")
			return false
		} else if threads < 1 || threads > 5 {
			fmt.Println("Threads must be 1-5")
			return false
		}

		return n_check(ip, port, time, psize) == true
	case "vpn-down":
		if len(args) != 6 {
			fmt.Println("Invalid args")
			return false
		}
		tc, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Error setting thread count to integer.")
			return false
		}
		if tc < 1 || tc > 5 {
			fmt.Println("Invalid thread count. Expected 1-5")
			return false
		}
		ip := args[2]
		port := args[3]
		psize := args[4]
		time := args[5]

		return n_check(ip, port, time, psize) == true
	case "home-stomp": // ip port psize time threads rpc
		if len(args) != 7 {
			fmt.Println("from util bad len")
			fmt.Println(len(args))
			return false
		}
		ip := args[1]
		port := args[2]
		psize := args[3]
		time := args[4]

		threads, err := strconv.Atoi(args[5])
		if err != nil {
			fmt.Println("Error setting threads on HOME-STOMP to int, returning false")
			return false
		}

		rpc, err := strconv.Atoi(args[6])
		if err != nil {
			fmt.Println("Error setting RPC to int, returning false")
			return false
		}
		if threads < 1 || threads > 5 {
			fmt.Println("Threads must be 1-5 on HOME-STOMP")
			return false
		} else if rpc < 0 || rpc > 25 {
			fmt.Println("RPC must be 1-25 on HOME-STOMP")
			return false
		}

		return n_check(ip, port, time, psize) == true
	case "stop":
		if len(args) != 2 {
			return false
		}
		return Valid_ip(args[1])
	default:
		fmt.Printf("Method not found in verification package, returning false (%s)", args[0])
		return false
	}
}
